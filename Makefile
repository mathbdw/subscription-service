# ============================================
# ПЕРЕМЕННЫЕ ОКРУЖЕНИЯ
# ============================================
ifneq ($(wildcard .env),)
include .env
export
else
$(warning WARNING: .env file not found! Using .env.example)
include .env.example
export
endif


BINARY_NAME := http-server
EXCLUDE_DIRS := vendor|docs|mocks|bin
TEST_PACKAGES := $(shell go list ./... | grep -v -E "($(EXCLUDE_DIRS))")
COVERAGE_FILE := coverage.out

.PHONY: setup
setup: ## Установка зависимостей (для devcontainer)
	@echo "📦 Setting up development environment..."
	go mod download
	go mod verify
	@echo "✅ Setup complete!"

.PHONY: dev
dev: fmt vet ## Запуск dev-окружения
	@echo "🚀 Dev environment ready!"

.PHONY: fmt
fmt: ## Форматирование кода
	@echo "📝 Formatting code..."
	go fmt ./...
	goimports -w .
	gci write -s standard -s default -s "prefix(github.com/mathbdw/subscription-service)" -s blank -s dot .
	@echo "✅ Formatting complete!"

.PHONY: vet
vet: ## Запуск go vet
	@echo "🔍 Running go vet..."
	go vet ./...
	@echo "✅ go vet passed!"

.PHONY: lint
lint: ## Запуск golangci-lint
	@echo "🔍 Running golangci-lint..."
	golangci-lint run --timeout=5m
	@echo "✅ Linting passed!"

.PHONY: test
test: ## Быстрые тесты (локально)
	@echo "🧪 Running tests (fast mode)..."
	CGO_ENABLED=0 go test -v $(TEST_PACKAGES)
	@echo "✅ Tests passed!"

.PHONY: test-cover
test-cover: ## Тесты с покрытием
	@echo "📊 Running tests with coverage..."
	go test -v $(TEST_PACKAGES) -race -cover -coverprofile=$(COVERAGE_FILE) -count=1
	go tool cover -func=$(COVERAGE_FILE)
	@echo "✅ Coverage check complete!"

.PHONY: test-cover-html
test-cover-html: test-cover ## Открыть HTML-отчет о покрытии
	go tool cover -html=$(COVERAGE_FILE)
	@echo "🌐 Coverage report opened in browser"

.PHONY: test-race
test-race: ## Тесты с детектором гонок
	@echo "🏎️ Running tests with race detector..."
	CGO_ENABLED=1 go test -race -v $(TEST_PACKAGES) -count=1
	@echo "✅ Race check passed!"

.PHONY: build
build: ## Сборка приложения
	@echo "🔨 Building application..."
	go mod download
	CGO_ENABLED=0 go build -tags='no_mysql no_sqlite3' -o ./bin/http-server ./cmd/app/main.go
	@echo "✅ Build complete: ./bin/http-server"

.PHONY: check-artifacts
check-artifacts: build ## Проверка артефактов сборки
	@echo "🔍 Checking build artifacts..."
	@if [ -f ./bin/http-server ]; then \
		echo "✅ BINARY FOUND: ./bin/http-server"; \
		ls -la ./bin/http-server; \
	else \
		echo "❌ ERROR: Binary not found!"; \
		exit 1; \
	fi
	@echo "✅ All artifacts checked!"

.PHONY: clean
clean: ## Очистка артефактов
	@echo "🧹 Cleaning artifacts..."
	rm -rf ./bin/
	rm -f $(COVERAGE_FILE)
	@echo "✅ Clean complete!"

.PHONY: swag
swag: ## Генерация Swagger-документации
	swag init -g internal/interfaces/http/router.go -o docs/swagger

.PHONY: pre-commit
pre-commit: fmt vet ## Быстрая проверка перед коммитом
	@echo "✅ Pre-commit checks passed!"

.PHONY: ci
ci: fmt vet lint test build check-artifacts ## Полная проверка (для CI)
	@echo "✅ All CI checks passed!"

.PHONY: all
all: pre-commit lint test-cover build ## Всё в одном (локально)
	@echo "✅ All tasks completed!"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'