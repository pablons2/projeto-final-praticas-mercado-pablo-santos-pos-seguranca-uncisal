// Package incident contém os usecases de CRUD de incidentes. Cada
// operação é um tipo próprio com um único método Execute (SRP). Todos
// recebem ownerID e o propagam ao repositório, que filtra por dono na
// query — a autorização por recurso é verificada aqui e no SQL
// (defesa em profundidade; OWASP A01, docs/plano.md §7).
package incident

import (
	"time"

	"incidenttrack/internal/domain/entity"
	"incidenttrack/internal/domain/repository"
)

// deps são as dependências compartilhadas pelos usecases deste pacote.
type deps struct {
	repo repository.IncidentRepository
	now  func() time.Time
}

func newDeps(repo repository.IncidentRepository) deps {
	return deps{repo: repo, now: time.Now}
}

// CreateIncident cria um incidente para o dono informado.
type CreateIncident struct{ deps }

// ListIncidents lista os incidentes do dono, com filtros opcionais.
type ListIncidents struct{ deps }

// GetIncident detalha um incidente do dono.
type GetIncident struct{ deps }

// UpdateIncident substitui os campos de um incidente do dono.
type UpdateIncident struct{ deps }

// DeleteIncident remove um incidente do dono.
type DeleteIncident struct{ deps }

func NewCreateIncident(r repository.IncidentRepository) *CreateIncident {
	return &CreateIncident{newDeps(r)}
}
func NewListIncidents(r repository.IncidentRepository) *ListIncidents {
	return &ListIncidents{newDeps(r)}
}
func NewGetIncident(r repository.IncidentRepository) *GetIncident { return &GetIncident{newDeps(r)} }
func NewUpdateIncident(r repository.IncidentRepository) *UpdateIncident {
	return &UpdateIncident{newDeps(r)}
}
func NewDeleteIncident(r repository.IncidentRepository) *DeleteIncident {
	return &DeleteIncident{newDeps(r)}
}

// CreateInput / UpdateInput são os dados já normalizados e validados
// pela camada dto (enums garantidamente dentro do conjunto fechado).
type CreateInput struct {
	Titulo     string
	Descricao  string
	Categoria  string
	Severidade string
	Status     string
}

type UpdateInput struct {
	Titulo     string
	Descricao  string
	Categoria  string
	Severidade string
	Status     string
}

func statusOrDefault(s string) entity.Status {
	if s == "" {
		return entity.StatusAberto
	}
	return entity.Status(s)
}
