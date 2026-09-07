# IncidentTrack — Frontend (React SPA)

Single-page application for the IncidentTrack security-incident registry.
Aligned with `docs/plano.md` §8. Stack: **Vite + React 18 + TypeScript**,
**Tailwind CSS**, **Radix Primitives** (custom-styled), **react-query**,
**react-hook-form + zod**.

## Run locally

```bash
npm install
cp .env.example .env      # adjust VITE_API_BASE_URL if backend is on another port
npm run dev               # http://localhost:5173
```

## Scripts

| Command          | Purpose                                  |
| ---------------- | ---------------------------------------- |
| `npm run dev`    | Dev server with HMR                      |
| `npm run build`  | Type-check + production build to `dist/` |
| `npm run preview`| Serve the production build locally       |
| `npm run analyze`| Build + open bundle-size report          |
| `npm run lint`   | oxlint                                   |

## Architecture

```
src/
├── components/ui/     Design-system primitives (Radix + Tailwind)
├── components/        Shared app components (TopBar, BrandMark)
├── features/auth/     Auth domain (schemas, hooks, forms)
├── features/incidents/ Incident domain (schemas, hooks, list, form, dialogs)
├── context/           AuthContext (user state via /api/auth/me)
├── routes/            ProtectedRoute / PublicOnlyRoute guards
├── pages/             Route-level page components (lazy-loaded)
├── lib/               http client, query client, formatters, constants, hooks
└── styles/            tokens.css (CSS vars for Radix)
```

## Design system

- **Light-minimal** theme. Tokens in `tailwind.config.js` + `src/styles/tokens.css`.
- One accent (indigo `#4F46E5`), neutral slate base, hairline borders.
- Semantic colors (severity/status) used only as 8px dots + text labels —
  color is never the only signal (WCAG 2.2 AA).
- Fluid type scale via `clamp`; 4px spacing grid; 44px min touch targets.

## Performance

- Route-level code splitting (`React.lazy`) — login first paint ~71KB gzip.
- Manual chunks: `react-vendor`, `router`, `radix`, `query`, `forms`.
- Budget: first-load gzip < 180KB. Run `npm run analyze` to inspect.
- No icon fonts, no runtime CSS, no Google Fonts CDN request.

## Security-relevant behaviours (OWASP, see docs/plano.md §7)

- **No token in JS storage.** Auth relies entirely on the `HttpOnly` cookie;
  `AuthContext` only knows "is there a user" via `GET /api/auth/me`.
- **XSS surface minimized.** No `dangerouslySetInnerHTML`; incident content
  rendered as plain text (React escapes by default).
- **CSRF.** Every mutating request carries `X-Requested-With: XMLHttpRequest`.
- **No secret leakage.** Generic error messages; raw server traces never
  displayed. Inputs trimmed client-side, but server validation is authoritative.

## Accessibility

- WCAG 2.2 AA: 4.5:1 contrast, full keyboard nav, visible focus rings.
- Radix primitives handle focus trap, `aria-*`, escape, scroll-lock.
- `prefers-reduced-motion` disables animations.
- Skip-to-content link on the dashboard.

## Responsive strategy

- Mobile-first. Dashboard list renders as **cards** < 1024px and a **table**
  >= 1024px (same hooks/schemas, presentation-only components).
- Incident form is a **bottom sheet** on mobile, **centered modal** on desktop.
