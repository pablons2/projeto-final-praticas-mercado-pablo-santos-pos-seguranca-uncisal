import { fileURLToPath, URL } from 'node:url'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'
import { visualizer } from 'rollup-plugin-visualizer'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    react(),
    // Bundle report — emitted on every build, opened only with --mode analyze.
    visualizer({
      filename: 'dist/stats.html',
      gzipSize: true,
      open: process.env.ANALYZE === 'true',
    }),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  build: {
    target: 'es2020',
    sourcemap: false,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules')) {
            if (id.includes('react-router')) return 'router'
            if (id.includes('react-dom') || id.includes('/react/'))
              return 'react-vendor'
            if (id.includes('@radix-ui')) return 'radix'
            if (id.includes('@tanstack')) return 'query'
            if (id.includes('zod') || id.includes('react-hook-form'))
              return 'forms'
          }
        },
      },
    },
  },
})
