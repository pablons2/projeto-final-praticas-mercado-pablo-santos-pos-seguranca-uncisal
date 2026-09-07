import * as DropdownMenu from '@radix-ui/react-dropdown-menu'
import { ExitIcon, PersonIcon } from '@radix-ui/react-icons'
import { useNavigate } from 'react-router-dom'
import { BrandMark } from '@/components/BrandMark'
import { useAuth } from '@/context/AuthContext'
import { useLogout } from '@/features/auth/hooks/useAuthMutations'
import { cn } from '@/lib/cn'

export function TopBar() {
  const { user } = useAuth()
  const logout = useLogout()
  const navigate = useNavigate()

  if (!user) return null

  const initials = user.name
    .split(' ')
    .map((p) => p[0])
    .slice(0, 2)
    .join('')
    .toUpperCase()

  return (
    <header className="sticky top-0 z-30 flex h-14 items-center justify-between border-b border-border bg-bg/95 px-4 backdrop-blur supports-[backdrop-filter]:bg-bg/80 pt-safe">
      <BrandMark size="sm" />

      <DropdownMenu.Root>
        <DropdownMenu.Trigger asChild>
          <button
            type="button"
            aria-label="Menu do usuário"
            className={cn(
              'inline-flex items-center gap-2 rounded-full py-1 pl-1 pr-3',
              'hover:bg-bg-subtle focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-accent-ring',
            )}
          >
            <span
              className="inline-flex h-8 w-8 items-center justify-center rounded-full bg-accent text-xs font-semibold text-white"
              aria-hidden="true"
            >
              {initials || <PersonIcon />}
            </span>
            <span className="hidden max-w-[160px] truncate text-sm text-ink-secondary sm:inline">
              {user.name}
            </span>
          </button>
        </DropdownMenu.Trigger>
        <DropdownMenu.Portal>
          <DropdownMenu.Content
            align="end"
            sideOffset={6}
            className={cn(
              'z-50 min-w-[200px] rounded-md border border-border bg-bg p-1 shadow-overlay',
              'animate-fade-in',
            )}
          >
            <div className="px-2 py-2">
              <p className="text-sm font-medium text-ink">{user.name}</p>
              <p className="truncate text-xs text-ink-muted">{user.email}</p>
            </div>
            <DropdownMenu.Separator className="my-1 h-px bg-border" />
            <DropdownMenu.Item
              onSelect={() => navigate('/dashboard')}
              className={cn(
                'flex cursor-pointer items-center gap-2 rounded-sm px-2 py-2 text-sm text-ink-secondary outline-none',
                'data-[highlighted]:bg-bg-subtle data-[highlighted]:text-ink',
              )}
            >
              <PersonIcon /> Meus incidentes
            </DropdownMenu.Item>
            <DropdownMenu.Item
              onSelect={() => logout.mutate()}
              disabled={logout.isPending}
              className={cn(
                'flex cursor-pointer items-center gap-2 rounded-sm px-2 py-2 text-sm text-danger outline-none',
                'data-[highlighted]:bg-danger-subtle',
              )}
            >
              <ExitIcon /> Sair
            </DropdownMenu.Item>
          </DropdownMenu.Content>
        </DropdownMenu.Portal>
      </DropdownMenu.Root>
    </header>
  )
}
