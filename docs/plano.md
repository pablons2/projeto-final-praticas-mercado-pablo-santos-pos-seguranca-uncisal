# Plano de Execução — Sistema Web com Login e CRUD (Go + React)

**Disciplina:** Projeto Aplicado: Práticas de Mercado — Pós-graduação Lato Sensu em Segurança da Informação e Análise Forense (UNCISAL)
**Referência:** `Escopo_e_elementos_obrigatorios.md`, versão 1.0
**Norma de segurança aplicável:** OWASP Top 10:2025 (lista vigente desde janeiro/2026 — substitui a de 2021)

---

## 1. Leitura crítica do escopo

O documento do professor divide a avaliação em três eixos interligados por uma esteira de CI/CD. Antes de decidir qualquer linha de código, vale separar o que é exigido em cada um, porque nem tudo cabe dentro de um repositório de aplicação — parte é operação de infraestrutura que só existe fora do código.

### Eixo 1 — Infraestrutura (Cloud Free Tier)
Exige VM em Ubuntu Server ou Debian, Nginx ou Apache como web server, acesso só por chave SSH (sem senha), firewall com portas mínimas, Fail2Ban na porta 22 (tolerância de 4 erros, ban de 24h), HTTPS via Certbot com renovação automática, redirecionamento HTTP→HTTPS, e nota A no Qualys SSL Labs com PQC ativado. Tudo isso é configuração de servidor real, executada dentro da conta de nuvem do aluno — não é algo que se resolve escrevendo código de aplicação.

### Eixo 2 — Repositório (GitHub)
Repositório público, `.gitignore` correto (proibido subir `.env`, chaves SSH, credenciais, senhas hardcoded), README.md funcionando como relatório técnico da entrega. Isso sim é parte do que entra no código-fonte.

### Eixo 3 — Desenvolvimento (a aplicação em si)
Aqui o escopo é propositalmente flexível: qualquer stack, qualquer arquitetura, sem exigência de banco de dados. Só três coisas são obrigatórias na aplicação: tela de login, página interna pós-autenticação, e botão de logout funcional. Além disso, o código precisa mitigar de forma comprovável no mínimo 3 categorias do OWASP Top 10:2025, documentadas no README com a indicação exata de onde estão no código.

### Integração contínua (CI/CD)
Pipeline no GitHub Actions disparado por `git push origin main`, usando GitHub Secrets para as credenciais de deploy — nunca hardcoded no workflow.

### O que efetivamente será construído por mim nesta entrega de código

Eu vou gerar o código-fonte completo do Eixo 3 (backend Go + frontend React), os artefatos de apoio ao Eixo 2 (`.gitignore`, README técnico com a seção de mitigação OWASP) e um esqueleto de workflow do GitHub Actions para o Eixo 3 do CI/CD. O Eixo 1 (provisionar a VM, gerar as chaves SSH reais, configurar Fail2Ban, rodar o Certbot, testar no Qualys) e a configuração final das credenciais reais no GitHub (2FA, Secrets com a chave SSH do servidor) dependem de contas e acessos que são seus — vou deixar um checklist detalhado dessas etapas no README, mas a execução em si é sua.

---

## 2. Tema escolhido para a aplicação

**IncidentTrack** — um sistema simples de registro de incidentes de segurança da informação.

Cada usuário autenticado cadastra, lista, edita e exclui incidentes (título, descrição, categoria, severidade, status). Escolhi esse tema em vez de algo genérico (lista de tarefas, agenda) porque ele conversa diretamente com o eixo da pós-graduação — dá para usar o próprio domínio do sistema (registro de incidentes) como pretexto narrativo dentro do README ao explicar as mitigações de segurança, o que reforça a leitura de que a aplicação foi pensada com intenção, não só para cumprir a exigência mínima de "ter um CRUD".

**Entidade principal — Incidente:**
`id`, `titulo`, `descricao`, `categoria` (phishing / malware / acesso_indevido / vazamento_dados / outro), `severidade` (baixa / media / alta / critica), `status` (aberto / em_analise / resolvido), `owner_id`, `created_at`, `updated_at`.

---

## 3. Arquitetura geral

