package httpapi

import (
	"log/slog"
	"net/http"

	"incidenttrack/internal/domain/repository"
	"incidenttrack/internal/dto"
	"incidenttrack/internal/middleware"
	"incidenttrack/internal/pkg/httpx"
	"incidenttrack/internal/pkg/token"
	authuc "incidenttrack/internal/usecase/auth"
)

// AuthHandler expõe registro, login, logout e "quem sou eu".
type AuthHandler struct {
	register *authuc.RegisterUser
	login    *authuc.LoginUser
	users    repository.UserRepository
	tokens   token.Service
	denylist *token.Denylist
	cookie   CookieConfig
	// loginLimiter é o mesmo limitador aplicado como middleware na rota de
	// login; o handler o zera após autenticação bem-sucedida.
	loginLimiter *middleware.Limiter
	trustProxy   bool
	log          *slog.Logger
}

// AuthHandlerDeps agrupa as dependências para manter a assinatura enxuta.
type AuthHandlerDeps struct {
	Register     *authuc.RegisterUser
	Login        *authuc.LoginUser
	Users        repository.UserRepository
	Tokens       token.Service
	Denylist     *token.Denylist
	Cookie       CookieConfig
	LoginLimiter *middleware.Limiter
	TrustProxy   bool
	Log          *slog.Logger
}

// NewAuthHandler monta o handler de autenticação.
func NewAuthHandler(d AuthHandlerDeps) *AuthHandler {
	return &AuthHandler{
		register:     d.Register,
		login:        d.Login,
		users:        d.Users,
		tokens:       d.Tokens,
		denylist:     d.Denylist,
		cookie:       d.Cookie,
		loginLimiter: d.LoginLimiter,
		trustProxy:   d.TrustProxy,
		log:          d.Log,
	}
}

// Register — POST /api/auth/register. Cria o usuário; não faz login
// automático (o frontend redireciona para a tela de login).
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in dto.RegisterInput
	if !decodeJSON(w, r, &in) {
		return
	}
	in.Normalize()
	if errs := in.Validate(); !errs.Ok() {
		httpx.WriteValidationError(w, "Dados inválidos.", errs.Fields())
		return
	}

	u, err := h.register.Execute(r.Context(), authuc.RegisterInput{
		Name:     in.Name,
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		writeDomainError(w, r, h.log, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto.NewUserResponse(u))
}

// Login — POST /api/auth/login. Emite o cookie de sessão. A rota já
// passa por rate limit por IP (middleware); aqui, no sucesso, zeramos o
// contador daquele IP.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in dto.LoginInput
	if !decodeJSON(w, r, &in) {
		return
	}
	in.Normalize()
	if errs := in.Validate(); !errs.Ok() {
		httpx.WriteValidationError(w, "Dados inválidos.", errs.Fields())
		return
	}

	u, err := h.login.Execute(r.Context(), authuc.LoginInput{
		Email:    in.Email,
		Password: in.Password,
	})
	if err != nil {
		writeDomainError(w, r, h.log, err)
		return
	}

	raw, _, err := h.tokens.Generate(u.ID)
	if err != nil {
		writeDomainError(w, r, h.log, err)
		return
	}
	h.cookie.setSessionCookie(w, raw)

	if h.loginLimiter != nil {
		h.loginLimiter.Reset(middleware.ClientIP(r, h.trustProxy))
	}

	httpx.WriteJSON(w, http.StatusOK, dto.NewUserResponse(u))
}

// Logout — POST /api/auth/logout. Limpa o cookie e revoga o jti atual na
// denylist em memória, invalidando o token antes do vencimento natural.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if claims, ok := middleware.ClaimsFrom(r.Context()); ok && h.denylist != nil {
		h.denylist.Revoke(claims.TokenID, claims.ExpiresAt)
	}
	h.cookie.clearSessionCookie(w)
	httpx.NoContent(w)
}

// Me — GET /api/auth/me. Retorna o usuário autenticado.
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	uid := middleware.UserIDFrom(r.Context())
	u, err := h.users.FindByID(r.Context(), uid)
	if err != nil {
		// Token válido mas usuário sumiu (deletado): trata como não autenticado.
		writeDomainError(w, r, h.log, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto.NewUserResponse(u))
}
