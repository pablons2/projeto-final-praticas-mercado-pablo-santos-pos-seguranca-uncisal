package auth_test

import (
	"context"
	"errors"
	"testing"

	"incidenttrack/internal/domain/errdomain"
	"incidenttrack/internal/pkg/hash"
	"incidenttrack/internal/repository/memory"
	authuc "incidenttrack/internal/usecase/auth"
)

// custo baixo de bcrypt só nos testes, para não deixá-los lentos.
func testHasher() hash.Hasher { return hash.NewBcrypt(4) }

func TestRegisterUser_Success(t *testing.T) {
	users := memory.NewUserRepo()
	uc := authuc.NewRegisterUser(users, testHasher())

	u, err := uc.Execute(context.Background(), authuc.RegisterInput{
		Name: "Ana", Email: "Ana@Example.com ", Password: "senha1234",
	})
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}
	if u.ID == "" {
		t.Error("esperava ID gerado")
	}
	if u.Email != "ana@example.com" {
		t.Errorf("e-mail deveria ser normalizado, veio %q", u.Email)
	}
	if u.PasswordHash == "" || u.PasswordHash == "senha1234" {
		t.Error("senha não pode ser armazenada em texto plano")
	}
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	users := memory.NewUserRepo()
	uc := authuc.NewRegisterUser(users, testHasher())
	in := authuc.RegisterInput{Name: "Ana", Email: "ana@example.com", Password: "senha1234"}

	if _, err := uc.Execute(context.Background(), in); err != nil {
		t.Fatalf("primeiro registro falhou: %v", err)
	}
	_, err := uc.Execute(context.Background(), in)
	if !errors.Is(err, errdomain.ErrEmailTaken) {
		t.Fatalf("esperava ErrEmailTaken, veio %v", err)
	}
}

func TestLoginUser_Success(t *testing.T) {
	users := memory.NewUserRepo()
	h := testHasher()
	reg := authuc.NewRegisterUser(users, h)
	if _, err := reg.Execute(context.Background(), authuc.RegisterInput{
		Name: "Ana", Email: "ana@example.com", Password: "senha1234",
	}); err != nil {
		t.Fatalf("setup: %v", err)
	}

	dummy, _ := h.Hash("dummy-000000")
	login := authuc.NewLoginUser(users, h, dummy)

	u, err := login.Execute(context.Background(), authuc.LoginInput{
		Email: "ana@example.com", Password: "senha1234",
	})
	if err != nil {
		t.Fatalf("login válido falhou: %v", err)
	}
	if u.Email != "ana@example.com" {
		t.Errorf("usuário errado: %q", u.Email)
	}
}

func TestLoginUser_WrongPassword(t *testing.T) {
	users := memory.NewUserRepo()
	h := testHasher()
	reg := authuc.NewRegisterUser(users, h)
	_, _ = reg.Execute(context.Background(), authuc.RegisterInput{
		Name: "Ana", Email: "ana@example.com", Password: "senha1234",
	})
	dummy, _ := h.Hash("dummy-000000")
	login := authuc.NewLoginUser(users, h, dummy)

	_, err := login.Execute(context.Background(), authuc.LoginInput{
		Email: "ana@example.com", Password: "errada9999",
	})
	if !errors.Is(err, errdomain.ErrInvalidCredentials) {
		t.Fatalf("esperava ErrInvalidCredentials, veio %v", err)
	}
}

func TestLoginUser_UnknownEmail(t *testing.T) {
	users := memory.NewUserRepo()
	h := testHasher()
	dummy, _ := h.Hash("dummy-000000")
	login := authuc.NewLoginUser(users, h, dummy)

	_, err := login.Execute(context.Background(), authuc.LoginInput{
		Email: "ninguem@example.com", Password: "qualquer123",
	})
	// Mesma resposta genérica de "senha errada" — sem vazar que o e-mail não existe.
	if !errors.Is(err, errdomain.ErrInvalidCredentials) {
		t.Fatalf("esperava ErrInvalidCredentials, veio %v", err)
	}
}
