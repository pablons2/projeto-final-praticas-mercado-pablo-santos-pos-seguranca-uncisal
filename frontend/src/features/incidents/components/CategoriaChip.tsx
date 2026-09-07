import { CATEGORIA_LABELS, type Categoria } from '@/lib/constants'

/** Categoria chip — subtle background pill, no color coding. */
export function CategoriaChip({ value }: { value: Categoria }) {
  return (
    <span className="inline-flex items-center rounded-full bg-bg-subtle px-2.5 py-0.5 text-xs font-medium text-ink-secondary">
      {CATEGORIA_LABELS[value]}
    </span>
  )
}
