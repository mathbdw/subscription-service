# Subscription Service

## Overview

Сервис управления подписками с REST API для CRUDL операций и аналитики.

**Основные возможности:**
- ✅ CRUDL операции над подписками (Create, Read, Update, Delete, List)
- ✅ Подсчет суммарной стоимости подписок за период
- ✅ Фильтрация по пользователю и названию подписки
- ✅ Пагинация и сортировка
- ✅ Валидация входных данных
- ✅ Логирование операций

## Tech Stack

**Backend:**
- **Language**: Go 1.24+
- **Framework**: [Fiber](https://gofiber.io/) v2
- **Database**: PostgreSQL
- **ORM**: [sqlx](https://github.com/jmoiron/sqlx)
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **Logging**: [zerolog](https://github.com/rs/zerolog)
- **Testing**: [testify](https://github.com/stretchr/testify), [gomock](https://github.com/golang/mock)
- **Documentation**: [Swagger](https://swagger.io/) (via swaggo)

**Infrastructure:**
- **Containerization**: Docker & Docker Compose
- **Migrations**: SQL миграции

## API Documentation

- **Swagger UI**: http://localhost:8080/swagger/index.html
- **API Base URL**: http://localhost:8080/api/v1

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST   | `/subscription/create` | Создать подписку |
| GET    | `/subscription/:id` | Получить подписку по ID |
| GET    | `/subscription/list` | Список подписок с пагинацией |
| PATCH  | `/subscription/:id` | Обновить подписку |
| DELETE | `/subscription/:id` | Удалить подписку |
| GET    | `/subscription/cost` | Суммарная стоимость по подпискам |

## Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- PostgreSQL 15+ (опционально для локальной разработки)

### Local Development (без Docker)

```bash
# 1. Клонируйте репозиторий
git clone <repository-url>
cd subscription-service

# 2. Настройте окружение
cp .env.example .env
# Отредактируйте .env под ваши нужды

# 3. Установите зависимости
go mod download

# 4. Запустите PostgreSQL (если нет запущенного)
docker-compose up -d postgres

# 5. Запустите приложение
make run

# 6. Swagger документация
open http://localhost:8080/swagger/index.html
```

### Local Development (Docker)
docker-compose up --build -d