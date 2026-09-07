// Comando api — ponto de entrada do backend do IncidentTrack.
//
// Responsabilidade única: ler a configuração, montar o grafo de
// dependências (injeção manual — quem escolhe as implementações
// concretas é aqui, não os usecases; DIP, docs/plano.md §4) e subir o
// servidor HTTP com encerramento gracioso.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"incidenttrack/internal/config"
	httpapi "incidenttrack/internal/handler/http"
	"incidenttrack/internal/middleware"
	"incidenttrack/internal/pkg/hash"
	"incidenttrack/internal/pkg/token"
	"incidenttrack/internal/repository/sqlite"
	authuc "incidenttrack/internal/usecase/auth"
	incuc "incidenttrack/internal/usecase/incident"
)

func main() {
	// Subcomando usado pelo HEALTHCHECK do container (imagem distroless,
	// sem shell/curl): `api healthcheck` faz um GET em /api/health e sai
	// com código 0/1.
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		os.Exit(healthcheck())
	}
	if err := run(); err != nil {
		slog.Error("falha fatal na inicialização", "err", err)
		os.Exit(1)
	}
}

// healthcheck bate no próprio endpoint de liveness. Deriva a porta de
// HTTP_ADDR (default :8080).
func healthcheck() int {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		port = "8080"
	}
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://127.0.0.1:" + port + "/api/health")
	if err != nil {
		return 1
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}

func run() error {
	// .env é conveniência de dev; em produção as vars vêm do ambiente.
	if err := config.LoadDotEnv(".env"); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log := newLogger(cfg)
	slog.SetDefault(log)

	// Contexto encerrado em SIGINT/SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// --- Infraestrutura: banco ---
	if err := ensureDBDir(cfg.DatabaseURL); err != nil {
		return err
	}
	db, err := sqlite.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()
	log.Info("banco pronto", "dsn", redactDSN(cfg.DatabaseURL))

	// --- Adaptadores / serviços de pkg ---
	hasher := hash.NewBcrypt(cfg.BcryptCost)
	tokens := token.NewHS256(cfg.JWTSecret, cfg.AccessTokenTTL)
	denylist := token.NewDenylist(time.Minute)

	// Hash "dummy" pré-computado: comparado no login quando o e-mail não
	// existe, para nivelar o tempo de resposta (anti-enumeração, A07).
	dummyHash, err := hasher.Hash("dummy-password-for-timing-safety-01")
	if err != nil {
		return err
	}

	// --- Repositórios (implementações concretas) ---
	userRepo := sqlite.NewUserRepo(db)
	incidentRepo := sqlite.NewIncidentRepo(db)

	// --- Usecases ---
	registerUC := authuc.NewRegisterUser(userRepo, hasher)
	loginUC := authuc.NewLoginUser(userRepo, hasher, dummyHash)

	// --- Rate limiters ---
	apiLimiter := middleware.NewLimiter(cfg.APIRateLimit, cfg.APIRateWindow)
	loginLimiter := middleware.NewLimiter(cfg.LoginRateLimit, cfg.LoginRateWindow)

	// Atrás de um proxy reverso confiável (Nginx no deploy, serviço "web"
	// no docker-compose) confiamos no X-Forwarded-For para o rate limiting.
	trustProxy := cfg.TrustProxy

	// --- Handlers ---
	authHandler := httpapi.NewAuthHandler(httpapi.AuthHandlerDeps{
		Register: registerUC,
		Login:    loginUC,
		Users:    userRepo,
		Tokens:   tokens,
		Denylist: denylist,
		Cookie: httpapi.CookieConfig{
			Name:   cfg.CookieName,
			Secure: cfg.CookieSecure,
			Domain: cfg.CookieDomain,
			TTL:    cfg.AccessTokenTTL,
		},
		LoginLimiter: loginLimiter,
		TrustProxy:   trustProxy,
		Log:          log,
	})

	incidentHandler := httpapi.NewIncidentHandler(httpapi.IncidentHandlerDeps{
		Create: incuc.NewCreateIncident(incidentRepo),
		List:   incuc.NewListIncidents(incidentRepo),
		Get:    incuc.NewGetIncident(incidentRepo),
		Update: incuc.NewUpdateIncident(incidentRepo),
		Delete: incuc.NewDeleteIncident(incidentRepo),
		Log:    log,
	})

	authenticator := middleware.NewAuthenticator(tokens, denylist, cfg.CookieName)

	router := httpapi.NewRouter(httpapi.RouterDeps{
		Auth:          authHandler,
		Incidents:     incidentHandler,
		Authenticator: authenticator,
		APILimiter:    apiLimiter,
		LoginLimiter:  loginLimiter,
		AllowedOrigin: cfg.CORSAllowedOrigin,
		TrustProxy:    trustProxy,
		IsProduction:  cfg.IsProduction(),
		Log:           log,
	})

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelError),
	}

	// Sobe o servidor numa goroutine e espera o sinal de parada.
	errCh := make(chan error, 1)
	go func() {
		log.Info("servidor ouvindo", "addr", cfg.HTTPAddr, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Info("sinal de encerramento recebido, drenando conexões")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	log.Info("servidor encerrado com sucesso")
	return nil
}

func newLogger(cfg *config.Config) *slog.Logger {
	level := slog.LevelInfo
	if !cfg.IsProduction() {
		level = slog.LevelDebug
	}
	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level}
	if cfg.IsProduction() {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}
	return slog.New(handler)
}

// ensureDBDir cria o diretório do arquivo SQLite, se o DSN for baseado
// em arquivo ("file:...").
func ensureDBDir(dsn string) error {
	path := strings.TrimPrefix(dsn, "file:")
	if i := strings.IndexByte(path, '?'); i >= 0 {
		path = path[:i]
	}
	if path == "" || path == ":memory:" || strings.HasPrefix(path, ":memory:") {
		return nil
	}
	dir := filepath.Dir(path)
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0o750)
}

// redactDSN evita registrar caminhos completos/segredos em log.
func redactDSN(dsn string) string {
	if i := strings.IndexByte(dsn, '?'); i >= 0 {
		return dsn[:i] + "?<pragmas>"
	}
	return dsn
}
