# RecSystem - система рекомендаций контента

**RecSystem** - микросервисное веб-приложение для любителей фильмов, сериалов и аниме.  
Система отслеживает взаимодействие пользователя с контентом и формирует персональные рекомендации на основе его предпочтений.  
Дополнительно доступен ИИ-ассистент на базе **LLaMA 3** (Hugging Face) для рекомендаций по произвольному запросу.

---

## Функциональность

| Сервис | Описание |
|--------|----------|
| **SSO** | Регистрация и авторизация. JWT-токены с хранением в Redis |
| **Catalog** | Управление каталогом контента. Импорт фильмов и сериалов с TMDB, аниме с AniList |
| **Recommendation** | Персональные рекомендации на основе действий пользователя. ИИ-рекомендации через Hugging Face |
| **Notification** | Отправка email-уведомлений при регистрации через SMTP |
| **API Gateway** | Единая точка входа, маршрутизация к микросервисам, JWT-авторизация |

---

## Архитектура
```
                        +-------------+
                        | API Gateway |  :8080
                        +------+------+
                               | gRPC
          +--------------------+--------------------+
          |                    |                    |
   +------+------+    +--------+-------+   +-------+------+
   |     SSO     |    |    Catalog     |   |Recommendation|
   |  :50051     |    |    :50051      |   |   :50051     |
   +------+------+    +--------+-------+   +------+-------+
          |                    |                   |
          |     user-registered|  content-genre    | user-action
          |         (kafka)    |     (kafka)       |  (kafka)
          |                    |                   |
          +----------+  +------+         +---------+
                     |  |                |
                  +--+--+----------------+--+
                  |          Kafka          |
                  +--+-----------------------+
                     | user-registered
              +------+------+
              |Notification |
              +-------------+
```

Потоки данных через Kafka:

| Топик | Отправитель | Получатель | Описание |
|-------|-------------|------------|----------|
| `user-registered` | SSO | Notification | Новый пользователь - приветственное письмо |
| `content-genre` | Catalog (Importer) | Recommendation | Жанры контента - для построения рекомендаций |
| `user-action` | API Gateway | Recommendation | Действия пользователя (лайк, дизлайк) - обновление предпочтений |

Синхронное взаимодействие через gRPC:

| Клиент | Сервер | Описание |
|--------|--------|----------|
| API Gateway | SSO | Регистрация, логин, валидация токена |
| API Gateway | Catalog | Получение контента |
| API Gateway | Recommendation | Получение рекомендаций |

---

## Стек технологий

- **Go 1.25** - основной язык разработки
- **gRPC + Protobuf** - межсервисное взаимодействие
- **PostgreSQL 16** - основное хранилище данных
- **Redis 7** - хранение и обновление JWT-токенов
- **Apache Kafka** - асинхронные события между сервисами
- **Hugging Face (LLaMA 3)** - ИИ-рекомендации
- **Docker & Docker Compose** - контейнеризация
- **Goose** - управление миграциями БД
- **testcontainers-go** - интеграционное тестирование

---

## Запуск приложения

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

---

## Первый запуск

После старта приложения база каталога пуста. Необходимо запустить импорт контента с TMDB:
```bash
docker-compose --profile import run --rm catalog-importer
```

или через Make:
```bash
make import
```

Импортер загрузит **1000 популярных фильмов** и отправит жанры в сервис рекомендаций через Kafka. Процесс занимает ~3-5 минут.

Без импорта эндпоинты каталога и рекомендаций будут возвращать пустые результаты.

---

## API Endpoints

Base URL: `http://localhost:8080`

### Публичные (без токена)

#### Регистрация
```
POST /auth/register
```
```json
{
  "email": "user@example.com",
  "password": "secret123",
  "name": "Ivan",
  "surname": "Ivanov",
  "role": "user"
}
```

Ответ:
```json
{
  "user_id": "uuid"
}
```

#### Авторизация
```
POST /auth/login
```
```json
{
  "email": "user@example.com",
  "password": "secret123"
}
```

Ответ:
```json
{
  "access_token": "eyJhbGci..."
}
```

---

### Защищённые (требуют `Authorization: Bearer <token>`)

#### Каталог

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/api/catalog/get_content` | Получить контент по ID |
| GET | `/api/catalog/get-content-by-external` | Найти контент по внешнему ID |

**GET /api/catalog/get_content**
```json
{
  "content_id": "uuid"
}
```

**GET /api/catalog/get-content-by-external**
```json
{
  "content_id": "uuid",
  "external_source": "tmdb"
}
```

---

#### Действия с контентом

Все эндпоинты принимают `content_id` из URL. `user_id` берётся из JWT токена автоматически.

| Метод | Endpoint | Описание |
|-------|----------|----------|
| POST | `/api/content/{content_id}/like` | Лайк |
| POST | `/api/content/{content_id}/dislike` | Дизлайк |
| POST | `/api/content/{content_id}/favorite` | Добавить в избранное |
| POST | `/api/content/{content_id}/view` | Отметить как просмотренное |

Пример:
```
POST /api/content/d3820b92-9d1e-4893-9495-03694a8bfd2e/like
```

---

#### Рекомендации

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/api/recommendations` | Персональные рекомендации на основе истории |
| POST | `/api/recommendations/explicit` | ИИ-рекомендации по текстовому запросу |
| POST | `/api/recommendations/genres` | Рекомендации по жанрам |
| POST | `/api/recommendations/similar` | Похожий контент |
| GET | `/api/recommendations/trending` | Трендовый контент (просмотры за 24ч) |
| GET | `/api/recommendations/popular` | Популярный контент |

**GET /api/recommendations**
```json
{
  "limit": 10
}
```

**POST /api/recommendations/explicit** - ИИ по текстовому запросу
```json
{
  "query": "хочу что-то про космос и путешествия во времени",
  "limit": 5
}
```

**POST /api/recommendations/genres**
```json
{
  "genres": ["Action", "Thriller"],
  "limit": 10
}
```

**POST /api/recommendations/similar**
```json
{
  "content_id": "d3820b92-9d1e-4893-9495-03694a8bfd2e",
  "limit": 10
}
```

**GET /api/recommendations/trending**
```json
{
  "limit": 10
}
```

**GET /api/recommendations/popular**
```json
{
  "limit": 10
}
```

---

#### Health Check
```
GET /health -> 200 OK
```

---

## Тестирование

В проекте реализованы юнит и интеграционные тесты.  
Интеграционные тесты используют **testcontainers-go** - автоматически поднимают изолированную PostgreSQL в Docker.
```bash
# юнит тесты (быстро, без Docker)
go test ./internal/... -v

# интеграционные тесты catalog
go test ./catalog-system-microservice/tests/... -v -timeout 120s

# интеграционные тесты recommendation
go test ./recommedation-system-microservice/tests/... -v -timeout 120s
```

---

## Структура проекта
```
recsystem/
├── api-gateway/                       # HTTP Gateway
├── sso-microservice/                  # Авторизация
├── catalog-system-microservice/       # Каталог контента
├── recommedation-system-microservice/ # Рекомендации
├── notification-microservice/         # Email уведомления
├── proto/                             # Protobuf схемы
├── docker-compose.yml
├── Makefile
└── .env
```

---

## Makefile команды

| Команда | Описание |
|---------|----------|
| `make up` | Запустить приложение |
| `make up-build` | Пересобрать и запустить |
| `make down` | Остановить |
| `make down-v` | Остановить и удалить данные |
| `make import` | Запустить импорт контента |
| `make logs s=catalog` | Логи конкретного сервиса |
| `make vendor` | Вендоризация зависимостей |