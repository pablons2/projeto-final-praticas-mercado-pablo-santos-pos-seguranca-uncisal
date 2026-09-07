package incident

import "context"

// Execute remove um incidente de ownerID. Retorna errdomain.ErrNotFound
// se não existe e errdomain.ErrForbidden se é de outro usuário.
func (uc *DeleteIncident) Execute(ctx context.Context, id, ownerID string) error {
	// Confirma existência e posse antes de apagar, para diferenciar
	// 404 de 403 na resposta (o DELETE cru não distingue os casos).
	if _, err := uc.repo.GetByIDForOwner(ctx, id, ownerID); err != nil {
		return err
	}
	return uc.repo.DeleteForOwner(ctx, id, ownerID)
}
