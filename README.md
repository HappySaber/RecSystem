# RecSystem – система рекомендаций контента

**RecSystem** — микросервисное веб-приложение для любителей фильмов, сериалов и аниме.  
Система отслеживает взаимодействие пользователя с контентом и формирует персональные рекомендации на основе его предпочтений.  
Дополнительно доступен ИИ-ассистент на базе **LLaMA 3** (Hugging Face) для рекомендаций по произвольному запросу.

---

## ✨ Функциональность

| Сервис | Описание |
|--------|----------|
| **SSO** | Регистрация и авторизация. JWT-токены с хранением в Redis |
| **Catalog** | Управление каталогом контента. Импорт фильмов и сериалов с TMDB, аниме с AniList |
| **Recommendation** | Персональные рекомендации на основе действий пользователя. ИИ-рекомендации через Hugging Face |
| **Notification** | Отправка email-уведомлений при регистрации через SMTP |
| **API Gateway** | Единая точка входа, маршрутизация к микросервисам, JWT-авторизация |

---

## 🏗 Архитектура
```
                        ┌─────────────┐
                        │ API Gateway │  :8080
                        └──────┬──────┘
                               │ gRPC
          ┌────────────────────┼────────────────────┐
          │                    │                    │
   ┌──────▼──────┐    ┌────────▼───────┐   ┌───────▼──────┐
   │     SSO     │    │    Catalog     │   │Recommendation│
   │  :50051     │    │    :50051      │   │   :50051     │
   └──────┬──────┘    └────────┬───────┘   └──────┬───────┘
          │                    │                   │
          │              ┌─────▼──────┐            │
          └──────────────►   Kafka    ◄────────────┘
                         └─────┬──────┘
                               │
                        ┌──────▼──────┐
                        │Notification │
                        └─────────────┘
```

Каждый микросервис имеет собственную базу данных PostgreSQL.  
Межсервисное взаимодействие — через **gRPC** (синхронно) и **Kafka** (асинхронно).

---

## 🛠 Стек технологий

- **Go 1.25** — основной язык разработки
- **gRPC + Protobuf** — межсервисное взаимодействие
- **PostgreSQL 16** — основное хранилище данных
- **Redis 7** — хранение и обновление JWT-токенов
- **Apache Kafka** — асинхронные события между сервисами
- **Hugging Face (LLaMA 3)** — ИИ-рекомендации
- **Docker & Docker Compose** — контейнеризация
- **Goose** — управление миграциями БД
- **testcontainers-go** — интеграционное тестирование

---

## 🚀 Запуск приложения

### Требования

- Docker Desktop
- Make (опционально)

### Быстрый старт

**1. Клонируй репозиторий**
```bash
git clone https://github.com/username/recsystem.git
cd recsystem
```

**2. Создай файл `.env` в корне проекта**
```bash
# JWT
JWT_SECRET_KEY=your_secret_key

# SSO DB
SSO_DB_USER=postgres
SSO_DB_PASSWORD=your_password
SSO_DB_NAME=sso_rec_system
SSO_DB_PORT=5432

# Catalog DB
CATALOG_DB_USER=postgres
CATALOG_DB_PASSWORD=your_password
CATALOG_DB_NAME=catalog_rec_system
CATALOG_DB_PORT=5432

# Recommendation DB
REC_DB_USER=postgres
REC_DB_PASSWORD=your_password
REC_DB_NAME=rec_system
REC_DB_PORT=5432

# Redis
REDIS_PASSWORD=
REDIS_USER=
REDIS_DB=0
REDIS_MAX_RETRIES=5
REDIS_DIAL_TIMEOUT=5s
REDIS_TIMEOUT=3s

# External APIs
TMDB_APIKEY=your_tmdb_api_key
HF_API_KEY=your_huggingface_api_key

# SMTP
SMTP_HOST=smtp.mail.ru
SMTP_PORT=587
SMTP_USER=your@email.com
SMTP_PASS=your_smtp_password
SMTP_FROM=your@email.com
```

**3. Запусти приложение**
```bash
docker-compose up --build
```

или через Make:
```bash
make up-build
```

API Gateway доступен на `http://localhost:8080`

### Импорт контента (одноразово)

После первого запуска необходимо загрузить контент из внешних API:
```bash
docker-compose --profile import run --rm catalog-importer
```

или через Make:
```bash
make import
```

---

## 🧪 Тестирование

В проекте реализованы **юнит** и **интеграционные** тесты.  
Интеграционные тесты используют **testcontainers-go** — автоматически поднимают изолированную PostgreSQL в Docker.
```bash
# юнит тесты (быстро, без Docker)
go test ./internal/... -v

# интеграционные тесты catalog
go test ./catalog-system-microservice/tests/... -v -timeout 120s

# интеграционные тесты recommendation
go test ./recommedation-system-microservice/tests/... -v -timeout 120s
```

---

## 📁 Структура проекта
```
recsystem/
├── api-gateway/                  # HTTP Gateway
├── sso-microservice/             # Авторизация
├── catalog-system-microservice/  # Каталог контента
├── recommedation-system-microservice/ # Рекомендации
├── notification-microservice/    # Email уведомления
├── proto/                        # Protobuf схемы
├── docker-compose.yml
├── Makefile
└── .env
```

---

## 📌 Makefile команды

| Команда | Описание |
|---------|----------|
| `make up` | Запустить приложение |
| `make up-build` | Пересобрать и запустить |
| `make down` | Остановить |
| `make down-v` | Остановить и удалить данные |
| `make import` | Запустить импорт контента |
| `make logs s=catalog` | Логи конкретного сервиса |
| `make vendor` | Вендоризация зависимостей |