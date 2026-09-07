# syntax=docker/dockerfile:1

# ---------- build ----------
FROM golang:1.23-alpine AS build

WORKDIR /src

# Cache de dependências: só reexecuta quando go.mod/go.sum mudam.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

# Binário estático (sem CGO — driver SQLite puro Go), sem símbolos de debug.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# Diretório de dados pré-criado com dono nonroot (uid 65532). Um volume
# nomeado vazio herda dono/permissões deste diretório na primeira criação.
RUN mkdir -p /out/data && chown 65532:65532 /out/data

# ---------- runtime ----------
# distroless: sem shell, sem gerenciador de pacotes — superfície mínima.
# A tag :nonroot já roda como usuário sem privilégios (uid 65532).
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app
COPY --from=build /out/api /app/api
COPY --from=build --chown=65532:65532 /out/data /data

# Diretório do banco SQLite (montado como volume no compose).
VOLUME ["/data"]

ENV HTTP_ADDR=:8080 \
    DATABASE_URL=file:/data/incidenttrack.db

EXPOSE 8080
USER nonroot:nonroot

# Healthcheck sem curl/wget: o próprio binário tem o subcomando.
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=5 \
    CMD ["/app/api", "healthcheck"]

ENTRYPOINT ["/app/api"]
