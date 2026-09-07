import { type ReactNode } from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'
import { Spinner } from '@/components/ui'

/**
 * Gate for auth-only routes (login, register). If the user is already
 * authenticated, bounce them to /dashboard so they don't see the login
 * form again after a fresh page load with a valid cookie.
 */
export function PublicOnlyRoute({ children }: { children: ReactNode }) {
  const { user, bootstrapped } = useAuth()

  if (!bootstrapped) {
    return (
      <div className="flex min-h-dvh items-center justify-center">
        <Spinner className="h-6 w-6 text-ink-muted" />
      </div>
    )
  }

  if (user) return <Navigate to="/dashboard" replace />

  return <>{children}</>
}
