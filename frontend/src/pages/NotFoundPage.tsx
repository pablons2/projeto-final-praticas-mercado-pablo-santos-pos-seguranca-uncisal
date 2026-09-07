import { Link } from 'react-router-dom'
import { BrandMark } from '@/components/BrandMark'

export function NotFoundPage() {
  return (
    <div className="flex min-h-dvh flex-col items-center justify-center gap-6 bg-bg-surface px-4 text-center">
      <BrandMark />
      <div>
        <p className="text-5xl font-semibold tracking-tight text-ink">404</p>
        <p className="mt-2 text-sm text-ink-secondary">
          A página que você procura não existe.
        </p>
      </div>
      <Link
        to="/"
        className="inline-flex h-11 items-center justify-center rounded-sm bg-accent px-4 text-sm font-medium text-white hover:bg-accent-hover focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring focus-visible:ring-offset-2"
      >
        Voltar ao início
      </Link>
    </div>
  )
}