Monorepo único, com separação física clara entre as duas aplicações:

```
projeto-incidenttrack/
├── backend/          → API em Go
├── frontend/          → SPA em React
├── .github/workflows/ → pipeline CI/CD
├── .gitignore
└── README.md
```

Backend e frontend são desacoplados por HTTP (REST + JSON). O frontend não sabe nada sobre a estrutura interna do Go, e o backend não sabe nada sobre React — a única interface entre os dois é o contrato da API.

---

## 4. Arquitetura do backend (Go) — Clean Architecture + SOLID básico

Camadas, de dentro para fora (regra de dependência: as camadas internas nunca importam as externas):

```
backend/
├── cmd/api/main.go              → ponto de entrada, injeção de dependência
├── internal/
│   ├── domain/                  → entidades + interfaces (portas). Zero dependências externas.
│   │   ├── entity/
│   │   └── repository/          → interfaces: IncidentRepository, UserRepository
│   ├── usecase/                 → regras de negócio, uma responsabilidade por arquivo
│   │   ├── auth/                → RegisterUser, LoginUser
│   │   └── incident/             → CreateIncident, ListIncidents, UpdateIncident, DeleteIncident
│   ├── repository/sqlite/        → implementação concreta das interfaces de domain
│   ├── handler/http/             → controllers HTTP, tradução request↔usecase
│   ├── middleware/               → auth (JWT), recover, cors, security headers, rate limit
│   ├── dto/                       → structs de entrada/saída com validação
│   └── pkg/
│       ├── hash/                  → bcrypt
│       ├── token/                  → geração/validação JWT
│       └── validator/               → sanitização e regras de validação
├── go.mod
└── .env.example
```

Onde entra cada princípio SOLID, de forma concreta (não como slide de aula):

- **SRP:** cada usecase resolve uma única ação de negócio (`CreateIncident` não lista, não atualiza). O handler HTTP só traduz request/response — não tem regra de negócio dentro dele.
- **OCP:** trocar a implementação de persistência (SQLite → Postgres, por exemplo) não exige tocar em nenhum usecase, porque eles dependem da interface `IncidentRepository`, não da implementação.
- **LSP:** qualquer implementação de `IncidentRepository` precisa cumprir o mesmo contrato — os testes dos usecases rodam contra a interface, com um repositório fake em memória, e continuam válidos se trocar o SQLite por outro banco.
- **ISP:** interfaces pequenas e específicas (`IncidentRepository`, `UserRepository`, `TokenService`) em vez de uma interface "repositório genérico" gigante que ninguém implementa por completo.
- **DIP:** o `usecase` depende da abstração definida em `domain`, nunca do pacote `repository/sqlite` diretamente. Quem decide qual implementação concreta usar é o `main.go`, na hora de montar o grafo de dependências.

**Persistência:** SQLite via `database/sql`, com *driver* puro Go (`modernc.org/sqlite`, sem CGO — mais fácil de compilar e implantar num Free Tier). O escopo não exige banco de dados, mas optei por usar um banco real (em vez de um mapa em memória) justamente para poder demonstrar a mitigação de Injection com *prepared statements* de forma genuína, não simulada.

---

## 5. Fluxo de autenticação e autorização

O mínimo exigido pelo escopo é: tela de login, página interna, logout funcional. Vou implementar também um registro simples de usuário — não é obrigatório, mas sem ele a "página interna" fica presa a um único usuário fixo, o que empobrece a demonstração de controle de acesso por dono do recurso (cada usuário só vê e edita os próprios incidentes).

- **Senha:** hash com `bcrypt` (cost 12), nunca armazenada nem logada em texto plano.
- **Token:** JWT (HS256), assinado com segredo de 256 bits vindo de variável de ambiente, expiração curta (15 minutos).
- **Transporte do token:** cookie `HttpOnly`, `Secure`, `SameSite=Strict` — não uso `localStorage` para o token, porque qualquer XSS no frontend conseguiria lê-lo ali. Guardado em cookie `HttpOnly`, o JavaScript do navegador simplesmente não tem acesso a ele.
- **CSRF:** como a autenticação passa a depender de cookie, adiciono verificação de um header customizado (`X-Requested-With`) nas rotas de mutação, complementando o `SameSite=Strict`.
- **Logout:** limpa o cookie no servidor (`Set-Cookie` com expiração no passado). Como JWT é stateless, uma extensão opcional (fora do mínimo) é manter uma lista de revogação em memória para invalidar o token antes do vencimento natural.
- **Rate limiting:** limite de tentativas de login por IP (ex.: 5 tentativas / 15 minutos) para dificultar força bruta.

