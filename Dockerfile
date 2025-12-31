# syntax=docker/dockerfile:1

ARG GITHUB_PATH=github.com/mathbdw/subscription-service

FROM golang:1.24-alpine AS builder

ARG GITHUB_PATH

# Создание рабочей директории
WORKDIR /home/${GITHUB_PATH}

# Установка зависимостей
RUN apk add --no-cache --update \
    make \
    git \
    curl \
    && rm -rf /var/cache/apk/*

# Build
COPY . .
RUN make build-go


# HTTP Server

FROM alpine:latest as server
RUN apk --no-cache add ca-certificates
WORKDIR /root/

ARG GITHUB_PATH

COPY --from=builder /home/${GITHUB_PATH}/bin/http-server .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .
COPY --from=builder /home/${GITHUB_PATH}/migrations/ ./migrations

RUN chown root:root http-server

EXPOSE 8080

CMD ["./http-server"]