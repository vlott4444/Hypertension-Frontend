# Hypertension Stages — REST API

Веб-сервис для работы с услугами (стадиями гипертонии) на Go + Gin + GORM + PostgreSQL + Minio.

## Стек

- **Go** — язык
- **Gin** — HTTP-роутинг
- **GORM** — ORM для PostgreSQL
- **PostgreSQL** — БД
- **Minio** — S3-хранилище для изображений и видео
- **Insomnia / Postman** — тестирование API

## Запуск

```bash
go run cmd/migrate/main.go   # миграции
go run cmd/seed/main.go      # заполнение БД
go run cmd/server/main.go    # запуск сервера на :8080
