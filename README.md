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

Código escrito e auditado com assistência de IA — ver
[Desenvolvimento assistido por IA](#desenvolvimento-assistido-por-ia).

---

## Aplicação em produção

A aplicação está no ar, acessível por **IP público** (sem domínio):

**https://163.176.179.49**

- HTTP redireciona automaticamente para HTTPS (`301`).
- Certificado Let's Encrypt emitido para o IP, renovado automaticamente.
- Hospedagem: Oracle Cloud (Free Tier), Ubuntu 26.04.1 LTS, região `sa-saopaulo-1`.

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

## Testes de ponta a ponta (smoke test)

Pasta [`e2e/`](e2e/): testes com **Playwright** que exercitam a aplicação
pela interface do navegador e pela API:

1. registrar → logar → criar, editar e excluir incidente → logout;
2. incidente de outro usuário retorna `404` (controle de acesso, OWASP A01).

Rodar com a stack local no ar (`make up`):

```bash
cd e2e
npm ci
npm test
```

**No CI** (job `e2e` em [`ci.yml`](.github/workflows/ci.yml)): a stack sobe com
`make up`, o teste roda com o Chromium do Playwright e o job **anexa a gravação
em vídeo e o relatório HTML** como artefato `e2e-gravacao-e-relatorio` (aba
*Actions* → execução → *Artifacts*). Uma falha no smoke test bloqueia o deploy.

**Gravação de referência:** [`evidencias/e2e/smoke-gravacao.webm`](evidencias/e2e/smoke-gravacao.webm).

**Localmente no Ubuntu 26.04** (o Playwright não baixa Chromium nem ffmpeg para
essa versão): use o Chrome do sistema com `CHROME_PATH=/usr/bin/google-chrome npm test`
e aponte o ffmpeg do Playwright para o ffmpeg do sistema:

```bash
D=~/.cache/ms-playwright/ffmpeg-1011
mkdir -p "$D" && ln -sf /usr/bin/ffmpeg "$D/ffmpeg-linux" && touch "$D/INSTALLATION_COMPLETE"
```

Cada execução usa e-mails únicos, então pode ser repetida.

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
├── .github/workflows/  esteira de CI/CD (deploy.yml)
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

---

## Desenvolvimento assistido por IA

Conforme o escopo da disciplina, todo o desenvolvimento foi feito com
inteligência artificial na **escrita** e na **auditoria** do código:

- **Escrita:** geração da estrutura Clean Architecture do backend Go, dos
  usecases, middlewares de segurança, handlers HTTP e da SPA React
  (componentes, hooks, schemas de validação).
- **Auditoria:** revisão iterativa das mitigações OWASP Top 10:2025
  (controle de acesso por `owner_id`, prepared statements, política de
  senha/bcrypt, headers de segurança, tratamento de `panic`), do
  `.gitignore` e da esteira de CI/CD, checando cada item contra o
  `Escopo_e_elementos_obrigatorios.md`.
- **Ambiente:** IDE baseada em IA (Google Antigravity / equivalente),
  operada pelo aluno.

A revisão final e a responsabilidade pelo código entregue são do autor.

---

## Testes TLS (acesso por IP)

Testes exigidos pelo escopo para o acesso por IP. O Qualys SSL Labs (nota A)
é o critério para domínio e não se aplica aqui.

| Teste | Alvo | Resultado | Print |
|---|---|---|---|
| [SSL.org — SSL Certificate Checker](https://www.ssl.org/) | `163.176.179.49` | Certificate Trusted: **Yes** · Good signature · Good key · cadeia completa · TLS 1.2/1.3 | [ver print](evidencias/tls/ssl-org-163.176.179.49.jpg) |
| [DigiCert — TLS quantum readiness check](https://www.digicert.com/pqc-checker) | `163.176.179.49` | **Pass**: TLS 1.3 habilitado e troca de chaves pós-quântica ativa | [ver print](evidencias/tls/digicert-pqc-163.176.179.49.jpg) |

Configuração que sustenta os resultados: proxy TLS em `nginx:1.29-alpine`
(OpenSSL 3.5) com `X25519MLKEM768` — ver [`deploy/README.md`](deploy/README.md) §9.

---

## CI/CD — GitHub Actions

Dois workflows encadeados em `.github/workflows/`:

**CI — [`ci.yml`](.github/workflows/ci.yml)** — roda em todo `push` e todo
pull request; é o portão de qualidade:

- **`backend`** — `gofmt` + `go vet` + `go test -race` + build de sanidade (Go 1.23).
- **`frontend`** — `npm ci` + `npm run lint` (oxlint) + type-check e build (Node 22).

**CD — [`deploy.yml`](.github/workflows/deploy.yml)** — disparado por
`workflow_run` **só quando a CI passa na `main`** (nunca implanta com
teste vermelho); também aceita disparo manual. Conecta na VM por SSH
(`IdentitiesOnly` + `known_hosts` fixo), faz `git reset --hard` no SHA
testado pela CI, gera o `.env` de produção (`APP_ENV=production`,
`COOKIE_SECURE=true`, `CERT_NAME`) com `umask 077`, sobe a stack com
`docker-compose.yml` + [`docker-compose.prod.yml`](docker-compose.prod.yml)
— este último adiciona o **proxy TLS** ([`deploy/proxy/`](deploy/proxy/),
`nginx:1.29-alpine` / OpenSSL 3.5) que termina HTTPS com key exchange
**pós-quântico** (`X25519MLKEM768`), redirect HTTP→HTTPS e HSTS — e valida
`GET /healthz`. Enquanto o certificado não existir na VM, sobe só a stack
base. Bootstrap da VM e do TLS em [`deploy/README.md`](deploy/README.md).

Nenhuma credencial fica nos workflows — tudo vem de **GitHub Secrets**:
`SSH_PRIVATE_KEY`, `SSH_KNOWN_HOSTS`, `SERVER_HOST`, `SERVER_USER`,
`DEPLOY_PATH`, `JWT_SECRET`, `APP_PORT`, `APP_ORIGIN` (formatos e setup
único da VM documentados no cabeçalho do `deploy.yml`).
