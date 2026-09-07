# IncidentTrack — Backend (Go)

API REST do IncidentTrack: autenticação por sessão em cookie e CRUD de
incidentes de segurança da informação, restrito ao dono do recurso.

Implementa o **Eixo 3** de `../docs/plano.md`. Stack: Go + `net/http` da
biblioteca padrão, SQLite via driver puro-Go (`modernc.org/sqlite`, sem
CGO). Sem framework web.

---

## Arquitetura (Clean Architecture + SOLID)

Regra de dependência: as camadas de dentro nunca importam as de fora.

```
cmd/api/main.go            ponto de entrada; injeção manual de dependências (DIP)
internal/
├── config/                configuração 100% via ambiente
├── domain/
│   ├── entity/            entidades puras (User, Incident, enums)
│   ├── repository/        PORTAS: IncidentRepository, UserRepository (ISP)
│   └── errdomain/         erros sentinela do domínio
├── usecase/
│   ├── auth/              RegisterUser, LoginUser        (1 ação por arquivo — SRP)
│   └── incident/          Create/List/Get/Update/Delete  (1 ação por arquivo — SRP)
├── repository/
│   ├── sqlite/            implementação concreta das portas (prepared statements)
│   └── memory/            implementação fake, usada nos testes de usecase (LSP)
├── handler/http/          controllers: só traduzem request↔usecase
├── middleware/            auth, recover, cors, security headers, rate limit, csrf, request-id
├── dto/                   structs de entrada/saída + validação/sanitização
└── pkg/
    ├── hash/              bcrypt
    ├── token/             JWT HS256 (implementação própria) + denylist em memória
    ├── validator/         regras de validação e sanitização reutilizáveis
    └── httpx/             escrita de resposta JSON e formato de erro
```

- **SRP** — cada usecase resolve uma única ação (`CreateIncident` não lista).
  O handler HTTP não contém regra de negócio.
- **OCP** — trocar SQLite por Postgres não toca nenhum usecase: eles dependem
  de `repository.IncidentRepository`, não da implementação.
- **LSP** — os testes de usecase rodam contra a interface, com o fake de
  `repository/memory`; continuariam válidos com outro banco.
- **ISP** — interfaces pequenas (`UserRepository`, `IncidentRepository`,
  `token.Service`, `hash.Hasher`) em vez de um "repositório genérico".
- **DIP** — o usecase depende da abstração de `domain`; quem escolhe a
  implementação concreta é `main.go`.

---

## Executar localmente

Pré-requisitos: Go 1.23+.

```bash
cd backend
cp .env.example .env
# gere um segredo forte e cole em JWT_SECRET:
#   openssl rand -base64 48
make run          # ou: go run ./cmd/api
```

A API sobe em `http://localhost:8080`. O arquivo SQLite é criado em
`backend/data/incidenttrack.db` (fora do Git).

### Comandos

```bash
make test         # testes (unitários de usecase + integração HTTP)
make cover        # testes com cobertura
make build        # binário estático em bin/api (CGO desligado)
make check        # fmt + vet + test
```

---

## Contrato da API

Base: `/api`. Corpo e respostas em JSON. Erros seguem
`{ "error": "mensagem", "fields": { "campo": ["..."] } }` — nunca stack
trace. Sessão trafega em cookie `HttpOnly`; o frontend envia
`X-Requested-With: XMLHttpRequest` e `credentials: 'include'`.

| Método | Rota                  | Auth | Descrição                          |
|--------|-----------------------|:----:|------------------------------------|
| GET    | `/api/health`         | não  | Liveness                           |
| POST   | `/api/auth/register`  | não  | Cria usuário → `201` + usuário     |
| POST   | `/api/auth/login`     | não  | Autentica → `200` + `Set-Cookie`   |
| POST   | `/api/auth/logout`    | sim  | Limpa cookie + revoga jti → `204`  |
| GET    | `/api/auth/me`        | sim  | Usuário autenticado                |
| GET    | `/api/incidents`      | sim  | Lista do usuário (`?status=&severidade=&categoria=&q=`) |
| GET    | `/api/incidents/{id}` | sim  | Detalha (404 se de outro dono)     |
| POST   | `/api/incidents`      | sim  | Cria → `201`                       |
| PUT    | `/api/incidents/{id}` | sim  | Substitui → `200`                  |
| DELETE | `/api/incidents/{id}` | sim  | Remove → `204`                     |

**Incidente:** `id`, `titulo`, `descricao`, `categoria`
(`phishing`/`malware`/`acesso_indevido`/`vazamento_dados`/`outro`),
`severidade` (`baixa`/`media`/`alta`/`critica`),
`status` (`aberto`/`em_analise`/`resolvido`), `owner_id`,
`created_at`, `updated_at`.

