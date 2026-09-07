import { QueryClient } from '@tanstack/react-query'

/**
 * React Query client — conservative defaults for a small, rate-limited API.
 * - staleTime 30s: avoid hammering /api/incidents on every focus.
 * - retry 1: don't amplify load on a rate-limited (429) endpoint.
 * - refetchOnWindowFocus: yes (cheap, catches stale data after app switch).
 */
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 30_000,
      gcTime: 5 * 60_000,
      retry: 1,
      refetchOnWindowFocus: true,
    },
    mutations: {
      retry: 0,
    },
  },
})
