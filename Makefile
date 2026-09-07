# IncidentTrack — atalhos de orquestração.
# Use Docker (padrão) ou Podman:  make COMPOSE="podman compose" up

COMPOSE ?= docker compose
.DEFAULT_GOAL := help

.PHONY: help
help: ## Lista os alvos
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

.PHONY: env
env: ## Cria .env a partir do exemplo (não sobrescreve), com JWT_SECRET gerado
	@if [ -f .env ]; then \
		echo "✓ .env já existe (mantido)"; \
	else \
		sed 's#^JWT_SECRET=.*#JWT_SECRET='"$$(openssl rand -base64 48)"'#' .env.example > .env; \
		echo "✓ .env criado com JWT_SECRET aleatório"; \
	fi

.PHONY: up
up: env ## Sobe a stack (build se necessário) — http://localhost:8080
	$(COMPOSE) up --build -d
	@echo "→ http://localhost:$${APP_PORT:-8080}"

.PHONY: up-fg
up-fg: env ## Sobe a stack em foreground (logs no terminal)
	$(COMPOSE) up --build

.PHONY: down
down: ## Para e remove os containers
	$(COMPOSE) down

.PHONY: clean
clean: ## Para tudo e APAGA o volume do banco
	$(COMPOSE) down -v

.PHONY: logs
logs: ## Segue os logs de todos os serviços
	$(COMPOSE) logs -f

.PHONY: ps
ps: ## Status dos serviços
	$(COMPOSE) ps

.PHONY: rebuild
rebuild: ## Recompila as imagens sem cache
	$(COMPOSE) build --no-cache

.PHONY: backend-check
backend-check: ## fmt + vet + test do backend (precisa de Go local)
	$(MAKE) -C backend check
