// Package token gera e valida JWTs HS256 para a sessão do usuário
// (docs/plano.md §5). Implementação própria e mínima, sem dependência
// externa: assinatura HMAC-SHA256, comparação em tempo constante,
// validação estrita de alg/exp/nbf.
//
// Segurança (OWASP A07):
//   - alg "none" e algoritmos assimétricos são rejeitados na verificação.
//   - exp curto (config: 15 min) limita a janela de um token vazado.
//   - jti permite revogação antecipada via denylist em memória.
package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Claims são as reivindicações que emitimos e validamos.
type Claims struct {
	Subject   string // ID do usuário
	TokenID   string // jti
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// Service é a porta consumida pelo middleware de auth e pelos usecases.
type Service interface {
	Generate(userID string) (raw string, c Claims, err error)
	Verify(raw string) (Claims, error)
}

var (
	ErrMalformed = errors.New("token: formato inválido")
	ErrSignature = errors.New("token: assinatura inválida")
	ErrExpired   = errors.New("token: expirado")
	ErrNotYet    = errors.New("token: ainda não válido")
	ErrAlg       = errors.New("token: algoritmo não suportado")
)

type header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type payload struct {
	Sub string `json:"sub"`
	Jti string `json:"jti"`
	Iat int64  `json:"iat"`
	Nbf int64  `json:"nbf"`
	Exp int64  `json:"exp"`
}

// HS256 assina com uma chave HMAC simétrica.
type HS256 struct {
	secret []byte
	ttl    time.Duration
	now    func() time.Time // injeção para testes
}

// NewHS256 cria o serviço. Faz cópia defensiva do segredo.
func NewHS256(secret []byte, ttl time.Duration) *HS256 {
	cp := make([]byte, len(secret))
	copy(cp, secret)
	return &HS256{secret: cp, ttl: ttl, now: time.Now}
}

var b64 = base64.RawURLEncoding

// Generate emite um token assinado para userID.
func (s *HS256) Generate(userID string) (string, Claims, error) {
	now := s.now().UTC()
	c := Claims{
		Subject:   userID,
		TokenID:   uuid.NewString(),
		IssuedAt:  now,
		ExpiresAt: now.Add(s.ttl),
	}

	h, err := json.Marshal(header{Alg: "HS256", Typ: "JWT"})
	if err != nil {
		return "", Claims{}, err
	}
	p, err := json.Marshal(payload{
		Sub: c.Subject,
		Jti: c.TokenID,
		Iat: now.Unix(),
		Nbf: now.Unix(),
		Exp: c.ExpiresAt.Unix(),
	})
	if err != nil {
		return "", Claims{}, err
	}

	signingInput := b64.EncodeToString(h) + "." + b64.EncodeToString(p)
	sig := s.sign(signingInput)
	return signingInput + "." + sig, c, nil
}

// Verify valida assinatura e janela temporal, devolvendo os claims.
func (s *HS256) Verify(raw string) (Claims, error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return Claims{}, ErrMalformed
	}

	var h header
	hb, err := b64.DecodeString(parts[0])
	if err != nil || json.Unmarshal(hb, &h) != nil {
		return Claims{}, ErrMalformed
	}
	// Rejeita qualquer algoritmo diferente do esperado, incluindo "none".
	if !strings.EqualFold(h.Alg, "HS256") {
		return Claims{}, ErrAlg
	}

	expectedSig := s.sign(parts[0] + "." + parts[1])
	if !hmac.Equal([]byte(expectedSig), []byte(parts[2])) {
		return Claims{}, ErrSignature
	}

	var p payload
	pb, err := b64.DecodeString(parts[1])
	if err != nil || json.Unmarshal(pb, &p) != nil {
		return Claims{}, ErrMalformed
	}
	if p.Sub == "" || p.Exp == 0 {
		return Claims{}, ErrMalformed
	}

	now := s.now().UTC()
	const skew = 30 * time.Second
	if now.After(time.Unix(p.Exp, 0).Add(skew)) {
		return Claims{}, ErrExpired
	}
	if p.Nbf != 0 && now.Add(skew).Before(time.Unix(p.Nbf, 0)) {
		return Claims{}, ErrNotYet
	}

	return Claims{
		Subject:   p.Sub,
		TokenID:   p.Jti,
		IssuedAt:  time.Unix(p.Iat, 0).UTC(),
		ExpiresAt: time.Unix(p.Exp, 0).UTC(),
	}, nil
}

func (s *HS256) sign(input string) string {
	m := hmac.New(sha256.New, s.secret)
	m.Write([]byte(input))
	return b64.EncodeToString(m.Sum(nil))
}

// compile-time: garante que HS256 satisfaz Service.
var _ Service = (*HS256)(nil)

// String facilita logs sem vazar dados sensíveis.
func (c Claims) String() string {
	return fmt.Sprintf("Claims{sub:%s jti:%s exp:%s}", c.Subject, c.TokenID, c.ExpiresAt.Format(time.RFC3339))
}
