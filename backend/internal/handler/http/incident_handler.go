package httpapi

import (
	"log/slog"
	"net/http"

	"incidenttrack/internal/dto"
	"incidenttrack/internal/middleware"
	"incidenttrack/internal/pkg/httpx"
	incuc "incidenttrack/internal/usecase/incident"
)

// IncidentHandler expõe o CRUD de incidentes. Toda rota abaixo já passou
// pelo middleware de auth, então middleware.UserIDFrom sempre retorna um
// ID não vazio.
type IncidentHandler struct {
	create *incuc.CreateIncident
	list   *incuc.ListIncidents
	get    *incuc.GetIncident
	update *incuc.UpdateIncident
	delete *incuc.DeleteIncident
	log    *slog.Logger
}

// IncidentHandlerDeps agrupa as dependências.
type IncidentHandlerDeps struct {
	Create *incuc.CreateIncident
	List   *incuc.ListIncidents
	Get    *incuc.GetIncident
	Update *incuc.UpdateIncident
	Delete *incuc.DeleteIncident
	Log    *slog.Logger
}

// NewIncidentHandler monta o handler.
func NewIncidentHandler(d IncidentHandlerDeps) *IncidentHandler {
	return &IncidentHandler{
		create: d.Create,
		list:   d.List,
		get:    d.Get,
		update: d.Update,
		delete: d.Delete,
		log:    d.Log,
	}
}

// List — GET /api/incidents. Filtros opcionais via query string.
func (h *IncidentHandler) List(w http.ResponseWriter, r *http.Request) {
	owner := middleware.UserIDFrom(r.Context())
	q := r.URL.Query()
	filter := dto.ListIncidentsQuery{
		Status:     q.Get("status"),
		Severidade: q.Get("severidade"),
		Categoria:  q.Get("categoria"),
		Query:      q.Get("q"),
	}.ToFilter()

	items, err := h.list.Execute(r.Context(), owner, filter)
	if err != nil {
		writeDomainError(w, r, h.log, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto.NewIncidentListResponse(items))
}

// Get — GET /api/incidents/{id}.
func (h *IncidentHandler) Get(w http.ResponseWriter, r *http.Request) {
	owner := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")

	inc, err := h.get.Execute(r.Context(), id, owner)
	if err != nil {
		writeDomainError(w, r, h.log, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto.NewIncidentResponse(inc))
}

// Create — POST /api/incidents.
func (h *IncidentHandler) Create(w http.ResponseWriter, r *http.Request) {
	owner := middleware.UserIDFrom(r.Context())

	var in dto.CreateIncidentInput
	if !decodeJSON(w, r, &in) {
		return
	}
	in.Normalize()
	if errs := in.Validate(); !errs.Ok() {
		httpx.WriteValidationError(w, "Dados inválidos.", errs.Fields())
		return
	}

	inc, err := h.create.Execute(r.Context(), owner, incuc.CreateInput{
		Titulo:     in.Titulo,
		Descricao:  in.Descricao,
		Categoria:  in.Categoria,
		Severidade: in.Severidade,
		Status:     in.Status,
	})
	if err != nil {
		writeDomainError(w, r, h.log, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, dto.NewIncidentResponse(inc))
}

// Update — PUT /api/incidents/{id} (substituição total).
func (h *IncidentHandler) Update(w http.ResponseWriter, r *http.Request) {
	owner := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")

	var in dto.UpdateIncidentInput
	if !decodeJSON(w, r, &in) {
		return
	}
	in.Normalize()
	if errs := in.Validate(); !errs.Ok() {
		httpx.WriteValidationError(w, "Dados inválidos.", errs.Fields())
		return
	}

	inc, err := h.update.Execute(r.Context(), id, owner, incuc.UpdateInput{
		Titulo:     in.Titulo,
		Descricao:  in.Descricao,
		Categoria:  in.Categoria,
		Severidade: in.Severidade,
		Status:     in.Status,
	})
	if err != nil {
		writeDomainError(w, r, h.log, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, dto.NewIncidentResponse(inc))
}

// Delete — DELETE /api/incidents/{id}.
func (h *IncidentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	owner := middleware.UserIDFrom(r.Context())
	id := r.PathValue("id")

	if err := h.delete.Execute(r.Context(), id, owner); err != nil {
		writeDomainError(w, r, h.log, err)
		return
	}
	httpx.NoContent(w)
}