---

## 6. Endpoints da API

| Método | Rota                  | Autenticado | Descrição                          |
|--------|-----------------------|:-----------:|-------------------------------------|
| POST   | `/api/auth/register`  | não         | Cria usuário                        |
| POST   | `/api/auth/login`     | não         | Autentica e emite cookie de sessão  |
| POST   | `/api/auth/logout`    | sim         | Invalida a sessão atual             |
| GET    | `/api/auth/me`        | sim         | Retorna o usuário autenticado       |
| GET    | `/api/incidents`      | sim         | Lista incidentes do usuário         |
| GET    | `/api/incidents/:id`  | sim         | Detalha um incidente próprio        |
| POST   | `/api/incidents`      | sim         | Cria incidente                      |
| PUT    | `/api/incidents/:id`  | sim         | Atualiza incidente próprio          |
| DELETE | `/api/incidents/:id`  | sim         | Remove incidente próprio            |

---

## 7. Mapeamento de mitigações — OWASP Top 10:2025

O escopo pede a comprovação de **no mínimo 3** categorias no README. Vou implementar mais do que o mínimo e, no README, destacar as três mais fortes (A01, A05, A07) como resposta direta à exigência, deixando as demais como reforço.

| Categoria | Onde no código | Como mitiga |
|---|---|---|
| **A01 — Broken Access Control** | `middleware/auth.go` + checagem de `owner_id` dentro de cada usecase de incidente | Toda rota protegida passa pelo middleware de JWT; além disso, update/delete verificam explicitamente se o incidente pertence ao usuário do token antes de agir — não basta estar logado, precisa ser o dono. |
| **A05 — Injection** | `repository/sqlite/*.go` | Todas as queries usam *prepared statements* com parâmetros posicionais, nunca concatenação de string. Nenhum input do usuário chega ao SQL sem passar antes pela camada de `dto`/`validator`. |
| **A07 — Authentication Failures** | `pkg/hash`, `pkg/token`, `middleware/ratelimit.go` | Senha com hash bcrypt (nunca reversível), política mínima de senha (8+ caracteres, letra e número), expiração curta de token, limite de tentativas de login por IP. |
| A04 — Cryptographic Failures | `pkg/hash`, `pkg/token`, cookies `Secure` | Bcrypt para senha, JWT assinado com segredo forte fora do código-fonte, cookies só trafegam sobre HTTPS. |
| A02 — Security Misconfiguration | `middleware/security_headers.go`, `middleware/cors.go` | CSP, `X-Content-Type-Options`, `X-Frame-Options`; CORS restrito à origem exata do frontend; modo debug desligado via variável de ambiente em produção; mensagens de erro genéricas para o cliente (detalhe fica só no log do servidor). |
| A10 — Mishandling of Exceptional Conditions | `middleware/recover.go` | Recuperação de `panic` em toda requisição, evitando queda do processo e vazamento de *stack trace* para o cliente; erros tratados explicitamente em todas as camadas, nenhum `err` descartado silenciosamente. |

---

## 8. Frontend — estrutura e validação

```
frontend/
├── src/
│   ├── pages/          → Login, Register, Dashboard (CRUD de incidentes)
│   ├── components/      → formulário de incidente, tabela, modal de confirmação
│   ├── context/           → AuthContext (estado do usuário logado)
│   ├── routes/             → ProtectedRoute (redireciona para /login se não autenticado)
│   ├── services/            → cliente HTTP (fetch/axios com `credentials: include`)
│   └── utils/                → schemas de validação
```

