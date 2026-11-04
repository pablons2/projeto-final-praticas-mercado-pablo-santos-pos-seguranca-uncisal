import { test, expect, request as apiRequest, type APIRequestContext } from '@playwright/test'

// Smoke test (docs/check-list-entrega.md — "Smoke test manual da aplicação"):
//   registrar → logar → criar, editar, excluir incidente → logout
//   + incidente de outro usuário deve retornar 404.
// Cada execução usa e-mails únicos, então pode ser repetida sem limpeza.

const run = Date.now()
const PASSWORD = 'Senha12345'
const userA = { name: 'Smoke A', email: `smoke-a-${run}@example.com`, password: PASSWORD }
const userB = { name: 'Smoke B', email: `smoke-b-${run}@example.com`, password: PASSWORD }
// Usuários próprios do teste de acesso cruzado (e-mails distintos do fluxo de UI).
const apiA = { name: 'Smoke API A', email: `smoke-api-a-${run}@example.com`, password: PASSWORD }
const apiB = { name: 'Smoke API B', email: `smoke-api-b-${run}@example.com`, password: PASSWORD }
const titulo = `Incidente smoke ${run}`
const tituloEditado = `${titulo} (editado)`

// Cabeçalho exigido pelo middleware CSRF em toda mutação da API.
const CSRF = { 'X-Requested-With': 'fetch' }

async function apiSession(baseURL: string, user: typeof userA): Promise<APIRequestContext> {
  const ctx = await apiRequest.newContext({ baseURL })
  const reg = await ctx.post('/api/auth/register', {
    headers: CSRF,
    data: { name: user.name, email: user.email, password: user.password },
  })
  expect(reg.status(), 'cadastro via API').toBe(201)
  const login = await ctx.post('/api/auth/login', {
    headers: CSRF,
    data: { email: user.email, password: user.password },
  })
  expect(login.status(), 'login via API').toBe(200)
  return ctx
}

test('fluxo completo: registrar, logar, CRUD de incidente e logout', async ({ page }) => {
  // 1. Registrar pela interface
  await page.goto('/register')
  await page.getByLabel('Nome').fill(userA.name)
  await page.getByLabel('E-mail').fill(userA.email)
  // Campos de senha: o rótulo tem marcador de obrigatório, então usa-se a ordem dos inputs.
  await page.locator('input[type="password"]').nth(0).fill(userA.password)
  await page.locator('input[type="password"]').nth(1).fill(userA.password)
  await page.getByRole('button', { name: 'Criar conta' }).click()

  // 2. Cadastro já faz login automático e leva ao painel. Espera a navegação
  //    (sem recarregar a página, para não interromper a requisição de cadastro).
  await expect(page).toHaveURL(/\/dashboard/, { timeout: 15_000 })

  // 2b. Logar de novo, com a sessão limpa, para validar o fluxo de login.
  await page.context().clearCookies()
  await page.goto('/login')
  await page.getByLabel('E-mail').fill(userA.email)
  await page.locator('input[name="password"]').fill(userA.password)
  await page.getByRole('button', { name: 'Entrar' }).click()
  await expect(page).toHaveURL(/\/dashboard/, { timeout: 15_000 })
  await expect(page.getByRole('heading', { name: 'Incidentes' })).toBeVisible()

  // 3. Criar incidente
  await page.getByRole('button', { name: 'Novo incidente' }).click()
  await page.getByLabel('Título').fill(titulo)
  await page.getByLabel('Descrição').fill('Criado pelo smoke test automatizado.')
  // Selects customizados (Radix): abre o combobox e escolhe a opção pelo texto exibido.
  await page.getByRole('combobox', { name: 'Categoria' }).click()
  await page.getByRole('option', { name: 'Phishing', exact: true }).click()
  await page.getByRole('combobox', { name: 'Severidade' }).click()
  await page.getByRole('option', { name: 'Alta', exact: true }).click()
  await page.getByRole('combobox', { name: 'Status' }).click()
  await page.getByRole('option', { name: 'Aberto', exact: true }).click()
  await page.getByRole('button', { name: 'Criar incidente' }).click()
  await expect(page.getByText(titulo).first()).toBeVisible()

  // 4. Editar incidente
  await page.getByRole('button', { name: `Editar ${titulo}` }).first().click()
  await page.getByLabel('Título').fill(tituloEditado)
  await page.getByRole('button', { name: 'Salvar alterações' }).click()
  await expect(page.getByText(tituloEditado).first()).toBeVisible()

  // 5. Excluir incidente (confirmação no diálogo)
  await page.getByRole('button', { name: `Excluir ${tituloEditado}` }).first().click()
  await page.getByRole('button', { name: 'Excluir', exact: true }).click()
  await expect(page.getByText(tituloEditado)).toHaveCount(0)

  // 6. Logout
  await page.getByRole('button', { name: 'Menu do usuário' }).click()
  await page.getByRole('menuitem', { name: /Sair/ }).click()
  await expect(page).toHaveURL(/\/login/)
})

test('incidente de outro usuário retorna 404 (controle de acesso, OWASP A01)', async ({ baseURL }) => {
  const base = baseURL ?? 'http://localhost:8080'

  // Usuário A cria um incidente
  const a = await apiSession(base, apiA)
  const created = await a.post('/api/incidents', {
    headers: CSRF,
    data: {
      titulo: `Privado do A ${run}`,
      descricao: 'Só o dono pode ver.',
      categoria: 'malware',
      severidade: 'critica',
      status: 'aberto',
    },
  })
  expect(created.status(), 'criação do incidente de A').toBe(201)
  const { id } = await created.json()

  // Usuário B tenta acessar o incidente de A
  const b = await apiSession(base, apiB)
  const res = await b.get(`/api/incidents/${id}`)
  expect(res.status(), 'incidente de outro usuário deve ser 404').toBe(404)

  // E não aparece na lista de B
  const list = await b.get('/api/incidents')
  expect(await list.text()).not.toContain(id)

  await a.dispose()
  await b.dispose()
})
