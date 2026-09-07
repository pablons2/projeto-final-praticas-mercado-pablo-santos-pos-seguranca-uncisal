import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import { http, ApiError } from '@/lib/http'
import type { User } from '@/lib/api-types'

interface AuthContextValue {
  user: User | null
  /** null while the initial /me check is in flight. */
  bootstrapped: boolean
  /** Sets the user after a successful login/register without a second round-trip. */
  setUser: (user: User | null) => void
  /** Clears local state and calls /logout. */
  logout: () => Promise<void>
}

const AuthContext = createContext<AuthContextValue | null>(null)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [bootstrapped, setBootstrapped] = useState(false)

  // Single bootstrap call on mount: ask the server who we are via the
  // HttpOnly cookie. 401 → not authenticated (expected, not an error).
  useEffect(() => {
    let active = true
    ;(async () => {
      try {
        const data = await http.get<User>('/api/auth/me')
        if (active) setUser(data)
      } catch (err) {
        if (active && err instanceof ApiError && err.status !== 401) {
          // Network blip etc. — leave user null, routes will handle.
        }
      } finally {
        if (active) setBootstrapped(true)
      }
    })()
    return () => {
      active = false
    }
  }, [])

  const logout = useCallback(async () => {
    try {
      await http.post('/api/auth/logout')
    } catch {
      // Even if the server call fails, clear local state so the UI
      // doesn't stay stuck on an authenticated view.
    } finally {
      setUser(null)
    }
  }, [])

  const value = useMemo<AuthContextValue>(
    () => ({ user, bootstrapped, setUser, logout }),
    [user, bootstrapped, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within <AuthProvider>')
  return ctx
}
