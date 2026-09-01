# ifneq ($(wildcard .env),)
# include .env
# export
# else
# $(warning WARNING: .env file not found! Using .env.example)
# include .env.example
# export
# endif

.PHONY: run
run:
	go run cmd/app/main.go


EXCLUDE_DIRS = vendor|docs|mocks
TEST_PACKAGES = $(shell go list ./... | grep -v -E "($(EXCLUDE_DIRS))")
COVERAGE_FILE = coverage.out

# ============================================
.PHONY: setup
setup:
	@echo "📦 Setting up development environment..."
	go mod download
	go mod verify
	@echo "✅ Setup complete!"

.PHONY: dev
dev:
	@echo "🚀 Starting development environment..."
# Можно поднять тестовые данные или проверить БД
	@echo "✅ Dev environment ready!"

# ============================================
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .
	gci write -s standard -s default -s "prefix(github.com/mathbdw/subscription-service)" .
	@echo "✅ Formatting complete!"

# ============================================
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...
	@echo "✅ go vet passed!"

.PHONY: lint
lint:
	@echo "🔍 Running golangci-lint..."
	golangci-lint run --timeout=5m
	@echo "✅ Linting passed!"

# ============================================
# 1. Быстрые тесты (для локальной разработки)
.PHONY: test
test:
	@echo "🧪 Running tests (fast mode)..."
	@echo "📦 Packages: $(TEST_PACKAGES)" | tr ' ' '\n'
	CGO_ENABLED=0 go test -v $(TEST_PACKAGES)
	@echo "✅ Tests passed!"

# 2. Тесты с покрытием (без -race)
.PHONY: test-cover
test-cover:
	@echo "📊 Running tests with coverage..."
	CGO_ENABLED=0 go test -v $(TEST_PACKAGES) -cover -coverprofile=$(COVERAGE_FILE) -count=1
	@echo "📈 Coverage report saved to $(COVERAGE_FILE)"
	go tool cover -func=$(COVERAGE_FILE)
	@echo "✅ Coverage check complete!"

# 3. Открыть HTML-отчет
.PHONY: test-cover-html
test-cover-html: test-cover
	go tool cover -html=$(COVERAGE_FILE)
	@echo "🌐 Coverage report opened in browser"

# 4. Тесты с детектором гонок (отдельно)
.PHONY: test-race
test-race:
	@echo "🏎️ Running tests with race detector..."
	CGO_ENABLED=1 go test -race -v $(TEST_PACKAGES) -count=1
	@echo "✅ Race check passed!"

# ============================================

.PHONY: build-go
build-go: .build

.build:
	go mod download && CGO_ENABLED=0  go build \
		-tags='no_mysql no_sqlite3' \
		-o ./bin/http-server$(shell go env GOEXE) ./cmd/app/main.go


.PHONY: swag-v1
swag-v1:
	swag init -g internal/interfaces/http/router.go -o docs/swagger

# ============================================
# СБОРНЫЕ ЦЕЛИ (ДЛЯ КОНКРЕТНЫХ СЦЕНАРИЕВ)
# ============================================

# ✅ Для pre-commit хука (быстро)
.PHONY: pre-commit
pre-commit: fmt vet
	@echo "✅ Pre-commit checks passed!"

# ✅ Для CI / перед PR (полная проверка)
.PHONY: ci
ci: fmt vet lint test
	@echo "✅ All CI checks passed!"