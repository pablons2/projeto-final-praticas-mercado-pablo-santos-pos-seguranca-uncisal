import { AlertDialog, AlertDialogContent } from '@/components/ui'
import { useDeleteIncident } from '../hooks/useIncidents'
import type { Incident } from '@/lib/api-types'

interface DeleteConfirmDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  incident: Incident | null
}

export function DeleteConfirmDialog({
  open,
  onOpenChange,
  incident,
}: DeleteConfirmDialogProps) {
  const del = useDeleteIncident()

  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent
        title="Excluir incidente"
        description={
          incident ? (
            <>
              Excluir <strong className="text-ink">«{incident.titulo}»</strong>? Esta ação
              não pode ser desfeita.
            </>
          ) : (
            'Esta ação não pode ser desfeita.'
          )
        }
        confirmLabel="Excluir"
        destructive
        loading={del.isPending}
        onConfirm={() => {
          if (!incident) return
          del.mutate(incident.id, {
            onSettled: () => onOpenChange(false),
          })
        }}
      />
    </AlertDialog>
  )
}
