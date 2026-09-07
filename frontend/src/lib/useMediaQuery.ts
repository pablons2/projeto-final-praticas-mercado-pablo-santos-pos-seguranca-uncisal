import { useEffect, useState } from 'react'

/**
 * SSR-safe media query hook. Returns false during the first paint to
 * avoid layout flash, then resolves to the real match on mount.
 *
 * Used to switch the Dashboard list between table (>=1024px) and cards,
 * and the incident form between modal and bottom-sheet.
 */
export function useMediaQuery(query: string): boolean {
  const [matches, setMatches] = useState(false)

  useEffect(() => {
    const mql = window.matchMedia(query)
    const update = () => setMatches(mql.matches)
    update()
    mql.addEventListener('change', update)
    return () => mql.removeEventListener('change', update)
  }, [query])

  return matches
}

/** Convenience: true on >=1024px (lg breakpoint). */
export function useIsDesktop(): boolean {
  return useMediaQuery('(min-width: 1024px)')
}
