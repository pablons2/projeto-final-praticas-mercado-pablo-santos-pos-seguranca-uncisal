package token

import (
	"sync"
	"time"
)

// Denylist mantém, em memória, os jti de tokens revogados antes do
// vencimento natural (usado no logout). JWT é stateless; esta lista é a
// extensão opcional citada em docs/plano.md §5 para invalidação imediata.
//
// Escopo: processo único. Num cluster seria necessário um store
// compartilhado (Redis), mas o Free Tier deste trabalho roda 1 instância.
type Denylist struct {
	mu      sync.RWMutex
	entries map[string]time.Time // jti -> instante de expiração original
	now     func() time.Time
}

// NewDenylist cria a lista e dispara um coletor periódico que remove
// entradas já expiradas (elas não precisam mais ser bloqueadas).
func NewDenylist(gcEvery time.Duration) *Denylist {
	d := &Denylist{entries: make(map[string]time.Time), now: time.Now}
	if gcEvery <= 0 {
		gcEvery = time.Minute
	}
	go func() {
		t := time.NewTicker(gcEvery)
		defer t.Stop()
		for range t.C {
			d.gc()
		}
	}()
	return d
}

// Revoke marca um jti como revogado até o instante exp.
func (d *Denylist) Revoke(jti string, exp time.Time) {
	if jti == "" {
		return
	}
	d.mu.Lock()
	d.entries[jti] = exp
	d.mu.Unlock()
}

// Revoked informa se o jti está revogado e ainda dentro da validade.
func (d *Denylist) Revoked(jti string) bool {
	if jti == "" {
		return false
	}
	d.mu.RLock()
	exp, ok := d.entries[jti]
	d.mu.RUnlock()
	return ok && d.now().Before(exp)
}

func (d *Denylist) gc() {
	now := d.now()
	d.mu.Lock()
	for jti, exp := range d.entries {
		if now.After(exp) {
			delete(d.entries, jti)
		}
	}
	d.mu.Unlock()
}
