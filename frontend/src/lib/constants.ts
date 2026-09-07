/**
 * Shared domain constants — single source of truth for enums used by
 * both the zod schemas and the Select/Filter UI. Mirrors the backend
 * entity defined in docs/plano.md §2.
 */

export const CATEGORIAS = [
  'phishing',
  'malware',
  'acesso_indevido',
  'vazamento_dados',
  'outro',
] as const
export type Categoria = (typeof CATEGORIAS)[number]

export const SEVERIDADES = ['baixa', 'media', 'alta', 'critica'] as const
export type Severidade = (typeof SEVERIDADES)[number]

export const STATUS = ['aberto', 'em_analise', 'resolvido'] as const
export type Status = (typeof STATUS)[number]

/** Human-readable labels for PT-BR UI. */
export const CATEGORIA_LABELS: Record<Categoria, string> = {
  phishing: 'Phishing',
  malware: 'Malware',
  acesso_indevido: 'Acesso indevido',
  vazamento_dados: 'Vazamento de dados',
  outro: 'Outro',
}

export const SEVERIDADE_LABELS: Record<Severidade, string> = {
  baixa: 'Baixa',
  media: 'Média',
  alta: 'Alta',
  critica: 'Crítica',
}

export const STATUS_LABELS: Record<Status, string> = {
  aberto: 'Aberto',
  em_analise: 'Em análise',
  resolvido: 'Resolvido',
}

/** Display order for selects/filters (most-severe first for severity). */
export const SEVERIDADE_ORDER: Severidade[] = ['critica', 'alta', 'media', 'baixa']
