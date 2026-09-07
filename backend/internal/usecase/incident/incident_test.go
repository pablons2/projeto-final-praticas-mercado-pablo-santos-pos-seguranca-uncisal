package incident_test

import (
	"context"
	"errors"
	"testing"

	"incidenttrack/internal/domain/entity"
	"incidenttrack/internal/domain/errdomain"
	"incidenttrack/internal/domain/repository"
	"incidenttrack/internal/repository/memory"
	incuc "incidenttrack/internal/usecase/incident"
)

const (
	ana   = "user-ana"
	bruno = "user-bruno"
)

func validCreate() incuc.CreateInput {
	return incuc.CreateInput{
		Titulo:     "Tentativa de phishing no RH",
		Descricao:  "E-mail falso pedindo troca de senha.",
		Categoria:  string(entity.CategoriaPhishing),
		Severidade: string(entity.SeveridadeAlta),
		Status:     string(entity.StatusAberto),
	}
}

func TestCreateAndList(t *testing.T) {
	repo := memory.NewIncidentRepo()
	create := incuc.NewCreateIncident(repo)
	list := incuc.NewListIncidents(repo)
	ctx := context.Background()

	inc, err := create.Execute(ctx, ana, validCreate())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if inc.OwnerID != ana {
		t.Errorf("owner errado: %q", inc.OwnerID)
	}

	// Incidente de outro usuário não deve aparecer na lista da Ana.
	if _, err := create.Execute(ctx, bruno, validCreate()); err != nil {
		t.Fatalf("create bruno: %v", err)
	}

	got, err := list.Execute(ctx, ana, repository.IncidentFilter{})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("esperava 1 incidente da Ana, veio %d", len(got))
	}
}

func TestGet_OtherUsersIncident_IsForbidden(t *testing.T) {
	repo := memory.NewIncidentRepo()
	create := incuc.NewCreateIncident(repo)
	get := incuc.NewGetIncident(repo)
	ctx := context.Background()

	inc, _ := create.Execute(ctx, ana, validCreate())

	_, err := get.Execute(ctx, inc.ID, bruno)
	if !errors.Is(err, errdomain.ErrForbidden) {
		t.Fatalf("esperava ErrForbidden ao acessar incidente alheio, veio %v", err)
	}
}

func TestUpdate_OnlyOwner(t *testing.T) {
	repo := memory.NewIncidentRepo()
	create := incuc.NewCreateIncident(repo)
	update := incuc.NewUpdateIncident(repo)
	ctx := context.Background()

	inc, _ := create.Execute(ctx, ana, validCreate())

	upd := incuc.UpdateInput{
		Titulo:     "Título atualizado pela Ana",
		Descricao:  "novo texto",
		Categoria:  string(entity.CategoriaMalware),
		Severidade: string(entity.SeveridadeCritica),
		Status:     string(entity.StatusEmAnalise),
	}

	// Bruno não pode atualizar.
	if _, err := update.Execute(ctx, inc.ID, bruno, upd); !errors.Is(err, errdomain.ErrForbidden) {
		t.Fatalf("esperava ErrForbidden para não-dono, veio %v", err)
	}

	// Ana pode.
	got, err := update.Execute(ctx, inc.ID, ana, upd)
	if err != nil {
		t.Fatalf("update dono: %v", err)
	}
	if got.Titulo != upd.Titulo || got.Severidade != entity.SeveridadeCritica {
		t.Errorf("update não aplicado: %+v", got)
	}
	if !got.UpdatedAt.After(inc.UpdatedAt) && !got.UpdatedAt.Equal(inc.UpdatedAt) {
		t.Error("updated_at deveria avançar")
	}
	if got.CreatedAt != inc.CreatedAt {
		t.Error("created_at não pode mudar no update")
	}
}

func TestDelete_OnlyOwner_ThenNotFound(t *testing.T) {
	repo := memory.NewIncidentRepo()
	create := incuc.NewCreateIncident(repo)
	del := incuc.NewDeleteIncident(repo)
	get := incuc.NewGetIncident(repo)
	ctx := context.Background()

	inc, _ := create.Execute(ctx, ana, validCreate())

	if err := del.Execute(ctx, inc.ID, bruno); !errors.Is(err, errdomain.ErrForbidden) {
		t.Fatalf("esperava ErrForbidden para não-dono, veio %v", err)
	}
	if err := del.Execute(ctx, inc.ID, ana); err != nil {
		t.Fatalf("delete dono: %v", err)
	}
	if _, err := get.Execute(ctx, inc.ID, ana); !errors.Is(err, errdomain.ErrNotFound) {
		t.Fatalf("esperava ErrNotFound após exclusão, veio %v", err)
	}
}

func TestDelete_Missing_IsNotFound(t *testing.T) {
	repo := memory.NewIncidentRepo()
	del := incuc.NewDeleteIncident(repo)
	if err := del.Execute(context.Background(), "nao-existe", ana); !errors.Is(err, errdomain.ErrNotFound) {
		t.Fatalf("esperava ErrNotFound, veio %v", err)
	}
}
