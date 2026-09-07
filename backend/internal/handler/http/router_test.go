package httpapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"incidenttrack/internal/domain/entity"
	httpapi "incidenttrack/internal/handler/http"
	"incidenttrack/internal/middleware"
	"incidenttrack/internal/pkg/hash"
	"incidenttrack/internal/pkg/token"
	"incidenttrack/internal/repository/memory"
	authuc "incidenttrack/internal/usecase/auth"
	incuc "incidenttrack/internal/usecase/incident"
)

const origin = "http://localhost:5173"

// buildServer monta o router completo com repositórios em memória.
func buildServer(t *testing.T) *httptest.Server {
	t.Helper()

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	users := memory.NewUserRepo()
	incidents := memory.NewIncidentRepo()
	hasher := hash.NewBcrypt(4)
	tokens := token.NewHS256([]byte("0123456789abcdef0123456789abcdef"), 15*time.Minute)
	denylist := token.NewDenylist(time.Minute)
	dummy, _ := hasher.Hash("dummy-000000-timing")

	authHandler := httpapi.NewAuthHandler(httpapi.AuthHandlerDeps{
		Register: authuc.NewRegisterUser(users, hasher),
		Login:    authuc.NewLoginUser(users, hasher, dummy),
		Users:    users,
		Tokens:   tokens,
		Denylist: denylist,
		Cookie: httpapi.CookieConfig{
			Name: "it_session", Secure: false, TTL: 15 * time.Minute,
		},
		LoginLimiter: middleware.NewLimiter(50, time.Minute),
		Log:          log,
	})
	incidentHandler := httpapi.NewIncidentHandler(httpapi.IncidentHandlerDeps{
		Create: incuc.NewCreateIncident(incidents),
		List:   incuc.NewListIncidents(incidents),
		Get:    incuc.NewGetIncident(incidents),
		Update: incuc.NewUpdateIncident(incidents),
		Delete: incuc.NewDeleteIncident(incidents),
		Log:    log,
	})

	router := httpapi.NewRouter(httpapi.RouterDeps{
		Auth:          authHandler,
		Incidents:     incidentHandler,
		Authenticator: middleware.NewAuthenticator(tokens, denylist, "it_session"),
		APILimiter:    middleware.NewLimiter(1000, time.Minute),
		LoginLimiter:  middleware.NewLimiter(50, time.Minute),
		AllowedOrigin: origin,
		Log:           log,
	})
	srv := httptest.NewServer(router)
	t.Cleanup(srv.Close)
	return srv
}

func newClient(t *testing.T) *http.Client {
	t.Helper()
	jar, _ := cookiejar.New(nil)
	return &http.Client{Jar: jar}
}

func do(t *testing.T, c *http.Client, method, url string, body any, csrf bool) *http.Response {
	t.Helper()
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(context.Background(), method, url, rdr)
	if err != nil {
		t.Fatalf("req: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Origin", origin)
	if csrf {
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
	}
	res, err := c.Do(req)
	if err != nil {
		t.Fatalf("do: %v", err)
	}
	return res
}

func decode[T any](t *testing.T, res *http.Response) T {
	t.Helper()
	defer res.Body.Close()
	var v T
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return v
}

func registerAndLogin(t *testing.T, srv *httptest.Server, email string) *http.Client {
	t.Helper()
	c := newClient(t)

	res := do(t, c, http.MethodPost, srv.URL+"/api/auth/register", map[string]string{
		"name": "Fulano", "email": email, "password": "senha1234",
	}, true)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("register: status %d", res.StatusCode)
	}
	res.Body.Close()

	res = do(t, c, http.MethodPost, srv.URL+"/api/auth/login", map[string]string{
		"email": email, "password": "senha1234",
	}, true)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("login: status %d", res.StatusCode)
	}
	res.Body.Close()
	return c
}

func sampleIncident() map[string]string {
	return map[string]string{
		"titulo":     "Acesso suspeito ao servidor",
		"descricao":  "Login fora do horário comercial.",
		"categoria":  string(entity.CategoriaAcessoIndevido),
		"severidade": string(entity.SeveridadeMedia),
		"status":     string(entity.StatusAberto),
	}
}

