// Package memory fornece implementações em memória das portas de
// domain/repository. Serve aos testes unitários de usecase (docs/plano.md
// §9): como os usecases dependem da interface, e não do SQLite, a regra
// de negócio é testável sem subir banco (LSP na prática).
package memory

import (
	"context"
	"sort"
	"strings"
	"sync"

	"incidenttrack/internal/domain/entity"
	"incidenttrack/internal/domain/errdomain"
	"incidenttrack/internal/domain/repository"
)

// UserRepo é um repositório de usuários em memória.
type UserRepo struct {
	mu    sync.RWMutex
	byID  map[string]entity.User
	byEml map[string]string // email -> id
}

// NewUserRepo cria o fake vazio.
func NewUserRepo() *UserRepo {
	return &UserRepo{byID: map[string]entity.User{}, byEml: map[string]string{}}
}

var _ repository.UserRepository = (*UserRepo)(nil)

func (r *UserRepo) Create(_ context.Context, u *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := strings.ToLower(u.Email)
	if _, ok := r.byEml[key]; ok {
		return errdomain.ErrEmailTaken
	}
	r.byID[u.ID] = *u
	r.byEml[key] = u.ID
	return nil
}

func (r *UserRepo) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	id, ok := r.byEml[strings.ToLower(email)]
	if !ok {
		return nil, errdomain.ErrNotFound
	}
	u := r.byID[id]
	return &u, nil
}

func (r *UserRepo) FindByID(_ context.Context, id string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, errdomain.ErrNotFound
	}
	return &u, nil
}

// IncidentRepo é um repositório de incidentes em memória.
type IncidentRepo struct {
	mu   sync.RWMutex
	data map[string]entity.Incident
}

// NewIncidentRepo cria o fake vazio.
func NewIncidentRepo() *IncidentRepo {
	return &IncidentRepo{data: map[string]entity.Incident{}}
}

var _ repository.IncidentRepository = (*IncidentRepo)(nil)

func (r *IncidentRepo) Create(_ context.Context, in *entity.Incident) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[in.ID] = *in
	return nil
}

func (r *IncidentRepo) ListByOwner(_ context.Context, ownerID string, f repository.IncidentFilter) ([]entity.Incident, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var out []entity.Incident
	for _, inc := range r.data {
		if inc.OwnerID != ownerID {
			continue
		}
		if f.Status != "" && string(inc.Status) != f.Status {
			continue
		}
		if f.Severidade != "" && string(inc.Severidade) != f.Severidade {
			continue
		}
		if f.Categoria != "" && string(inc.Categoria) != f.Categoria {
			continue
		}
		if f.Query != "" {
			q := strings.ToLower(f.Query)
			if !strings.Contains(strings.ToLower(inc.Titulo), q) &&
				!strings.Contains(strings.ToLower(inc.Descricao), q) {
				continue
			}
		}
		out = append(out, inc)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.After(out[j].CreatedAt) })
	if f.Limit > 0 && len(out) > f.Limit {
		out = out[:f.Limit]
	}
	return out, nil
}

func (r *IncidentRepo) GetByIDForOwner(_ context.Context, id, ownerID string) (*entity.Incident, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	inc, ok := r.data[id]
	if !ok {
		return nil, errdomain.ErrNotFound
	}
	if inc.OwnerID != ownerID {
		return nil, errdomain.ErrForbidden
	}
	return &inc, nil
}

func (r *IncidentRepo) UpdateForOwner(_ context.Context, in *entity.Incident) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cur, ok := r.data[in.ID]
	if !ok {
		return errdomain.ErrNotFound
	}
	if cur.OwnerID != in.OwnerID {
		return errdomain.ErrForbidden
	}
	r.data[in.ID] = *in
	return nil
}

func (r *IncidentRepo) DeleteForOwner(_ context.Context, id, ownerID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	inc, ok := r.data[id]
	if !ok {
		return errdomain.ErrNotFound
	}
	if inc.OwnerID != ownerID {
		return errdomain.ErrForbidden
	}
	delete(r.data, id)
	return nil
}
