/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{ts,tsx}'],
  theme: {
    extend: {
      colors: {
        // Neutral base (slate)
        bg: {
          DEFAULT: '#FFFFFF',
          surface: '#F8FAFC',
          subtle: '#F1F5F9',
        },
        border: {
          DEFAULT: '#E2E8F0',
          hairline: '#EDF2F7',
          strong: '#CBD5E1',
        },
        ink: {
          DEFAULT: '#0F172A',
          secondary: '#475569',
          muted: '#94A3B8',
        },
        // Single accent (indigo)
        accent: {
          DEFAULT: '#4F46E5',
          hover: '#4338CA',
          subtle: '#EEF2FF',
          ring: '#6366F1',
        },
        // Semantic — incident domain (used as dots + labels, never full-bleed)
        sev: {
          baixa: '#10B981',
          media: '#F59E0B',
          alta: '#F97316',
          critica: '#EF4444',
        },
        stat: {
          aberto: '#F59E0B',
          em_analise: '#4F46E5',
          resolvido: '#10B981',
        },
        danger: {
          DEFAULT: '#EF4444',
          hover: '#DC2626',
          subtle: '#FEF2F2',
        },
      },
      fontFamily: {
        sans: [
          'Inter',
          '-apple-system',
          'BlinkMacSystemFont',
          'Segoe UI',
          'Roboto',
          'Helvetica Neue',
          'Arial',
          'sans-serif',
        ],
        mono: [
          'JetBrains Mono',
          'ui-monospace',
          'SFMono-Regular',
          'Menlo',
          'Consolas',
          'monospace',
        ],
      },
      fontSize: {
        // Fluid scale via clamp
        '2xl': ['clamp(1.25rem, 1.15rem + 0.5vw, 1.5rem)', { lineHeight: '1.2' }],
        '3xl': ['clamp(1.5rem, 1.3rem + 1vw, 1.875rem)', { lineHeight: '1.2' }],
      },
      borderRadius: {
        sm: '6px',
        md: '8px',
        lg: '12px',
      },
      boxShadow: {
        sm: '0 1px 2px 0 rgb(15 23 42 / 0.05), 0 1px 3px 0 rgb(15 23 42 / 0.06)',
        overlay: '0 4px 12px -2px rgb(15 23 42 / 0.10), 0 2px 6px -2px rgb(15 23 42 / 0.06)',
      },
      transitionDuration: {
        150: '150ms',
      },
      maxWidth: {
        prose: '65ch',
        form: '400px',
        dialog: '560px',
      },
      minHeight: {
        touch: '44px',
      },
      keyframes: {
        'fade-in': {
          from: { opacity: '0' },
          to: { opacity: '1' },
        },
        'slide-up': {
          from: { transform: 'translateY(8px)', opacity: '0' },
          to: { transform: 'translateY(0)', opacity: '1' },
        },
        'sheet-up': {
          from: { transform: 'translateY(100%)' },
          to: { transform: 'translateY(0)' },
        },
        'toast-in': {
          from: { transform: 'translateX(8px)', opacity: '0' },
          to: { transform: 'translateX(0)', opacity: '1' },
        },
      },
      animation: {
        'fade-in': 'fade-in 150ms ease-out',
        'slide-up': 'slide-up 200ms ease-out',
        'sheet-up': 'sheet-up 200ms ease-out',
        'toast-in': 'toast-in 200ms ease-out',
      },
    },
  },
  plugins: [],
}
