package entity

import "time"

// User é o usuário autenticável do sistema. A senha nunca trafega nem é
// persistida em texto plano: apenas PasswordHash (bcrypt) é armazenado
// (OWASP A02 — Cryptographic Failures; docs/plano.md §5).
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"` // nunca serializado em resposta
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
