import { useMutation } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { http, ApiError } from '@/lib/http'
import { useAuth } from '@/context/AuthContext'
import { useToast } from '@/components/ui'
import type { LoginPayload, RegisterPayload, User } from '@/lib/api-types'
import type { LoginValues, RegisterValues } from '../schemas/auth.schema'

/** Login mutation — on success, set the user from /me (cookie is now set). */
export function useLogin() {
  const { setUser } = useAuth()
  const { toast } = useToast()
  const navigate = useNavigate()

  return useMutation({
    mutationFn: async (values: LoginValues) => {
      const payload: LoginPayload = {
        email: values.email.trim().toLowerCase(),
        password: values.password,
      }
      await http.post('/api/auth/login', payload)
      // Cookie is set; fetch the user once. /me returns the user directly.
      return http.get<User>('/api/auth/me')
    },
    onSuccess: (user: User) => {
      setUser(user)
      navigate('/dashboard', { replace: true })
    },
    onError: (err) => {
      const message =
        err instanceof ApiError ? err.message : 'Falha ao entrar. Tente novamente.'
      toast({ title: 'Não foi possível entrar', description: message, variant: 'error' })
    },
  })
}

/** Register mutation — on success, log the user in directly. */
export function useRegister() {
  const { setUser } = useAuth()
  const { toast } = useToast()
  const navigate = useNavigate()

  return useMutation({
    mutationFn: async (values: RegisterValues) => {
      const payload: RegisterPayload = {
        name: values.name,
        email: values.email.trim().toLowerCase(),
        password: values.password,
      }
      await http.post('/api/auth/register', payload)
      // Auto-login after register — cookie is set by the register endpoint
      // if the backend issues it, otherwise call /login. We try /me first;
      // if it 401s, we explicitly log in.
      try {
        return await http.get<User>('/api/auth/me')
      } catch (err) {
        if (err instanceof ApiError && err.status === 401) {
          await http.post('/api/auth/login', {
            email: payload.email,
            password: payload.password,
          })
          return await http.get<User>('/api/auth/me')
        }
        throw err
      }
    },
    onSuccess: (user: User) => {
      setUser(user)
      toast({
        title: 'Conta criada',
        description: 'Bem-vindo ao IncidentTrack.',
        variant: 'success',
      })
      navigate('/dashboard', { replace: true })
    },
    onError: (err) => {
      const message =
        err instanceof ApiError ? err.message : 'Falha ao registrar. Tente novamente.'
      toast({ title: 'Não foi possível registrar', description: message, variant: 'error' })
    },
  })
}

/** Logout — clears local state and calls /logout. */
export function useLogout() {
  const { logout } = useAuth()
  const navigate = useNavigate()
  const { toast } = useToast()

  return useMutation({
    mutationFn: () => http.post('/api/auth/logout'),
    onMutate: () => {
      // Optimistic: clear local state immediately for snappy UX.
      // The AuthContext.logout already clears user; we just navigate.
    },
    onSuccess: () => {
      toast({ title: 'Sessão encerrada', variant: 'info' })
      navigate('/login', { replace: true })
    },
    onError: () => {
      // Even on error, AuthContext.logout clears local state.
      navigate('/login', { replace: true })
    },
    onSettled: () => {
      void logout()
    },
  })
}