func TestCRUDFlow(t *testing.T) {
	srv := buildServer(t)
	ana := registerAndLogin(t, srv, "ana@example.com")

	// cria
	res := do(t, ana, http.MethodPost, srv.URL+"/api/incidents", sampleIncident(), true)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create: status %d", res.StatusCode)
	}
	created := decode[map[string]any](t, res)
	id, _ := created["id"].(string)
	if id == "" {
		t.Fatal("create: sem id")
	}

	// lista
	res = do(t, ana, http.MethodGet, srv.URL+"/api/incidents", nil, false)
	list := decode[[]map[string]any](t, res)
	if len(list) != 1 {
		t.Fatalf("list: esperava 1, veio %d", len(list))
	}

	// detalha
	res = do(t, ana, http.MethodGet, srv.URL+"/api/incidents/"+id, nil, false)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("get: status %d", res.StatusCode)
	}
	res.Body.Close()

	// atualiza
	upd := sampleIncident()
	upd["titulo"] = "Acesso suspeito confirmado"
	upd["status"] = string(entity.StatusEmAnalise)
	res = do(t, ana, http.MethodPut, srv.URL+"/api/incidents/"+id, upd, true)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("update: status %d", res.StatusCode)
	}
	got := decode[map[string]any](t, res)
	if got["titulo"] != "Acesso suspeito confirmado" {
		t.Errorf("update não aplicado: %v", got["titulo"])
	}

	// remove
	res = do(t, ana, http.MethodDelete, srv.URL+"/api/incidents/"+id, nil, true)
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("delete: status %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestOwnershipIsolation(t *testing.T) {
	srv := buildServer(t)
	ana := registerAndLogin(t, srv, "ana@example.com")
	bruno := registerAndLogin(t, srv, "bruno@example.com")

	res := do(t, ana, http.MethodPost, srv.URL+"/api/incidents", sampleIncident(), true)
	id := decode[map[string]any](t, res)["id"].(string)

	// Bruno não enxerga o incidente da Ana — deve receber 404 (não 403,
	// para não confirmar a existência do recurso).
	res = do(t, bruno, http.MethodGet, srv.URL+"/api/incidents/"+id, nil, false)
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("get alheio: esperava 404, veio %d", res.StatusCode)
	}
	res.Body.Close()

	res = do(t, bruno, http.MethodDelete, srv.URL+"/api/incidents/"+id, nil, true)
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("delete alheio: esperava 404, veio %d", res.StatusCode)
	}
	res.Body.Close()

	// A lista do Bruno continua vazia.
	res = do(t, bruno, http.MethodGet, srv.URL+"/api/incidents", nil, false)
	if l := decode[[]map[string]any](t, res); len(l) != 0 {
		t.Fatalf("lista do Bruno deveria estar vazia, veio %d", len(l))
	}
}

func TestUnauthenticatedIsRejected(t *testing.T) {
	srv := buildServer(t)
	c := newClient(t)
	res := do(t, c, http.MethodGet, srv.URL+"/api/incidents", nil, false)
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("esperava 401, veio %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestCSRFHeaderRequiredOnMutations(t *testing.T) {
	srv := buildServer(t)
	ana := registerAndLogin(t, srv, "ana@example.com")

	// POST sem X-Requested-With deve ser bloqueado pelo middleware CSRF.
	res := do(t, ana, http.MethodPost, srv.URL+"/api/incidents", sampleIncident(), false)
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("esperava 403 sem header CSRF, veio %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestLogoutRevokesSession(t *testing.T) {
	srv := buildServer(t)
	ana := registerAndLogin(t, srv, "ana@example.com")

	// /me funciona antes do logout.
	res := do(t, ana, http.MethodGet, srv.URL+"/api/auth/me", nil, false)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("me antes do logout: %d", res.StatusCode)
	}
	res.Body.Close()

	res = do(t, ana, http.MethodPost, srv.URL+"/api/auth/logout", nil, true)
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("logout: %d", res.StatusCode)
	}
	res.Body.Close()

	// Cookie limpo → /me agora 401.
	res = do(t, ana, http.MethodGet, srv.URL+"/api/auth/me", nil, false)
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("me após logout: esperava 401, veio %d", res.StatusCode)
	}
	res.Body.Close()
}

func TestRegisterValidation(t *testing.T) {
	srv := buildServer(t)
	c := newClient(t)
	res := do(t, c, http.MethodPost, srv.URL+"/api/auth/register", map[string]string{
		"name": "x", "email": "nao-eh-email", "password": "curta",
	}, true)
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("esperava 422, veio %d", res.StatusCode)
	}
	body := decode[struct {
		Error  string              `json:"error"`
		Fields map[string][]string `json:"fields"`
	}](t, res)
	for _, f := range []string{"name", "email", "password"} {
		if len(body.Fields[f]) == 0 {
			t.Errorf("esperava erro no campo %q", f)
		}
	}
}
