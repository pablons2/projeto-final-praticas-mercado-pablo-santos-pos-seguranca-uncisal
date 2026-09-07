// Package dto define as estruturas de entrada/saída da API com sua
// validação. É a fronteira onde o input não confiável do cliente é
// sanitizado e checado antes de alcançar usecases e persistência
// (docs/plano.md §4, §8).
package dto

import (
	"time"

	"incidenttrack/internal/domain/entity"
	"incidenttrack/internal/pkg/validator"
)

// RegisterInput é o corpo de POST /api/auth/register.
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Normalize aplica sanitização in-place (trim, lowercase de e-mail).
func (in *RegisterInput) Normalize() {
	in.Name = validator.CollapseSpaces(in.Name)
	in.Email = validator.NormalizeEmail(in.Email)
	// Password não sofre trim: espaços podem ser intencionais na senha.
}

// Validate retorna os erros por campo (nil-safe via Errors.Ok()).
func (in *RegisterInput) Validate() *validator.Errors {
	e := validator.New()
	if !validator.RuneLenBetween(in.Name, 2, 80) {
		e.Add("name", "Deve ter entre 2 e 80 caracteres.")
	}
	if !validator.IsEmail(in.Email) {
		e.Add("email", "Informe um e-mail válido.")
	}
	if msg := validator.PasswordPolicy(in.Password); msg != "" {
		e.Add("password", msg)
	}
	return e
}

// LoginInput é o corpo de POST /api/auth/login.
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Normalize sanitiza o e-mail.
func (in *LoginInput) Normalize() {
	in.Email = validator.NormalizeEmail(in.Email)
}

// Validate faz apenas checagem de presença/forma. A verificação real de
// credenciais acontece no usecase, com mensagem genérica única (A07).
func (in *LoginInput) Validate() *validator.Errors {
	e := validator.New()
	if in.Email == "" {
		e.Add("email", "Campo obrigatório.")
	}
	if in.Password == "" {
		e.Add("password", "Campo obrigatório.")
	}
	return e
}

// UserResponse é a projeção pública de um usuário (sem hash de senha).
type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// NewUserResponse converte a entidade em resposta.
func NewUserResponse(u *entity.User) UserResponse {
	return UserResponse{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt}
}
