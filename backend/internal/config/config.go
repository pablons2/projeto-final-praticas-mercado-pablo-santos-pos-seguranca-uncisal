// Package config carrega a configuração da aplicação exclusivamente a
// partir de variáveis de ambiente. Nenhum segredo é embutido no código
// (OWASP A02 — Security Misconfiguration; ver docs/plano.md §7).
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config agrega todos os parâmetros de execução do servidor.
type Config struct {
	// Env indica o ambiente ("development" ou "production"). Em produção
	// o servidor exige cookies Secure e recusa segredos fracos.
	Env string

	// HTTPAddr é o endereço de escuta (ex.: ":8080").
	HTTPAddr string

	// DatabaseURL é o caminho do arquivo SQLite (ex.: "file:data/incidenttrack.db").
	DatabaseURL string

	// JWTSecret é a chave HMAC (mínimo 32 bytes) usada para assinar tokens.
	JWTSecret []byte

	// AccessTokenTTL é a validade curta do token de sessão.
	AccessTokenTTL time.Duration

	// CookieName é o nome do cookie de sessão.
	CookieName string

	// CookieSecure controla o atributo Secure do cookie. Deve ser true em
	// produção; pode ser false em desenvolvimento local sobre HTTP.
	CookieSecure bool

	// CookieDomain é opcional; vazio deixa o browser inferir o host.
	CookieDomain string

	// CORSAllowedOrigin é a origem exata do frontend autorizada a chamar a API.
	CORSAllowedOrigin string

	// TrustProxy habilita a leitura de X-Forwarded-For / X-Real-IP para
	// identificar o IP do cliente (rate limiting). Só deve ser true quando
	// a API está atrás de um proxy reverso confiável (Nginx no deploy, ou
	// o serviço "web" no docker-compose). Default: true em produção.
	TrustProxy bool

	// LoginRateLimit / LoginRateWindow: tentativas de login por IP.
	LoginRateLimit  int
	LoginRateWindow time.Duration

	// APIRateLimit / APIRateWindow: teto global de requisições por IP.
	APIRateLimit  int
	APIRateWindow time.Duration

	// BcryptCost é o custo do hash de senha (>= 10).
	BcryptCost int

	// ReadTimeout / WriteTimeout / IdleTimeout do http.Server.
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration

	// ShutdownTimeout é o prazo para drenar conexões no encerramento gracioso.
	ShutdownTimeout time.Duration
}

// IsProduction indica se o ambiente é de produção.
func (c *Config) IsProduction() bool { return c.Env == "production" }

// Load lê a configuração do ambiente e valida os campos obrigatórios.
func Load() (*Config, error) {
	c := &Config{
		Env:               getStr("APP_ENV", "development"),
		HTTPAddr:          getStr("HTTP_ADDR", ":8080"),
		DatabaseURL:       getStr("DATABASE_URL", "file:data/incidenttrack.db"),
		AccessTokenTTL:    getDur("ACCESS_TOKEN_TTL", 15*time.Minute),
		CookieName:        getStr("COOKIE_NAME", "it_session"),
		CookieSecure:      getBool("COOKIE_SECURE", true),
		CookieDomain:      getStr("COOKIE_DOMAIN", ""),
		CORSAllowedOrigin: getStr("CORS_ALLOWED_ORIGIN", "http://localhost:5173"),
		LoginRateLimit:    getInt("LOGIN_RATE_LIMIT", 5),
		LoginRateWindow:   getDur("LOGIN_RATE_WINDOW", 15*time.Minute),
		APIRateLimit:      getInt("API_RATE_LIMIT", 120),
		APIRateWindow:     getDur("API_RATE_WINDOW", time.Minute),
		BcryptCost:        getInt("BCRYPT_COST", 12),
		ReadTimeout:       getDur("HTTP_READ_TIMEOUT", 10*time.Second),
		WriteTimeout:      getDur("HTTP_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:       getDur("HTTP_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout:   getDur("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
	}
	c.TrustProxy = getBool("TRUST_PROXY", c.IsProduction())

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("config: JWT_SECRET é obrigatório")
	}
	if len(secret) < 32 {
		return nil, fmt.Errorf("config: JWT_SECRET deve ter >= 32 bytes (tem %d)", len(secret))
	}
	c.JWTSecret = []byte(secret)

	if c.BcryptCost < 10 || c.BcryptCost > 15 {
		return nil, fmt.Errorf("config: BCRYPT_COST fora da faixa 10..15 (%d)", c.BcryptCost)
	}
	if c.CORSAllowedOrigin == "" || c.CORSAllowedOrigin == "*" {
		return nil, errors.New("config: CORS_ALLOWED_ORIGIN deve ser uma origem exata, nunca vazia ou '*'")
	}
	if c.IsProduction() && !c.CookieSecure {
		return nil, errors.New("config: COOKIE_SECURE deve ser true em produção")
	}

	return c, nil
}

func getStr(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

func getInt(key string, def int) int {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func getBool(key string, def bool) bool {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}

func getDur(key string, def time.Duration) time.Duration {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return def
}
