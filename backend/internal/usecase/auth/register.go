// Package auth contém os usecases de autenticação. Cada arquivo resolve
// uma única ação de negócio (SRP; docs/plano.md §4). Os usecases
// dependem apenas de abstrações de domain e pkg — nunca de implementação
// concreta de banco ou HTTP (DIP).
package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"incidenttrack/internal/domain/entity"
	"incidenttrack/internal/domain/errdomain"
	"incidenttrack/internal/domain/repository"
	"incidenttrack/internal/pkg/hash"
)

// RegisterInput são os dados já normalizados e validados pela camada dto.
type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

// RegisterUser cria um novo usuário com senha hasheada.
type RegisterUser struct {
	users  repository.UserRepository
	hasher hash.Hasher
	now    func() time.Time
}

// NewRegisterUser injeta as dependências.
func NewRegisterUser(users repository.UserRepository, hasher hash.Hasher) *RegisterUser {
	return &RegisterUser{users: users, hasher: hasher, now: time.Now}
}

// Execute persiste o usuário. Retorna errdomain.ErrEmailTaken se o
// e-mail já existir.
func (uc *RegisterUser) Execute(ctx context.Context, in RegisterInput) (*entity.User, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))

	// Checagem prévia amigável; a garantia real é o índice UNIQUE no banco,
	// tratado abaixo para evitar TOCTOU.
	if _, err := uc.users.FindByEmail(ctx, email); err == nil {
		return nil, errdomain.ErrEmailTaken
	} else if !errors.Is(err, errdomain.ErrNotFound) {
		return nil, err
	}

	hashed, err := uc.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	now := uc.now().UTC()
	u := &entity.User{
		ID:           uuid.NewString(),
		Email:        email,
		Name:         in.Name,
		PasswordHash: hashed,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := uc.users.Create(ctx, u); err != nil {
		return nil, err // pode ser errdomain.ErrEmailTaken vindo do UNIQUE
	}
	return u, nil
}
