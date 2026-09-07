import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { http, ApiError } from '@/lib/http'
import { useToast } from '@/components/ui'
import type {
  Incident,
  IncidentCreatePayload,
  IncidentUpdatePayload,
} from '@/lib/api-types'
import type { IncidentValues } from '../schemas/incident.schema'

const QK = ['incidents'] as const

/** List query — the backend already filters by owner_id via the JWT. */
export function useIncidents() {
  return useQuery({
    queryKey: QK,
    queryFn: () => http.get<Incident[]>('/api/incidents'),
    select: (data) => data ?? [],
  })
}

function toPayload(values: IncidentValues): IncidentCreatePayload {
  return {
    titulo: values.titulo,
    descricao: values.descricao,
    categoria: values.categoria,
    severidade: values.severidade,
    status: values.status,
  }
}

export function useCreateIncident() {
  const qc = useQueryClient()
  const { toast } = useToast()

  return useMutation({
    mutationFn: (values: IncidentValues) =>
      http.post<Incident>('/api/incidents', toPayload(values)),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: QK })
      toast({ title: 'Incidente criado', variant: 'success' })
    },
    onError: (err) => {
      toast({
        title: 'Não foi possível criar',
        description: err instanceof ApiError ? err.message : undefined,
        variant: 'error',
      })
    },
  })
}

export function useUpdateIncident() {
  const qc = useQueryClient()
  const { toast } = useToast()

  return useMutation({
    mutationFn: ({ id, values }: { id: string; values: IncidentValues }) =>
      http.put<Incident>(`/api/incidents/${id}`, toPayload(values) as IncidentUpdatePayload),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: QK })
      toast({ title: 'Incidente atualizado', variant: 'success' })
    },
    onError: (err) => {
      toast({
        title: 'Não foi possível atualizar',
        description: err instanceof ApiError ? err.message : undefined,
        variant: 'error',
      })
    },
  })
}

export function useDeleteIncident() {
  const qc = useQueryClient()
  const { toast } = useToast()

  return useMutation({
    mutationFn: (id: string) => http.delete(`/api/incidents/${id}`),
    onMutate: async (id) => {
      // Optimistic delete — remove from cache immediately, rollback on error.
      await qc.cancelQueries({ queryKey: QK })
      const prev = qc.getQueryData<Incident[]>(QK)
      qc.setQueryData<Incident[]>(QK, (old) => (old ?? []).filter((i) => i.id !== id))
      return { prev }
    },
    onError: (_err, _id, ctx) => {
      if (ctx?.prev) qc.setQueryData(QK, ctx.prev)
      toast({
        title: 'Não foi possível excluir',
        description: 'Tente novamente.',
        variant: 'error',
      })
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: QK })
      toast({ title: 'Incidente excluído', variant: 'info' })
    },
  })
}
