import { z } from 'zod'
import { CATEGORIAS, SEVERIDADES, STATUS } from '@/lib/constants'

/**
 * Incident validation — mirrors the backend entity (docs/plano.md §2).
 * Client-side validation is for fast feedback only; the server is
 * authoritative.
 */
export const incidentSchema = z.object({
  titulo: z
    .string()
    .min(3, 'Título muito curto (mín. 3)')
    .max(120, 'Título muito longo (máx. 120)')
    .transform((v) => v.trim()),
  descricao: z
    .string()
    .min(10, 'Descrição muito curta (mín. 10)')
    .max(1000, 'Descrição muito longa (máx. 1000)')
    .transform((v) => v.trim()),
  categoria: z.enum(CATEGORIAS, { message: 'Selecione uma categoria' }),
  severidade: z.enum(SEVERIDADES, { message: 'Selecione a severidade' }),
  status: z.enum(STATUS, { message: 'Selecione o status' }),
})

export type IncidentValues = z.infer<typeof incidentSchema>

/** Filters for the dashboard list — all optional. */
export const incidentFiltersSchema = z.object({
  search: z.string().default(''),
  categoria: z.enum(CATEGORIAS).optional().or(z.literal('')),
  severidade: z.enum(SEVERIDADES).optional().or(z.literal('')),
  status: z.enum(STATUS).optional().or(z.literal('')),
})

export type IncidentFilters = z.infer<typeof incidentFiltersSchema>
