package token

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func newSvc(t *testing.T, ttl time.Duration) *HS256 {
	t.Helper()
	return NewHS256([]byte("0123456789abcdef0123456789abcdef"), ttl)
}

func TestGenerateVerify_RoundTrip(t *testing.T) {
	s := newSvc(t, 15*time.Minute)
	raw, claims, err := s.Generate("user-1")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	got, err := s.Verify(raw)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	if got.Subject != "user-1" || got.TokenID != claims.TokenID {
		t.Errorf("claims divergentes: %+v vs %+v", got, claims)
	}
}

func TestVerify_RejectsTamperedPayload(t *testing.T) {
	s := newSvc(t, 15*time.Minute)
	raw, _, _ := s.Generate("user-1")
	parts := strings.Split(raw, ".")

	forged, _ := json.Marshal(payload{Sub: "admin", Exp: time.Now().Add(time.Hour).Unix()})
	parts[1] = base64.RawURLEncoding.EncodeToString(forged)

	if _, err := s.Verify(strings.Join(parts, ".")); err == nil {
		t.Fatal("payload adulterado deveria falhar na verificação de assinatura")
	}
}

func TestVerify_RejectsAlgNone(t *testing.T) {
	s := newSvc(t, 15*time.Minute)
	h, _ := json.Marshal(header{Alg: "none", Typ: "JWT"})
	p, _ := json.Marshal(payload{Sub: "user-1", Exp: time.Now().Add(time.Hour).Unix()})
	forged := base64.RawURLEncoding.EncodeToString(h) + "." +
		base64.RawURLEncoding.EncodeToString(p) + "."

	if _, err := s.Verify(forged); err == nil {
		t.Fatal(`token com alg "none" deveria ser rejeitado`)
	}
}

func TestVerify_RejectsExpired(t *testing.T) {
	s := newSvc(t, time.Minute)
	past := time.Now().Add(-2 * time.Hour)
	s.now = func() time.Time { return past }
	raw, _, _ := s.Generate("user-1")

	s.now = time.Now // "agora" real, bem depois do exp
	if _, err := s.Verify(raw); err != ErrExpired {
		t.Fatalf("esperava ErrExpired, veio %v", err)
	}
}

func TestVerify_RejectsWrongSecret(t *testing.T) {
	s := newSvc(t, 15*time.Minute)
	raw, _, _ := s.Generate("user-1")

	other := NewHS256([]byte("ffffffffffffffffffffffffffffffff"), 15*time.Minute)
	if _, err := other.Verify(raw); err != ErrSignature {
		t.Fatalf("esperava ErrSignature, veio %v", err)
	}
}

func TestDenylist_Revoke(t *testing.T) {
	d := &Denylist{entries: map[string]time.Time{}, now: time.Now}
	exp := time.Now().Add(time.Hour)
	d.Revoke("jti-1", exp)
	if !d.Revoked("jti-1") {
		t.Error("jti revogado deveria constar como revogado")
	}
	if d.Revoked("jti-2") {
		t.Error("jti não revogado não deveria constar")
	}
}
