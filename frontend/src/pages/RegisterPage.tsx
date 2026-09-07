import { Link } from 'react-router-dom'
import { AuthLayout } from '@/pages/AuthLayout'
import { RegisterForm } from '@/features/auth/components/RegisterForm'

export function RegisterPage() {
  return (
    <AuthLayout
      title="Criar conta"
      subtitle="Comece a registrar incidentes em segundos."
      footer={
        <>
          Já tem conta?{' '}
          <Link
            to="/login"
            className="font-medium text-accent hover:text-accent-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring"
          >
            Entrar
          </Link>
        </>
      }
    >
      <RegisterForm />
    </AuthLayout>
  )
}
