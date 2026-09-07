// Package validator reúne sanitização e regras de validação reutilizáveis.
// A validação no backend é a linha de defesa definitiva: o frontend
// valida por UX, mas nada que chega aqui é considerado confiável
// (docs/plano.md §8; OWASP A03/A05 — entrada tratada antes de qualquer uso).
package validator

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Errors acumula erros por campo no formato que o frontend consome:
// { "fields": { "email": ["mensagem"] } } (ver frontend/src/lib/http.ts).
type Errors struct {
	fields map[string][]string
}

// New cria um acumulador vazio.
func New() *Errors { return &Errors{fields: map[string][]string{}} }

// Add registra uma mensagem para o campo.
func (e *Errors) Add(field, msg string) {
	e.fields[field] = append(e.fields[field], msg)
}

// Ok indica ausência de erros.
func (e *Errors) Ok() bool { return len(e.fields) == 0 }

// Fields devolve o mapa (nil se vazio) para serialização.
func (e *Errors) Fields() map[string][]string {
	if len(e.fields) == 0 {
		return nil
	}
	return e.fields
}

// --- Sanitização ---

// Trim remove espaços das pontas.
func Trim(s string) string { return strings.TrimSpace(s) }

// NormalizeEmail: trim + lowercase. E-mail é case-insensitive no domínio.
func NormalizeEmail(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

// CollapseSpaces troca sequências de espaço em branco por um único espaço
// e remove caracteres de controle (exceto \n e \t em textos longos).
func CollapseSpaces(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	prevSpace := false
	for _, r := range s {
		// Whitespace (inclui \t e \n, que também são "control") vira um
		// único espaço; os demais caracteres de controle são descartados.
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteRune(' ')
			}
			prevSpace = true
			continue
		}
		if unicode.IsControl(r) {
			continue
		}
		prevSpace = false
		b.WriteRune(r)
	}
	return strings.TrimSpace(b.String())
}

// StripControl remove apenas caracteres de controle perigosos, preservando
// quebras de linha e tabs (para campos de texto multilinha como descrição).
func StripControl(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' || r == '\r' {
			return r
		}
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, s)
}

// --- Regras ---

var emailRe = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// IsEmail valida o formato (não a existência) do e-mail.
func IsEmail(s string) bool {
	return len(s) <= 254 && emailRe.MatchString(s)
}

// RuneLenBetween verifica o comprimento em runas (não bytes) no intervalo
// fechado [min, max].
func RuneLenBetween(s string, min, max int) bool {
	n := utf8.RuneCountInString(s)
	return n >= min && n <= max
}

// PasswordPolicy: mínimo 8 caracteres, com ao menos uma letra e um dígito
// (docs/plano.md §7, A07). Retorna a mensagem de erro ou "" se ok.
func PasswordPolicy(pw string) string {
	if utf8.RuneCountInString(pw) < 8 {
		return "Deve ter ao menos 8 caracteres."
	}
	if len(pw) > 72 {
		return "Deve ter no máximo 72 caracteres."
	}
	var hasLetter, hasDigit bool
	for _, r := range pw {
		switch {
		case unicode.IsLetter(r):
			hasLetter = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	if !hasLetter || !hasDigit {
		return "Deve conter ao menos uma letra e um número."
	}
	return ""
}

// InSet indica se v está na lista de valores permitidos.
func InSet[T comparable](v T, allowed ...T) bool {
	for _, a := range allowed {
		if v == a {
			return true
		}
	}
	return false
}
