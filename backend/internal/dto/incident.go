package dto

import (
	"time"

	"incidenttrack/internal/domain/entity"
	"incidenttrack/internal/domain/repository"
	"incidenttrack/internal/pkg/validator"
)

const (
	tituloMin    = 3
	tituloMax    = 140
	descricaoMax = 5000
)

// CreateIncidentInput é o corpo de POST /api/incidents.
type CreateIncidentInput struct {
	Titulo     string `json:"titulo"`
	Descricao  string `json:"descricao"`
	Categoria  string `json:"categoria"`
	Severidade string `json:"severidade"`
	Status     string `json:"status"` // opcional; default "aberto"
}

// Normalize sanitiza os campos textuais.
func (in *CreateIncidentInput) Normalize() {
	in.Titulo = validator.CollapseSpaces(in.Titulo)
	in.Descricao = validator.StripControl(in.Descricao)
	in.Categoria = validator.Trim(in.Categoria)
	in.Severidade = validator.Trim(in.Severidade)
	in.Status = validator.Trim(in.Status)
	if in.Status == "" {
		in.Status = string(entity.StatusAberto)
	}
}

// Validate confere tamanho e pertencimento aos enums fechados.
func (in *CreateIncidentInput) Validate() *validator.Errors {
	e := validator.New()
	validateIncidentCommon(e, in.Titulo, in.Descricao, in.Categoria, in.Severidade, in.Status)
	return e
}

// UpdateIncidentInput é o corpo de PUT /api/incidents/{id} (substituição total).
type UpdateIncidentInput struct {
	Titulo     string `json:"titulo"`
	Descricao  string `json:"descricao"`
	Categoria  string `json:"categoria"`
	Severidade string `json:"severidade"`
	Status     string `json:"status"`
}

// Normalize sanitiza os campos textuais.
func (in *UpdateIncidentInput) Normalize() {
	in.Titulo = validator.CollapseSpaces(in.Titulo)
	in.Descricao = validator.StripControl(in.Descricao)
	in.Categoria = validator.Trim(in.Categoria)
	in.Severidade = validator.Trim(in.Severidade)
	in.Status = validator.Trim(in.Status)
}

// Validate exige todos os campos (PUT é substituição completa).
func (in *UpdateIncidentInput) Validate() *validator.Errors {
	e := validator.New()
	validateIncidentCommon(e, in.Titulo, in.Descricao, in.Categoria, in.Severidade, in.Status)
	if in.Status == "" {
		e.Add("status", "Campo obrigatório.")
	}
	return e
}

func validateIncidentCommon(e *validator.Errors, titulo, descricao, categoria, severidade, status string) {
	if !validator.RuneLenBetween(titulo, tituloMin, tituloMax) {
		e.Add("titulo", "Deve ter entre 3 e 140 caracteres.")
	}
	if !validator.RuneLenBetween(descricao, 0, descricaoMax) {
		e.Add("descricao", "Deve ter no máximo 5000 caracteres.")
	}
	if !entity.Categoria(categoria).Valida() {
		e.Add("categoria", "Categoria inválida.")
	}
	if !entity.Severidade(severidade).Valida() {
		e.Add("severidade", "Severidade inválida.")
	}
	if status != "" && !entity.Status(status).Valida() {
		e.Add("status", "Status inválido.")
	}
}

// ListIncidentsQuery são os filtros opcionais de GET /api/incidents.
type ListIncidentsQuery struct {
	Status     string
	Severidade string
	Categoria  string
	Query      string
}

// ToFilter converte a query string validada em filtro de repositório.
// Valores fora dos enums são descartados silenciosamente (não é erro
// filtrar por algo inexistente — apenas não filtra).
func (q ListIncidentsQuery) ToFilter() repository.IncidentFilter {
	f := repository.IncidentFilter{Limit: 500}
	if entity.Status(q.Status).Valida() {
		f.Status = q.Status
	}
	if entity.Severidade(q.Severidade).Valida() {
		f.Severidade = q.Severidade
	}
	if entity.Categoria(q.Categoria).Valida() {
		f.Categoria = q.Categoria
	}
	if s := validator.CollapseSpaces(q.Query); s != "" && validator.RuneLenBetween(s, 1, 140) {
		f.Query = s
	}
	return f
}

// IncidentResponse é a projeção de saída de um incidente.
type IncidentResponse struct {
	ID         string    `json:"id"`
	Titulo     string    `json:"titulo"`
	Descricao  string    `json:"descricao"`
	Categoria  string    `json:"categoria"`
	Severidade string    `json:"severidade"`
	Status     string    `json:"status"`
	OwnerID    string    `json:"owner_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// NewIncidentResponse converte a entidade em resposta.
func NewIncidentResponse(in *entity.Incident) IncidentResponse {
	return IncidentResponse{
		ID:         in.ID,
		Titulo:     in.Titulo,
		Descricao:  in.Descricao,
		Categoria:  string(in.Categoria),
		Severidade: string(in.Severidade),
		Status:     string(in.Status),
		OwnerID:    in.OwnerID,
		CreatedAt:  in.CreatedAt,
		UpdatedAt:  in.UpdatedAt,
	}
}

// NewIncidentListResponse converte uma lista.
func NewIncidentListResponse(items []entity.Incident) []IncidentResponse {
	out := make([]IncidentResponse, 0, len(items))
	for i := range items {
		out = append(out, NewIncidentResponse(&items[i]))
	}
	return out
}
