/**
 * Display formatters — pure functions, no I18n lib needed for this scope.
 * All dates are formatted relative to the user's locale (pt-BR default).
 */

/** "há 2h", "há 5min", "agora mesmo" — concise relative time. */
export function relativeTime(iso: string, now: Date = new Date()): string {
  const then = new Date(iso)
  const diffMs = now.getTime() - then.getTime()
  const sec = Math.round(diffMs / 1000)

  if (sec < 45) return 'agora mesmo'
  const min = Math.round(sec / 60)
  if (min < 60) return `há ${min}min`
  const hr = Math.round(min / 60)
  if (hr < 24) return `há ${hr}h`
  const day = Math.round(hr / 24)
  if (day < 30) return `há ${day}d`
  const month = Math.round(day / 30)
  if (month < 12) return `há ${month} mês${month > 1 ? 'es' : ''}`
  const year = Math.round(month / 12)
  return `há ${year} ano${year > 1 ? 's' : ''}`
}

/** "07/09/2026 14:32" — absolute fallback for tooltips/details. */
export function absoluteTime(iso: string): string {
  return new Date(iso).toLocaleString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/** Short id chip: "a3f1" from a UUID. */
export function shortId(id: string): string {
  return id.replace(/-/g, '').slice(0, 4)
}
