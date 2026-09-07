// Package repository declara as portas de persistência do domínio. São
// interfaces pequenas e específicas (OWASP-agnóstico; princípio ISP do
// docs/plano.md §4). As implementações concretas vivem em
// internal/repository/sqlite e são injetadas pelo main.go (DIP).
package repository

import (
	"context"

	"incidenttrack/internal/domain/entity"
)

// UserRepository persiste usuários.
type UserRepository interface {
	// Create insere um novo usuário. Deve retornar errdomain.ErrEmailTaken
	// se o e-mail já existir (violação de UNIQUE).
	Create(ctx context.Context, u *entity.User) error
	// FindByEmail retorna errdomain.ErrNotFound quando não há usuário.
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	// FindByID retorna errdomain.ErrNotFound quando não há usuário.
	FindByID(ctx context.Context, id string) (*entity.User, error)
}

// IncidentFilter são filtros opcionais de listagem. Campos zero são ignorados.
type IncidentFilter struct {
	Status     string
	Severidade string
	Categoria  string
	Query      string // busca textual em título/descrição
	Limit      int
}

// IncidentRepository persiste incidentes. Todo método recebe ownerID e
// filtra por ele na própria query — a checagem de dono não é opcional
// nem delegada à camada de cima (defesa em profundidade, OWASP A01).
type IncidentRepository interface {
	Create(ctx context.Context, in *entity.Incident) error
	// ListByOwner devolve os incidentes do dono, aplicando f.
	ListByOwner(ctx context.Context, ownerID string, f IncidentFilter) ([]entity.Incident, error)
	// GetByIDForOwner retorna errdomain.ErrNotFound se o id não existe e
	// errdomain.ErrForbidden se existe mas é de outro dono.
	GetByIDForOwner(ctx context.Context, id, ownerID string) (*entity.Incident, error)
	// UpdateForOwner aplica as mudanças somente se in.OwnerID for o dono.
	UpdateForOwner(ctx context.Context, in *entity.Incident) error
	// DeleteForOwner remove somente se pertencer ao dono.
	DeleteForOwner(ctx context.Context, id, ownerID string) error
}
