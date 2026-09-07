package incident

import (
	"context"

	"incidenttrack/internal/domain/entity"
)

// Execute substitui os campos de um incidente de ownerID (semântica PUT).
// Carrega o registro antes para (a) confirmar existência e posse e
// (b) preservar created_at. Retorna errdomain.ErrNotFound /
// errdomain.ErrForbidden conforme o caso.
func (uc *UpdateIncident) Execute(ctx context.Context, id, ownerID string, in UpdateInput) (*entity.Incident, error) {
	current, err := uc.repo.GetByIDForOwner(ctx, id, ownerID)
	if err != nil {
		return nil, err
	}

	current.Titulo = in.Titulo
	current.Descricao = in.Descricao
	current.Categoria = entity.Categoria(in.Categoria)
	current.Severidade = entity.Severidade(in.Severidade)
	current.Status = entity.Status(in.Status)
	current.UpdatedAt = uc.now().UTC()

	if err := uc.repo.UpdateForOwner(ctx, current); err != nil {
		return nil, err
	}
	return current, nil
}
