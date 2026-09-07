import { Link } from 'react-router-dom'
import { AuthLayout } from '@/pages/AuthLayout'
import { LoginForm } from '@/features/auth/components/LoginForm'

export function LoginPage() {
  return (
    <AuthLayout
      title="Entrar"
      subtitle="Acesse seu painel de incidentes."
      footer={
        <>
          Não tem conta?{' '}
          <Link
            to="/register"
            className="font-medium text-accent hover:text-accent-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring"
          >
            Criar conta
          </Link>
        </>
      }
    >
      <LoginForm />
    </AuthLayout>
  )
}
