package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"incidenttrack/internal/domain/entity"
	"incidenttrack/internal/domain/errdomain"
	"incidenttrack/internal/domain/repository"
)

// IncidentRepo implementa repository.IncidentRepository sobre SQLite.
type IncidentRepo struct{ db *DB }

// NewIncidentRepo cria o repositório de incidentes.
func NewIncidentRepo(db *DB) *IncidentRepo { return &IncidentRepo{db: db} }

var _ repository.IncidentRepository = (*IncidentRepo)(nil)

const colsIncident = "id, titulo, descricao, categoria, severidade, status, owner_id, created_at, updated_at"

// Create insere um novo incidente.
func (r *IncidentRepo) Create(ctx context.Context, in *entity.Incident) error {
	const q = `INSERT INTO incidents (` + colsIncident + `)
	           VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, q,
		in.ID, in.Titulo, in.Descricao, string(in.Categoria), string(in.Severidade),
		string(in.Status), in.OwnerID, rfc3339(in.CreatedAt), rfc3339(in.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("sqlite: criar incidente: %w", err)
	}
	return nil
}

// ListByOwner lista os incidentes do dono. O WHERE é montado dinamicamente
// APENAS com fragmentos fixos e placeholders `?`; nenhum valor de usuário
// entra na string SQL (OWASP A03/A05).
func (r *IncidentRepo) ListByOwner(ctx context.Context, ownerID string, f repository.IncidentFilter) ([]entity.Incident, error) {
	where := []string{"owner_id = ?"}
	args := []any{ownerID}

	if f.Status != "" {
		where = append(where, "status = ?")
		args = append(args, f.Status)
	}
	if f.Severidade != "" {
		where = append(where, "severidade = ?")
		args = append(args, f.Severidade)
	}
	if f.Categoria != "" {
		where = append(where, "categoria = ?")
		args = append(args, f.Categoria)
	}
	if f.Query != "" {
		where = append(where, "(titulo LIKE ? ESCAPE '\\' OR descricao LIKE ? ESCAPE '\\')")
		like := "%" + escapeLike(f.Query) + "%"
		args = append(args, like, like)
	}

	limit := f.Limit
	if limit <= 0 || limit > 500 {
		limit = 500
	}

	q := `SELECT ` + colsIncident + ` FROM incidents WHERE ` +
		strings.Join(where, " AND ") +
		` ORDER BY created_at DESC LIMIT ?`
	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("sqlite: listar incidentes: %w", err)
	}
	defer rows.Close()

	var out []entity.Incident
	for rows.Next() {
		inc, err := scanIncident(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *inc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("sqlite: iterar incidentes: %w", err)
	}
	return out, nil
}

// GetByIDForOwner busca um incidente e distingue inexistência de posse alheia.
func (r *IncidentRepo) GetByIDForOwner(ctx context.Context, id, ownerID string) (*entity.Incident, error) {
	const q = `SELECT ` + colsIncident + ` FROM incidents WHERE id = ? LIMIT 1`
	inc, err := scanIncidentRow(r.db.QueryRowContext(ctx, q, id))
	if err != nil {
		return nil, err
	}
	if inc.OwnerID != ownerID {
		// Existe, mas não é do solicitante (OWASP A01).
		return nil, errdomain.ErrForbidden
	}
	return inc, nil
}

// UpdateForOwner aplica a atualização somente na linha do dono. Se
// nenhuma linha for afetada, decide entre 404 e 403.
func (r *IncidentRepo) UpdateForOwner(ctx context.Context, in *entity.Incident) error {
	const q = `UPDATE incidents
	           SET titulo = ?, descricao = ?, categoria = ?, severidade = ?, status = ?, updated_at = ?
	           WHERE id = ? AND owner_id = ?`
	res, err := r.db.ExecContext(ctx, q,
		in.Titulo, in.Descricao, string(in.Categoria), string(in.Severidade),
		string(in.Status), rfc3339(in.UpdatedAt), in.ID, in.OwnerID,
	)
	if err != nil {
		return fmt.Errorf("sqlite: atualizar incidente: %w", err)
	}
	return r.assertAffected(ctx, res, in.ID, in.OwnerID)
}

// DeleteForOwner remove somente a linha do dono.
func (r *IncidentRepo) DeleteForOwner(ctx context.Context, id, ownerID string) error {
	const q = `DELETE FROM incidents WHERE id = ? AND owner_id = ?`
	res, err := r.db.ExecContext(ctx, q, id, ownerID)
	if err != nil {
		return fmt.Errorf("sqlite: remover incidente: %w", err)
	}
	return r.assertAffected(ctx, res, id, ownerID)
}

// assertAffected traduz "0 linhas afetadas" em ErrNotFound ou ErrForbidden.
func (r *IncidentRepo) assertAffected(ctx context.Context, res sql.Result, id, ownerID string) error {
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("sqlite: rows affected: %w", err)
	}
	if n > 0 {
		return nil
	}
	// Nada mudou: ou o id não existe, ou pertence a outro dono.
	const q = `SELECT owner_id FROM incidents WHERE id = ? LIMIT 1`
	var actualOwner string
	switch err := r.db.QueryRowContext(ctx, q, id).Scan(&actualOwner); {
	case isNoRows(err):
		return errdomain.ErrNotFound
	case err != nil:
		return fmt.Errorf("sqlite: verificar posse: %w", err)
	default:
		if actualOwner != ownerID {
			return errdomain.ErrForbidden
		}
		// id existe e é do dono mas UPDATE não alterou nada (valores iguais).
		return nil
	}
}

type scanner interface{ Scan(dest ...any) error }

func scanIncident(s scanner) (*entity.Incident, error) {
	var (
		inc                  entity.Incident
		cat, sev, st         string
		createdAt, updatedAt string
	)
	err := s.Scan(&inc.ID, &inc.Titulo, &inc.Descricao, &cat, &sev, &st,
		&inc.OwnerID, &createdAt, &updatedAt)
	if isNoRows(err) {
		return nil, errdomain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("sqlite: ler incidente: %w", err)
	}
	inc.Categoria = entity.Categoria(cat)
	inc.Severidade = entity.Severidade(sev)
	inc.Status = entity.Status(st)
	inc.CreatedAt = parseTime(createdAt)
	inc.UpdatedAt = parseTime(updatedAt)
	return &inc, nil
}

func scanIncidentRow(row *sql.Row) (*entity.Incident, error) { return scanIncident(row) }

// escapeLike neutraliza os coringas do LIKE no termo de busca do usuário.
func escapeLike(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	return r.Replace(s)
}
