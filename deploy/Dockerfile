# ============================================
# Stage 1: Build (Сборка)
# ============================================
FROM golang:1.26-alpine AS builder

# Устанавливаем только необходимое для сборки
RUN apk add --no-cache git make

WORKDIR /app

# Копируем манифесты зависимостей ПЕРВЫМИ
# Это критично для кэширования слоев Docker
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем бинарник
# CGO_ENABLED=0 - статическая линковка (важно для alpine/distroless)
# -tags='no_mysql no_sqlite3' - твои теги сборки из Makefile
# -ldflags="-s -w" - убираем отладочную информацию
RUN CGO_ENABLED=0 GOOS=linux go build \
    -tags='no_mysql no_sqlite3' \
    -ldflags="-s -w" \
    -o /app/bin/http-server \
    ./cmd/app/main.go

# ============================================
# Stage 2: Runtime (Запуск)
# ============================================
# Используем distroless - минимальный образ без shell
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /

# Копируем файлы
COPY --from=builder /app/bin/http-server /http-server
COPY --from=builder /app/config.yml /config.yml
COPY --from=builder /app/migrations /migrations

USER nonroot:nonroot
EXPOSE 8080

ENTRYPOINT ["/http-server"]