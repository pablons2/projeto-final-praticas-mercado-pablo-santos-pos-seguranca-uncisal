package validator_test

import (
	"testing"

	"incidenttrack/internal/pkg/validator"
)

func TestPasswordPolicy(t *testing.T) {
	cases := map[string]bool{ // senha -> ok?
		"senha1234":  true,
		"abc12345":   true,
		"curto1":     false, // < 8
		"semnumeros": false, // sem dígito
		"12345678":   false, // sem letra
		"":           false,
	}
	for pw, wantOK := range cases {
		gotOK := validator.PasswordPolicy(pw) == ""
		if gotOK != wantOK {
			t.Errorf("PasswordPolicy(%q): ok=%v, queria %v", pw, gotOK, wantOK)
		}
	}
}

func TestIsEmail(t *testing.T) {
	valid := []string{"a@b.com", "ana.silva+tag@example.co.uk"}
	invalid := []string{"", "sem-arroba", "a@b", "a@@b.com", "a b@c.com"}
	for _, e := range valid {
		if !validator.IsEmail(e) {
			t.Errorf("%q deveria ser válido", e)
		}
	}
	for _, e := range invalid {
		if validator.IsEmail(e) {
			t.Errorf("%q deveria ser inválido", e)
		}
	}
}

func TestNormalizeEmail(t *testing.T) {
	if got := validator.NormalizeEmail("  Ana@Example.COM "); got != "ana@example.com" {
		t.Errorf("normalização errada: %q", got)
	}
}

func TestCollapseSpaces(t *testing.T) {
	if got := validator.CollapseSpaces("  ola\t\tmundo\n  "); got != "ola mundo" {
		t.Errorf("colapso errado: %q", got)
	}
}

func TestErrors_Accumulate(t *testing.T) {
	e := validator.New()
	if !e.Ok() {
		t.Fatal("novo acumulador deveria estar Ok")
	}
	e.Add("email", "obrigatório")
	e.Add("email", "inválido")
	if e.Ok() {
		t.Fatal("deveria ter erros")
	}
	if got := len(e.Fields()["email"]); got != 2 {
		t.Errorf("esperava 2 mensagens, veio %d", got)
	}
}
