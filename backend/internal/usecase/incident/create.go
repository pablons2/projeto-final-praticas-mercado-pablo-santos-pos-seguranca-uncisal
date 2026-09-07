package incident

import (
	"context"

	"github.com/google/uuid"

	"incidenttrack/internal/domain/entity"
)

// Execute cria e persiste um novo incidente pertencente a ownerID.
func (uc *CreateIncident) Execute(ctx context.Context, ownerID string, in CreateInput) (*entity.Incident, error) {
	now := uc.now().UTC()
	inc := &entity.Incident{
		ID:         uuid.NewString(),
		Titulo:     in.Titulo,
		Descricao:  in.Descricao,
		Categoria:  entity.Categoria(in.Categoria),
		Severidade: entity.Severidade(in.Severidade),
		Status:     statusOrDefault(in.Status),
		OwnerID:    ownerID,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := uc.repo.Create(ctx, inc); err != nil {
		return nil, err
	}
	return inc, nil
}
