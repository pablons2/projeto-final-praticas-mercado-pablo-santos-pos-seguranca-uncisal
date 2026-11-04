import { defineConfig } from '@playwright/test'

// Smoke test do IncidentTrack. Por padrão roda contra a stack local
// (`make up` → http://localhost:8080). Para outro alvo: APP_URL=https://...
// O navegador usa o Google Chrome do sistema: o Chromium empacotado pelo
// Playwright não tem build para o Ubuntu 26.04. Caminho via CHROME_PATH.
export default defineConfig({
  testDir: '.',
  testMatch: /.*\.spec\.ts/,
  timeout: 90_000,
  fullyParallel: false,
  workers: 1,
  reporter: [['list']],
  use: {
    baseURL: process.env.APP_URL ?? 'http://localhost:8080',
    screenshot: 'only-on-failure',
    trace: 'off',
    launchOptions: {
      executablePath: process.env.CHROME_PATH ?? '/usr/bin/google-chrome',
    },
  },
})
