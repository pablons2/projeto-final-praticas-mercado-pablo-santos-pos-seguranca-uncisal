import { clsx, type ClassValue } from 'clsx'

/** Tailwind-aware className combiner. */
export function cn(...parts: ClassValue[]): string {
  return clsx(parts)
}
