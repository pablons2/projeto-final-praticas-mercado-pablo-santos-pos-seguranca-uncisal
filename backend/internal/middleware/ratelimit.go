package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"incidenttrack/internal/pkg/httpx"
)

// Limiter é um limitador de taxa por chave, com janela fixa, em memória.
// Simples e suficiente para uma instância única no Free Tier
// (docs/plano.md §5). Para brute force de login usamos chave = IP.
type Limiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	hits   map[string]*window
	now    func() time.Time
}

type window struct {
	count int
	start time.Time
}

// NewLimiter cria um limitador de `limit` eventos por `window` e dispara
// um coletor que descarta janelas expiradas.
func NewLimiter(limit int, w time.Duration) *Limiter {
	l := &Limiter{
		limit:  limit,
		window: w,
		hits:   make(map[string]*window),
		now:    time.Now,
	}
	go l.gcLoop()
	return l
}

// Allow registra um evento para key e informa se ele é permitido. Quando
// negado, retryAfter é o tempo até a janela reabrir.
func (l *Limiter) Allow(key string) (ok bool, retryAfter time.Duration) {
	now := l.now()
	l.mu.Lock()
	defer l.mu.Unlock()

	w, exists := l.hits[key]
	if !exists || now.Sub(w.start) >= l.window {
		l.hits[key] = &window{count: 1, start: now}
		return true, 0
	}
	if w.count >= l.limit {
		return false, l.window - now.Sub(w.start)
	}
	w.count++
	return true, 0
}

// Reset zera o contador de uma chave (ex.: após login bem-sucedido).
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	delete(l.hits, key)
	l.mu.Unlock()
}

func (l *Limiter) gcLoop() {
	t := time.NewTicker(l.window)
	defer t.Stop()
	for range t.C {
		now := l.now()
		l.mu.Lock()
		for k, w := range l.hits {
			if now.Sub(w.start) >= l.window {
				delete(l.hits, k)
			}
		}
		l.mu.Unlock()
	}
}

// RateLimit aplica o limitador usando o IP do cliente como chave. Em
// excesso responde 429 com Retry-After (consumido por
// frontend/src/lib/http.ts).
func RateLimit(l *Limiter, trustProxy bool, msg string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := ClientIP(r, trustProxy)
			ok, retry := l.Allow(ip)
			if !ok {
				secs := int(retry.Seconds())
				if secs < 1 {
					secs = 1
				}
				w.Header().Set("Retry-After", strconv.Itoa(secs))
				httpx.WriteError(w, http.StatusTooManyRequests, msg)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
