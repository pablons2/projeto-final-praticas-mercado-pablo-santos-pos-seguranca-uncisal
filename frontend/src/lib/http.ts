/**
 * HTTP client — thin fetch wrapper.
 *
 * Security-relevant behaviours (maps to docs/plano.md §5, §7):
 *  - credentials: 'include' so the HttpOnly session cookie is sent.
 *  - X-Requested-With header on every request (CSRF mitigation for
 *    cookie-based auth, complementing SameSite=Strict on the server).
 *  - Never reads or writes token from JS-accessible storage.
 *  - Surfaces a typed ApiError so callers can branch on status without
 *    leaking raw server stack traces to the UI.
 */

export class ApiError extends Error {
  readonly status: number
  readonly fields?: Record<string, string[]>
  constructor(status: number, message: string, fields?: Record<string, string[]>) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.fields = fields
  }
}

export interface RequestOptions extends Omit<RequestInit, 'body'> {
  body?: unknown
  /** Expected ok-statuses; anything else throws. Defaults to [200,201]. */
  expect?: number[]
}

const DEFAULT_EXPECT = [200, 201]

const BASE_URL = import.meta.env.VITE_API_BASE_URL ?? ''

async function request<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const { body, expect = DEFAULT_EXPECT, headers, ...rest } = opts

  const init: RequestInit = {
    ...rest,
    credentials: 'include',
    headers: {
      'X-Requested-With': 'XMLHttpRequest',
      Accept: 'application/json',
      ...(body !== undefined ? { 'Content-Type': 'application/json' } : {}),
      ...headers,
    },
    ...(body !== undefined ? { body: JSON.stringify(body) } : {}),
  }

  let res: Response
  try {
    res = await fetch(`${BASE_URL}${path}`, init)
  } catch {
    throw new ApiError(0, 'Falha de conexão. Verifique sua rede.')
  }

  if (!expect.includes(res.status)) {
    throw await toApiError(res)
  }

  // 204 No Content
  if (res.status === 204 || res.headers.get('content-length') === '0') {
    return undefined as T
  }

  const contentType = res.headers.get('content-type') ?? ''
  if (!contentType.includes('application/json')) {
    return undefined as T
  }
  return (await res.json()) as T
}

async function toApiError(res: Response): Promise<ApiError> {
  let message = `Erro inesperado (${res.status}).`
  let fields: Record<string, string[]> | undefined

  try {
    const data = await res.json()
    if (typeof data?.error === 'string') message = data.error
    if (typeof data?.message === 'string') message = data.message
    if (data?.fields && typeof data.fields === 'object') fields = data.fields
  } catch {
    // Non-JSON error body — keep generic message, never expose raw text.
  }

  // Friendly mapping for common statuses
  if (res.status === 401) message = 'Sessão expirada. Faça login novamente.'
  if (res.status === 403) message = 'Você não tem permissão para esta ação.'
  if (res.status === 429) {
    const retry = res.headers.get('retry-after')
    message = retry
      ? `Muitas tentativas. Tente novamente em ${retry} s.`
      : 'Muitas tentativas. Tente novamente mais tarde.'
  }
  if (res.status >= 500) message = 'Algo deu errado no servidor. Tente novamente.'

  return new ApiError(res.status, message, fields)
}

export const http = {
  get: <T>(path: string, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: 'GET' }),
  post: <T>(path: string, body?: unknown, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: 'POST', body }),
  put: <T>(path: string, body?: unknown, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: 'PUT', body }),
  delete: <T>(path: string, opts?: RequestOptions) =>
    request<T>(path, { ...opts, method: 'DELETE', expect: [200, 204] }),
}