---

## Mitigações OWASP Top 10:2025 — localização no código

As três destacadas pelo escopo são **A01**, **A05** e **A07**; as demais
são reforço.

| Categoria | Arquivo(s) | Como mitiga |
|---|---|---|
| **A01 — Broken Access Control** | `internal/middleware/auth.go`; `internal/repository/sqlite/incident_repo.go` (`GetByIDForOwner`, `UpdateForOwner`, `DeleteForOwner`, `assertAffected`); `internal/usecase/incident/*.go` | Toda rota de incidente passa pelo middleware de sessão. Além disso, **todas as queries filtram por `owner_id`** e `update`/`delete` só afetam a linha do dono; recurso de outro usuário retorna `404` (não `403`), evitando enumeração de IDs. Testes: `internal/usecase/incident/incident_test.go`, `internal/handler/http/router_test.go` (`TestOwnershipIsolation`). |
| **A05 — Injection** | `internal/repository/sqlite/*.go`; `internal/handler/http/decode.go`; `internal/pkg/validator/validator.go`; `internal/repository/sqlite/migrations/0001_init.sql` | 100% das queries usam placeholders posicionais (`?`) resolvidos pelo driver — **zero concatenação de input em SQL**, inclusive no filtro dinâmico de listagem (só fragmentos fixos + `?`). Termo de busca do `LIKE` tem coringas escapados (`escapeLike`). Corpo JSON limitado a 1 MiB, campos desconhecidos rejeitados. `CHECK` de enum no schema como segunda barreira. |
| **A07 — Authentication Failures** | `internal/pkg/hash/hash.go`; `internal/pkg/token/jwt.go` + `denylist.go`; `internal/middleware/ratelimit.go`; `internal/usecase/auth/login.go`; `internal/pkg/validator/validator.go` (`PasswordPolicy`) | Senha com **bcrypt cost 12**, nunca em texto plano nem em log. Política de senha (8+ com letra e número). JWT **HS256** com segredo ≥ 256 bits fora do código, expiração de **15 min**, `alg:"none"` e troca de algoritmo rejeitados, comparação de assinatura em tempo constante. **Rate limit de login por IP** (5 / 15 min). Login usa hash *dummy* quando o e-mail não existe (nivela timing, anti-enumeração); mensagem de erro única e genérica. Denylist em memória revoga o token no logout. Testes: `internal/pkg/token/jwt_test.go`. |
| A04 — Cryptographic Failures | `internal/pkg/hash/hash.go`; `internal/pkg/token/jwt.go`; `internal/handler/http/session.go`; `internal/config/config.go` | Bcrypt para senha; JWT assinado com segredo forte validado na inicialização (≥ 32 bytes, obrigatório); cookie `Secure` (exigido em produção via `COOKIE_SECURE`), `HttpOnly`, `SameSite=Strict`. |
| A02 — Security Misconfiguration | `internal/middleware/security_headers.go`; `internal/middleware/cors.go`; `internal/handler/http/errors.go`; `internal/config/config.go` | CSP `default-src 'none'`, `X-Frame-Options: DENY`, `X-Content-Type-Options: nosniff`, `Referrer-Policy`, `Permissions-Policy`, HSTS em produção. CORS restrito à **origem exata** do frontend (`*` recusado na config); credenciais liberadas só para ela. Erros genéricos ao cliente, detalhe só no log do servidor. Segredos apenas via ambiente. |
| A10 — Mishandling of Exceptional Conditions | `internal/middleware/recover.go`; `internal/handler/http/errors.go`; `internal/repository/sqlite/*.go` | `panic` recuperado em toda requisição: stack trace só no log (com `request_id`), cliente recebe `500` genérico. Nenhum `err` descartado silenciosamente; erros de domínio mapeados explicitamente para status HTTP. |

---

## CSRF

Como a sessão é cookie-based, além de `SameSite=Strict` o middleware
`internal/middleware/csrf.go` exige, nas rotas de mutação, o header
`X-Requested-With` (não enviável por `<form>`/`<img>` cross-site sem
disparar preflight) e valida o header `Origin` quando presente.

---

## Testes

- **Usecase (unitário)** — `internal/usecase/**/*_test.go`: registro, login
  (senha errada e e-mail inexistente com a mesma resposta), e as quatro
  operações de CRUD, incluindo acesso a incidente de outro usuário.
  Rodam contra o fake em memória, sem subir banco.
- **HTTP (integração)** — `internal/handler/http/router_test.go`: fluxo
  completo pelo router real (middlewares, cookie, CSRF, isolamento por
  dono, 401 sem sessão, logout).
- **Cripto** — `internal/pkg/token/jwt_test.go`: adulteração de payload,
  `alg:"none"`, expiração, segredo errado, denylist.
