.DEFAULT_GOAL := help
BIN := bin/api
PKG := ./...

.PHONY: help
help: ## Lista os alvos disponíveis
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

.PHONY: tidy
tidy: ## Resolve dependências e gera go.sum
	go mod tidy

.PHONY: run
run: ## Sobe a API localmente (usa .env)
	go run ./cmd/api

.PHONY: build
build: ## Compila o binário estático em bin/api
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o $(BIN) ./cmd/api

.PHONY: test
test: ## Roda os testes
	go test -race -count=1 $(PKG)

.PHONY: cover
cover: ## Testes com relatório de cobertura
	go test -race -count=1 -coverprofile=coverage.out $(PKG)
	go tool cover -func=coverage.out | tail -n 1

.PHONY: vet
vet: ## go vet
	go vet $(PKG)

.PHONY: fmt
fmt: ## Formata o código
	gofmt -w -s .

.PHONY: check
check: fmt vet test ## fmt + vet + test
