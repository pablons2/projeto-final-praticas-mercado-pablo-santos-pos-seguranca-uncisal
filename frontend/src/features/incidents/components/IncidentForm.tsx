import { useForm, Controller, type Control } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { incidentSchema, type IncidentValues } from '../schemas/incident.schema'
import {
  CATEGORIAS,
  CATEGORIA_LABELS,
  SEVERIDADES,
  SEVERIDADE_LABELS,
  STATUS,
  STATUS_LABELS,
} from '@/lib/constants'
import { Button, Field, Input, Textarea, Select, type SelectOption } from '@/components/ui'
import { useCreateIncident, useUpdateIncident } from '../hooks/useIncidents'
import type { Incident } from '@/lib/api-types'

const CATEGORIA_OPTIONS: SelectOption[] = CATEGORIAS.map((c) => ({
  value: c,
  label: CATEGORIA_LABELS[c],
}))
const SEVERIDADE_OPTIONS: SelectOption[] = SEVERIDADES.map((s) => ({
  value: s,
  label: SEVERIDADE_LABELS[s],
}))
const STATUS_OPTIONS: SelectOption[] = STATUS.map((s) => ({
  value: s,
  label: STATUS_LABELS[s],
}))

const TITULO_MAX = 120
const DESCRICAO_MAX = 1000

interface IncidentFormProps {
  /** When provided, the form is in edit mode. */
  incident?: Incident
  onDone: () => void
  onCancel: () => void;
}

export function IncidentForm({ incident, onDone, onCancel }: IncidentFormProps) {
  const isEdit = Boolean(incident)
  const create = useCreateIncident()
  const update = useUpdateIncident()

  const {
    register,
    handleSubmit,
    setError,
    watch,
    control,
    formState: { errors, isSubmitting },
  } = useForm<IncidentValues>({
    resolver: zodResolver(incidentSchema),
    defaultValues: incident
      ? {
          titulo: incident.titulo,
          descricao: incident.descricao,
          categoria: incident.categoria as IncidentValues['categoria'],
          severidade: incident.severidade as IncidentValues['severidade'],
          status: incident.status as IncidentValues['status'],
        }
      : {
          titulo: '',
          descricao: '',
          categoria: undefined as unknown as IncidentValues['categoria'],
          severidade: undefined as unknown as IncidentValues['severidade'],
          status: 'aberto',
        },
  })

  const tituloLen = (watch('titulo') ?? '').length
  const descricaoLen = (watch('descricao') ?? '').length

  const onSubmit = handleSubmit(async (values) => {
    try {
      if (isEdit && incident) {
        await update.mutateAsync({ id: incident.id, values })
      } else {
        await create.mutateAsync(values)
      }
      onDone()
    } catch (err) {
      if (err && typeof err === 'object' && 'fields' in err) {
        const fields = (err as { fields?: Record<string, string[]> }).fields
        if (fields) {
          for (const [key, msgs] of Object.entries(fields)) {
            if (
              key === 'titulo' ||
              key === 'descricao' ||
              key === 'categoria' ||
              key === 'severidade' ||
              key === 'status'
            ) {
              setError(key, { message: msgs[0] })
            }
          }
        }
      }
    }
  })

  const pending = isSubmitting || create.isPending || update.isPending

  return (
    <form onSubmit={onSubmit} noValidate className="flex flex-col gap-4">
      <Field
        label="Título"
        required
        error={errors.titulo?.message}
        footer={
          <span className={tituloLen > TITULO_MAX ? 'text-danger' : undefined}>
            {tituloLen}/{TITULO_MAX}
          </span>
        }
      >
        {(props) => (
          <Input
            {...props}
            {...register('titulo')}
            placeholder="Ex.: E-mail de phishing recebido"
            maxLength={TITULO_MAX}
            disabled={pending}
          />
        )}
      </Field>

      <Field
        label="Descrição"
        required
        error={errors.descricao?.message}
        footer={
          <span className={descricaoLen > DESCRICAO_MAX ? 'text-danger' : undefined}>
            {descricaoLen}/{DESCRICAO_MAX}
          </span>
        }
      >
        {(props) => (
          <Textarea
            {...props}
            {...register('descricao')}
            placeholder="Descreva o que aconteceu, quando e o impacto estimado."
            maxLength={DESCRICAO_MAX}
            rows={5}
            disabled={pending}
          />
        )}
      </Field>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
        <Field label="Categoria" required error={errors.categoria?.message}>
          {(props) => (
            <ControllerSelect
              control={control}
              name="categoria"
              id={props.id}
              ariaLabel="Categoria"
              invalid={Boolean(errors.categoria)}
              options={CATEGORIA_OPTIONS}
              disabled={pending}
            />
          )}
        </Field>
        <Field label="Severidade" required error={errors.severidade?.message}>
          {(props) => (
            <ControllerSelect
              control={control}
              name="severidade"
              id={props.id}
              ariaLabel="Severidade"
              invalid={Boolean(errors.severidade?.message)}
              options={SEVERIDADE_OPTIONS}
              disabled={pending}
            />
          )}
        </Field>
        <Field label="Status" required error={errors.status?.message}>
          {(props) => (
            <ControllerSelect
              control={control}
              name="status"
              id={props.id}
              ariaLabel="Status"
              invalid={Boolean(errors.status?.message)}
              options={STATUS_OPTIONS}
              disabled={pending}
            />
          )}
        </Field>
      </div>

      <div className="mt-2 flex items-center justify-end gap-2">
        <Button type="button" variant="secondary" onClick={onCancel} disabled={pending}>
          Cancelar
        </Button>
        <Button type="submit" loading={pending}>
          {isEdit ? 'Salvar alterações' : 'Criar incidente'}
        </Button>
      </div>
    </form>
  )
}

/** Bridges react-hook-form's Controller with our Select. */
function ControllerSelect({
  control,
  name,
  id,
  ariaLabel,
  invalid,
  options,
  disabled,
}: {
  control: Control<IncidentValues>
  name: keyof IncidentValues
  id: string
  ariaLabel: string
  invalid: boolean
  options: SelectOption[]
  disabled?: boolean
}) {
  return (
    <Controller
      control={control}
      name={name}
      render={({ field }) => (
        <Select
          id={id}
          ariaLabel={ariaLabel}
          invalid={invalid}
          options={options}
          value={field.value ?? ''}
          onChange={field.onChange}
          disabled={disabled}
        />
      )}
    />
  )
}
