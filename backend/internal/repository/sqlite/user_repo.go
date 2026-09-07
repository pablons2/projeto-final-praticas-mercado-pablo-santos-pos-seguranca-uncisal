package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"incidenttrack/internal/domain/entity"
	"incidenttrack/internal/domain/errdomain"
	"incidenttrack/internal/domain/repository"
)

// UserRepo implementa repository.UserRepository sobre SQLite.
type UserRepo struct{ db *DB }

// NewUserRepo cria o repositório de usuários.
func NewUserRepo(db *DB) *UserRepo { return &UserRepo{db: db} }

var _ repository.UserRepository = (*UserRepo)(nil)

const colsUser = "id, email, name, password_hash, created_at, updated_at"

// Create insere o usuário. Traduz violação de UNIQUE(email) em ErrEmailTaken.
func (r *UserRepo) Create(ctx context.Context, u *entity.User) error {
	const q = `INSERT INTO users (id, email, name, password_hash, created_at, updated_at)
	           VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, q,
		u.ID, u.Email, u.Name, u.PasswordHash, rfc3339(u.CreatedAt), rfc3339(u.UpdatedAt),
	)
	if isUniqueViolation(err) {
		return errdomain.ErrEmailTaken
	}
	if err != nil {
		return fmt.Errorf("sqlite: criar usuário: %w", err)
	}
	return nil
}

// FindByEmail busca por e-mail (case-insensitive).
func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	const q = `SELECT ` + colsUser + ` FROM users WHERE email = ? COLLATE NOCASE LIMIT 1`
	return r.scanOne(r.db.QueryRowContext(ctx, q, email))
}

// FindByID busca por identificador.
func (r *UserRepo) FindByID(ctx context.Context, id string) (*entity.User, error) {
	const q = `SELECT ` + colsUser + ` FROM users WHERE id = ? LIMIT 1`
	return r.scanOne(r.db.QueryRowContext(ctx, q, id))
}

func (r *UserRepo) scanOne(row *sql.Row) (*entity.User, error) {
	var u entity.User
	var createdAt, updatedAt string
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &createdAt, &updatedAt)
	if isNoRows(err) {
		return nil, errdomain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("sqlite: ler usuário: %w", err)
	}
	u.CreatedAt = parseTime(createdAt)
	u.UpdatedAt = parseTime(updatedAt)
	return &u, nil
}
