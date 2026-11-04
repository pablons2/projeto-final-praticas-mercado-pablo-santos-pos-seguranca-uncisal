import { defineConfig } from '@playwright/test'

// Smoke test do IncidentTrack. Por padrão roda contra a stack local
// (`make up` → http://localhost:8080). Para outro alvo: APP_URL=https://...
//
// Navegador:
//  - CI (GitHub Actions): Chromium do próprio Playwright (`playwright install`).
//  - Local no Ubuntu 26.04: o Chromium empacotado não tem build para essa
//    versão; use o Chrome do sistema com CHROME_PATH=/usr/bin/google-chrome.
//
// Gravação: `video: 'on'` grava um vídeo de cada teste em test-results/,
// e o relatório HTML fica em playwright-report/. Ambos vão para os
// artefatos do workflow de CI como prova de execução.
const chromePath = process.env.CHROME_PATH

export default defineConfig({
  testDir: '.',
  testMatch: /.*\.spec\.ts/,
  timeout: 90_000,
  fullyParallel: false,
  workers: 1,
  reporter: [['list'], ['html', { open: 'never' }]],
  use: {
    baseURL: process.env.APP_URL ?? 'http://localhost:8080',
    video: 'on',
    screenshot: 'only-on-failure',
    trace: 'retain-on-failure',
    launchOptions: chromePath ? { executablePath: chromePath } : {},
  },
})