- Validação de formulário com schema (mesma regra de senha e campos obrigatórios espelhando o backend), para dar feedback imediato ao usuário sem depender só da resposta da API.
- Sanitização de input antes de enviar (trim de espaços, normalização), embora a sanitização definitiva sempre aconteça no backend — o frontend nunca é a última linha de defesa.
- React já faz *escaping* de conteúdo por padrão no JSX, o que ajuda contra XSS refletido a partir dos próprios dados exibidos.
- Nenhum dado sensível (token, senha) fica em `localStorage` ou `sessionStorage`.

---

## 9. Testes

Testes unitários dos usecases usando as interfaces de `domain` com um repositório fake em memória (sem precisar subir banco real para testar regra de negócio) — isso só é possível porque a camada de usecase depende de abstração, não de implementação concreta. Cobertura mínima: registro, login, e as quatro operações de CRUD de incidente, incluindo o caso de tentativa de acesso a um incidente de outro usuário (deve retornar 403/404).

---

## 10. `.gitignore` e proteção de segredos (Eixo 2)

Itens obrigatórios fora do repositório: `.env` (backend e frontend), arquivo do banco SQLite (`*.db`), `node_modules/`, build do frontend (`dist/`), qualquer chave `.pem`/`.key`. Vou entregar `.env.example` em vez de `.env` real, com placeholders.

---

## 11. README.md — o que ele vai conter

Descrição do projeto e stack, instruções de execução local (backend e frontend), e a seção obrigatória pelo escopo com as 3 categorias OWASP escolhidas e a localização exata no código de cada mitigação (a tabela da seção 7 vai virar essa seção do README, resumida).

---

## 12. Esqueleto de CI/CD (GitHub Actions)

Vou entregar um `deploy.yml` inicial: build e teste do backend, build do frontend, e um job de deploy via SSH usando `secrets.SSH_PRIVATE_KEY` e `secrets.SERVER_HOST` — sem nenhum valor real, só a estrutura. A configuração dos Secrets de verdade no GitHub, com a chave SSH do seu servidor, é etapa sua.

---

## 13. O que fica exclusivamente com você (fora do que este chat pode executar)

- Criar a conta na cloud e provisionar a VM (Ubuntu/Debian).
- Gerar e instalar a chave SSH real, desabilitar login por senha.
- Configurar Nginx/Apache, Fail2Ban (4 tentativas, ban de 24h) e firewall.
- Rodar o Certbot e obter certificado válido para o IP público.
- Rodar o teste no Qualys SSL Labs e ajustar a configuração até nota A com PQC.
- Criar o repositório no GitHub, ativar 2FA na conta, cadastrar os Secrets reais do workflow.

Vou deixar essas etapas como checklist no README, na ordem em que fazem sentido operacionalmente.

---

## 14. Ordem de execução do código (próximas entregas)

1. Estrutura de pastas + `go.mod` + configuração de ambiente
2. Camada `domain` (entidades e interfaces)
3. Camada `usecase`
4. Camada `repository` (SQLite)
5. Middlewares de segurança (auth, recover, cors, headers, rate limit)
6. Handlers HTTP + rotas
7. `main.go` com injeção de dependência
8. Frontend — estrutura base, `AuthContext`, rotas protegidas
9. Telas de Login, Registro e Dashboard (CRUD)
10. Validações de formulário no frontend
11. Testes unitários do backend
12. `README.md`, `.gitignore`, `.env.example` e workflow de CI/CD
13. Revisão final cruzando com o checklist do professor

---

## 15. Checklist de conferência (Eixo 3, antes de considerar pronto)

- [ ] Login funcional, página interna só acessível autenticado, logout limpa a sessão
- [ ] CRUD completo de incidentes, restrito ao dono do recurso
- [ ] Senha com hash, nunca em texto plano nem em log
- [ ] Token não acessível via JavaScript (cookie `HttpOnly`)
- [ ] Queries parametrizadas, sem concatenação de SQL
- [ ] Validação de input espelhada no backend e no frontend
- [ ] Headers de segurança e CORS restritos configurados
- [ ] `panic` recuperado, sem *stack trace* exposto ao cliente
- [ ] `.env` fora do repositório, `.env.example` presente
- [ ] README com as 3 categorias OWASP e a localização de cada mitigação
- [ ] Workflow de CI/CD versionado, sem nenhuma credencial hardcoded