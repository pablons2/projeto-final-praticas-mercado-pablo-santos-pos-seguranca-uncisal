package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"strings"

	"incidenttrack/internal/domain/entity"
	"incidenttrack/internal/domain/errdomain"
	"incidenttrack/internal/domain/repository"
	"incidenttrack/internal/pkg/hash"
)

// LoginInput são os dados já normalizados pela camada dto.
type LoginInput struct {
	Email    string
	Password string
}

// LoginUser autentica um usuário por e-mail e senha.
type LoginUser struct {
	users  repository.UserRepository
	hasher hash.Hasher
	// dummyHash é um hash bcrypt válido de valor fixo. Comparamos contra
	// ele quando o e-mail não existe, para que o tempo de resposta de
	// "usuário inexistente" seja próximo do de "senha errada"
	// (mitiga user enumeration por timing; OWASP A07).
	dummyHash string
}

// NewLoginUser injeta as dependências. dummyHash deve ser um hash bcrypt
// pré-computado (gerado uma vez na inicialização).
func NewLoginUser(users repository.UserRepository, hasher hash.Hasher, dummyHash string) *LoginUser {
	return &LoginUser{users: users, hasher: hasher, dummyHash: dummyHash}
}

// Execute retorna o usuário autenticado ou errdomain.ErrInvalidCredentials.
// A mensagem é sempre a mesma, não distinguindo e-mail de senha.
func (uc *LoginUser) Execute(ctx context.Context, in LoginInput) (*entity.User, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	u, err := uc.users.FindByEmail(ctx, email)
	if errors.Is(err, errdomain.ErrNotFound) {
		// Gasta tempo comparável ao caso real e descarta o resultado.
		_ = uc.hasher.Compare(uc.dummyHash, in.Password)
		return nil, errdomain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if cmpErr := uc.hasher.Compare(u.PasswordHash, in.Password); cmpErr != nil {
		return nil, errdomain.ErrInvalidCredentials
	}

	// Barreira extra irrelevante para segurança mas explícita: garante que
	// não retornamos usuário se o e-mail divergir por algum bug de camada.
	if subtle.ConstantTimeCompare([]byte(u.Email), []byte(email)) != 1 {
		return nil, errdomain.ErrInvalidCredentials
	}
	return u, nil
}
