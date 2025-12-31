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

# Исключаемые директории
EXCLUDE_DIRS = vendor|docs|mocks

# Получаем пакеты для тестирования
TEST_PACKAGES = $(shell go list ./... | grep -v -E "($(EXCLUDE_DIRS))")

.PHONY: test
test:
	@echo "📦 Testing packages:"
	@echo "$(TEST_PACKAGES)" | tr ' ' '\n'
	@echo ""
	go test -v $(TEST_PACKAGES) -cover -coverprofile=./coverage.out

.PHONY: test-cover
test-cover: test
	go tool cover -html=./coverage.out
	@echo "✅ Coverage report generated: coverage.html"

.PHONY: test-short
test-short:
	go test -v $(TEST_PACKAGES) -short -cover

.PHONY: test-race
test-race:
	go test -v $(TEST_PACKAGES) -race -cover

.PHONY: build-go
build-go: .build

.build:
	go mod download && CGO_ENABLED=0  go build \
		-tags='no_mysql no_sqlite3' \
		-o ./bin/http-server$(shell go env GOEXE) ./cmd/app/main.go


.PHONY: swag-v1
swag-v1:
	swag init -g internal/interfaces/http/router.go -o docs/swagger