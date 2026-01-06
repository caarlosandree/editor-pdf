.PHONY: help build run test lint format clean migrate-up migrate-down migrate-create swagger swagger-clean

# Variáveis
BINARY_NAME=server
MAIN_PATH=./cmd/server
MIGRATIONS_PATH=./migrations

help: ## Mostra esta mensagem de ajuda
	@echo "Comandos disponíveis:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Compila o projeto
	@echo "Compilando..."
	@go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

run: deps migrate-up build ## Executa o servidor (instala deps, executa migrations e compila antes)
	@echo "Executando servidor..."
	@./bin/$(BINARY_NAME)

test: ## Executa os testes
	@echo "Executando testes..."
	@go test -v ./...

test-coverage: ## Executa testes com cobertura
	@echo "Executando testes com cobertura..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

lint: ## Executa o linter
	@echo "Executando linter..."
	@golangci-lint run

format: ## Formata o código
	@echo "Formatando código..."
	@go fmt ./...
	@goimports -w .

clean: ## Limpa arquivos gerados
	@echo "Limpando..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

migrate-up: ## Executa migrations para cima
	@echo "Executando migrations..."
	@if [ -f .env.local ]; then \
		export $$(grep -v '^#' .env.local | xargs); \
	fi; \
	PATH="$$HOME/go/bin:$$PATH" migrate -path $(MIGRATIONS_PATH) -database "postgres://$${DB_USER}:$${DB_PASSWORD}@$${DB_HOST}:$${DB_PORT}/$${DB_NAME}?sslmode=$${DB_SSLMODE:-disable}" up

migrate-down: ## Reverte a última migration
	@echo "Revertendo última migration..."
	@if [ -f .env.local ]; then \
		export $$(grep -v '^#' .env.local | xargs); \
	fi; \
	PATH="$$HOME/go/bin:$$PATH" migrate -path $(MIGRATIONS_PATH) -database "postgres://$${DB_USER}:$${DB_PASSWORD}@$${DB_HOST}:$${DB_PORT}/$${DB_NAME}?sslmode=$${DB_SSLMODE:-disable}" down 1

migrate-create: ## Cria uma nova migration (use: make migrate-create NAME=nome_da_migration)
	@echo "Criando migration: $(NAME)"
	@PATH="$$HOME/go/bin:$$PATH" migrate create -ext sql -dir $(MIGRATIONS_PATH) -seq $(NAME)

deps: ## Instala dependências
	@echo "Instalando dependências..."
	@go mod download
	@go mod tidy

deps-update: ## Atualiza dependências
	@echo "Atualizando dependências..."
	@go get -u ./...
	@go mod tidy

swagger: ## Gera documentação Swagger
	@echo "Gerando documentação Swagger..."
	@swag init -g cmd/server/main.go -o cmd/server/docs

swagger-clean: ## Limpa arquivos gerados do Swagger
	@echo "Limpando documentação Swagger..."
	@rm -rf cmd/server/docs/
