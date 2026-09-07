# syntax=docker/dockerfile:1

# ---------- build ----------
FROM node:22-alpine AS build

WORKDIR /app

# Instala dependências com lockfile (reprodutível). Camada cacheada
# enquanto package*.json não mudam.
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci

COPY . .

# A URL da API é embutida em build time pelo Vite. Vazio = mesma origem
# (o Nginx deste container faz o proxy de /api para o backend).
ARG VITE_API_BASE_URL=""
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

RUN npm run build

# ---------- runtime ----------
FROM nginx:1.27-alpine

# Config que serve o SPA (fallback para index.html) e faz proxy de /api.
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
    CMD wget -qO- http://127.0.0.1/healthz >/dev/null 2>&1 || exit 1
