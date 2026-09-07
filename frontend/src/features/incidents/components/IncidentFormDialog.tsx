import { useEffect } from 'react'
import { Dialog, DialogContent } from '@/components/ui'
import { IncidentForm } from './IncidentForm'
import { useIsDesktop } from '@/lib/useMediaQuery'
import type { Incident } from '@/lib/api-types'

interface IncidentFormDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  /** When provided, edit mode. */
  incident?: Incident
}

/**
 * Wraps the IncidentForm in a Dialog — modal on desktop, bottom-sheet on mobile.
 */
export function IncidentFormDialog({
  open,
  onOpenChange,
  incident,
}: IncidentFormDialogProps) {
  const isDesktop = useIsDesktop()
  const isEdit = Boolean(incident)

  // Close on Escape is handled by Radix; we just need to stop the form
  // from submitting when the dialog is closing.
  useEffect(() => {
    if (!open) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onOpenChange(false)
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [open, onOpenChange])

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent
        title={isEdit ? 'Editar incidente' : 'Novo incidente'}
        description={
          isEdit
            ? 'Atualize os campos abaixo.'
            : 'Registre um novo incidente de segurança.'
        }
        variant={isDesktop ? 'modal' : 'sheet'}
        hideClose={false}
      >
        <IncidentForm
          incident={incident}
          onDone={() => onOpenChange(false)}
          onCancel={() => onOpenChange(false)}
        />
      </DialogContent>
    </Dialog>
  )
}
