# Лабораторная № 6 — Hypertension Stages

**Полная исходная сборка:** React + Redux Toolkit/Thunk + Axios + Vite (frontend); Go/Gin + GORM (backend); PostgreSQL, Redis (отзыв JWT) и MinIO (изображения/видео).

## Запуск одним способом (Docker Desktop)

Из корня распакованного архива:

```bash
docker compose up --build -d
docker compose ps
```

Дождитесь, пока контейнеры станут running. Первый запуск скачивает образы и Go/npm-зависимости и может занимать продолжительное время.

* Фронтенд: **http://localhost:3000**
* Go API: **http://localhost:8080/api/hypertension_stages**
* Администрирование БД: **http://localhost:8081** (System: PostgreSQL, Server: `postgres`, User: `postgres`, Password: `postgres1234`, Database: `test_db`)
* MinIO: **http://localhost:9001** (root / rootpassword; bucket `hypertension-media`)
* Redis: localhost:6379

Для просмотра ошибок:

```bash
docker compose logs -f backend
docker compose logs -f frontend
```

Для остановки без удаления базы:

```bash
docker compose down
```

`docker compose down -v` удалит и тестовые данные PostgreSQL, Redis и MinIO.

## Тестовые пользователи

После первого запуска сервис автоматически выполняет миграцию и seed. Используйте `test_user` / `test_pass` (владеет всеми исходными опубликованными карточками и черновиком). Другие тестовые пользователи: `elizabeth` / `pass1`, `anna` / `pass2`. Можно регистрировать новых через `/login` (новый пароль от 6 символов). Пароли при новой установке хранятся в bcrypt-хешах; первоначальные пароли предназначены исключительно для демонстрации.

Войти -> плитка -> лайк/удаление принадлежащих карточек -> /add -> создать черновик -> опубликовать. Для скриншота с F5: Redux возвращается в состояние гостя, JWT сохраняется в Local Storage для использования в Insomnia. Это учебное поведение, а не рекомендация по безопасности.

## Mock без бэкенда

Если нужно показать первые два скриншота (мок без сервиса), остановите контейнер `backend`:

```bash
docker compose stop backend
```

После обновления страницы **без входа** плитка и лента отображают mock-данные. Далее:

```bash
docker compose start backend
```

## Как проверить API

```bash
curl -i 'http://localhost:8080/api/hypertension_stages?min_sbp=140&max_sbp=180&title=стадия'
# POST /api/auth/login:
curl -s -X POST http://localhost:8080/api/auth/login -H 'Content-Type: application/json' -d '{"username":"test_user","password":"test_pass"}'
```

Ответ login содержит `token`. Для защищённых операций передавайте `Authorization: Bearer <token>`. В ответе списка с токеном у собственных карточек будет `is_owner: true`.

## API и медиа

* `GET /api/hypertension_stages?min_sbp=&max_sbp=&title=` — серверный фильтр, список со статусом и is_owner/is_liked/likes_count
* `GET /api/feed` и `GET /api/feed/:id?next=true` — одна карточка или следующая
* `GET /api/draft` — свой черновик (JWT)
* `POST /api/hypertension_stages` — multipart FormData, поля title/description/systolic_bp/diastolic_bp, опциональные image/video (JWT)
* `PUT /api/hypertension_stages/:id/publish`, `DELETE /api/hypertension_stages/:id`, `POST /api/hypertension_stages/:id/like` — JWT и проверка прав
* `POST /api/auth/register`, `/api/auth/login`, `/api/auth/logout` — авторизация, JWT, отзыв через Redis

Прокси Vite: `/api/*` на Go, `/s3/*` на MinIO (публичное чтение объектов в demo bucket, запись только с S3-подписью). У файлов новых карточек имена случайные. Исходные картинки и видео загружаются при инициализации MinIO из `frontend/public/media`.

Axios-клиент создан из `frontend/openapi/lab4.contract.json`, его исходник и генератор лежат в `frontend/src/api/generated` и `frontend/scripts/generate-api.mjs`. Для регенерации: `cd frontend && npm run generate-api`. Файл OpenAPI отражает интерфейс объединённого приложения; экспорт Swagger UI из Go не добавлялся.

## Проверка ЛР6

Последовательность для скриншотов см. `SCREENSHOTS.md`. Диаграммы не включены по запросу.

## Отдельный запуск без Docker

Нужны Go 1.27, Node 22+, PostgreSQL, Redis и MinIO. Создайте `backend/.env` с `DB_HOST=localhost`, `DB_PORT=5455`, `DB_USER=postgres`, `DB_PASS=postgres1234`, `DB_NAME=test_db`, `REDIS_ADDR=localhost:6379`, `JWT_SECRET=<не менее 32 символов>`, `MINIO_ENDPOINT=http://localhost:9000`, `MINIO_ROOT_USER=root`, `MINIO_ROOT_PASSWORD=rootpassword`.

```bash
cd backend
go run ./cmd/migrate
go run ./cmd/seed
go run ./cmd/server
# в отдельном терминале из корня проекта:
cd frontend
npm install
npm run dev
```

*Для этого варианта контейнеры PostgreSQL, Redis, MinIO и бакет MinIO всё равно должны быть запущены отдельно. Проще использовать основной docker compose.*

## Проверка сборки

В этом окружении исходники проверены статически, но `docker compose up --build` не запускался: Docker не установлен. Поэтому работу всей связки с контейнерами необходимо подтвердить на машине с Docker Desktop. Конфигурация рассчитана на локальную учебную демонстрацию, не на публикацию в Интернет.
