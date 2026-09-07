// Package sqlite implementa as portas de domain/repository sobre SQLite,
// usando database/sql com o driver puro-Go modernc.org/sqlite (sem CGO,
// compila e implanta fácil num Free Tier; docs/plano.md §4).
//
// OWASP A03/A05 — Injection: toda query usa placeholders posicionais (?)
// resolvidos pelo driver. Não há concatenação de input do usuário em SQL
// em nenhum ponto deste pacote.
package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// DB embrulha *sql.DB com a configuração adequada ao SQLite.
type DB struct{ *sql.DB }

// Open abre o banco, aplica PRAGMAs de segurança/robustez e roda as
// migrações embutidas. dsn é algo como "file:data/incidenttrack.db".
func Open(ctx context.Context, dsn string) (*DB, error) {
	// PRAGMAs via querystring do driver:
	//  - busy_timeout: espera em vez de falhar com "database is locked".
	//  - journal_mode=WAL: leituras concorrentes com escrita.
	//  - foreign_keys=ON: respeita o ON DELETE CASCADE.
	//  - synchronous=NORMAL: durabilidade adequada com WAL.
	sep := "?"
	if strings.Contains(dsn, "?") {
		sep = "&"
	}
	full := dsn + sep + strings.Join([]string{
		"_pragma=busy_timeout(5000)",
		"_pragma=journal_mode(WAL)",
		"_pragma=foreign_keys(ON)",
		"_pragma=synchronous(NORMAL)",
	}, "&")

	sqlDB, err := sql.Open("sqlite", full)
	if err != nil {
		return nil, fmt.Errorf("sqlite: open: %w", err)
	}

	// SQLite lida melhor com um único escritor. Mantemos o pool pequeno
	// para evitar contenção de lock em cargas de trabalho pequenas.
	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetConnMaxLifetime(time.Hour)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("sqlite: ping: %w", err)
	}

	db := &DB{sqlDB}
	if err := db.migrate(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}
	return db, nil
}

func (db *DB) migrate(ctx context.Context) error {
	entries, err := migrationsFS.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("sqlite: ler migrações: %w", err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		content, err := migrationsFS.ReadFile("migrations/" + e.Name())
		if err != nil {
			return fmt.Errorf("sqlite: ler %s: %w", e.Name(), err)
		}
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("sqlite: aplicar %s: %w", e.Name(), err)
		}
	}
	return nil
}

// isUniqueViolation detecta erro de UNIQUE constraint do driver modernc.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unique constraint failed") ||
		strings.Contains(msg, "constraint failed: unique")
}

// rfc3339 formata timestamps para armazenamento textual estável.
func rfc3339(t time.Time) string { return t.UTC().Format(time.RFC3339Nano) }

// parseTime lê o timestamp textual do banco.
func parseTime(s string) time.Time {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02 15:04:05.999999999-07:00", "2006-01-02 15:04:05"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC()
		}
	}
	return time.Time{}
}

// errNoRows normaliza sql.ErrNoRows para quem chama comparar.
var errNoRows = sql.ErrNoRows

func isNoRows(err error) bool { return errors.Is(err, errNoRows) }
