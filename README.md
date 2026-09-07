# IncidentTrack

> **Aluno:** Pablo Santos
> **Disciplina:** Projeto Aplicado: Práticas de Mercado
> **Curso:** Pós-graduação Lato Sensu em Segurança da Informação e Análise Forense — UNCISAL
> **Repositório:** <https://github.com/pablons2/projeto-final-praticas-mercado-pablo-santos-pos-seguranca-uncisal>

Sistema web de registro de incidentes de segurança da informação: login,
página interna autenticada, logout e CRUD de incidentes restrito ao dono
do recurso. Entrega da disciplina **Projeto Aplicado: Práticas de
Mercado** (Pós — Segurança da Informação e Análise Forense, UNCISAL).

Stack: **Go** (API REST, `net/http` puro, SQLite sem CGO) + **React**
(Vite, TypeScript). Norma de segurança aplicada: **OWASP Top 10:2025**.

Planejamento completo em [`docs/plano.md`](docs/plano.md).

---

## Executar com Docker (recomendado)

Pré-requisitos: Docker + Compose v2 (ou Podman: adicione
`COMPOSE="podman compose"` aos comandos `make`).

```bash
cp .env.example .env
# edite .env e defina JWT_SECRET:  openssl rand -base64 48

docker compose up --build
```

Aplicação em **http://localhost:8080**. Um único serviço Nginx (`web`)
serve o SPA e faz proxy de `/api` para o serviço `backend` — mesma
origem, sem CORS, cookie de sessão funcionando como em produção.

Atalhos (`make help` lista todos):

| Comando | Ação |
|---|---|
| `make up` | gera `.env` (com `JWT_SECRET` aleatório), sobe a stack em background |
| `make up-fg` | sobe em foreground (logs no terminal) |
| `make logs` | segue os logs |
| `make down` | para os containers |
| `make clean` | para tudo e **apaga o volume do banco** |
| `make rebuild` | recompila as imagens sem cache |

Para usar outra porta: ajuste `APP_PORT` **e** `APP_ORIGIN` no `.env`
(precisam bater).

---

## Executar sem Docker

**Backend** (Go 1.23+):

```bash
cd backend
cp .env.example .env      # defina JWT_SECRET
go run ./cmd/api          # :8080  — veja backend/README.md
```

**Frontend** (Node 22+):

```bash
cd frontend
cp .env.example .env      # VITE_API_BASE_URL=http://localhost:8080
npm ci && npm run dev     # :5173
```

Nesse modo, defina no `.env` do backend
`CORS_ALLOWED_ORIGIN=http://localhost:5173`.

---

## Estrutura do repositório

```
.
├── backend/            API em Go (Clean Architecture — ver backend/README.md)
│   ├── cmd/api/         ponto de entrada + injeção de dependência
│   ├── internal/        domain · usecase · repository · handler · middleware · dto · pkg
│   └── Dockerfile       build multi-stage → imagem distroless não-root
├── frontend/           SPA em React + Vite
│   ├── src/             pages · features · components · context · routes · lib
│   ├── nginx.conf       serve o SPA + proxy reverso de /api
│   └── Dockerfile       build multi-stage → Nginx
├── docs/plano.md       plano de execução e leitura crítica do escopo
├── docker-compose.yml  orquestração (origem única via Nginx)
└── .env.example        variáveis do compose (JWT_SECRET, APP_PORT, APP_ORIGIN)
```

---

## Segurança — OWASP Top 10:2025

As três categorias destacadas pelo escopo são **A01**, **A05** e **A07**.
A tabela completa, com o arquivo e a função exata de cada mitigação e os
testes que a cobrem, está em **[`backend/README.md`](backend/README.md)**.
Resumo:

| Categoria | Mitigação | Onde |
|---|---|---|
| **A01 — Broken Access Control** | JWT em toda rota protegida + **todas as queries filtram por `owner_id`**; recurso de outro usuário → `404` | `backend/internal/middleware/auth.go`, `backend/internal/repository/sqlite/incident_repo.go`, `backend/internal/usecase/incident/` |
| **A05 — Injection** | 100% das queries com placeholders `?` (zero concatenação), body JSON limitado, campos desconhecidos rejeitados, `CHECK` de enum no schema | `backend/internal/repository/sqlite/`, `backend/internal/handler/http/decode.go` |
| **A07 — Authentication Failures** | bcrypt cost 12, política de senha, JWT HS256 15 min (`alg:none` rejeitado), rate limit de login por IP, resposta genérica anti-enumeração | `backend/internal/pkg/hash/`, `backend/internal/pkg/token/`, `backend/internal/middleware/ratelimit.go` |
| A04 — Cryptographic Failures | segredo JWT ≥ 256 bits fora do código; cookie `HttpOnly` + `Secure` + `SameSite=Strict` | `backend/internal/handler/http/session.go`, `backend/internal/config/` |
| A02 — Security Misconfiguration | CSP, `X-Frame-Options`, HSTS; CORS de origem exata; erros genéricos ao cliente | `backend/internal/middleware/security_headers.go`, `cors.go` |
| A10 — Exceptional Conditions | `panic` recuperado sem stack trace ao cliente; nenhum `err` ignorado | `backend/internal/middleware/recover.go` |

Complemento CSRF (sessão em cookie): header `X-Requested-With` obrigatório
nas mutações + verificação de `Origin` —
`backend/internal/middleware/csrf.go`.
