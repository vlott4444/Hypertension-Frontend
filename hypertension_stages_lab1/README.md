# ЛР №1 — стадии гипертонической болезни

Проект переделан из исходного Go/Gin-приложения под тему **стадий гипертонической болезни**. В этой учебной модели стадия определяется по **систолическому АД (SBP/САД)**:

- 1 стадия: 140–159 мм рт. ст.
- 2 стадия: 160–179 мм рт. ст.
- 3 стадия: >180 мм рт. ст. (как в макете)

Интерфейс повторяет приложенный Figma/PDF-макет: мобильный портретный экран, красные карточки, светло-розовая нижняя навигация, лента с видео и плитка в 2 столбца.

## Запуск

1. Запустить Docker Desktop.
2. В корне проекта выполнить:

```bash
docker compose up -d
```

`minio-init` автоматически создаст публичный bucket `hypertension-media` и загрузит туда изображения и MP4 из `resources/minio`.

Minio Console: `http://localhost:9001`

- login: `root`
- password: `rootpassword`

3. Запустить Go-приложение:

```bash
go mod tidy
go run ./cmd/hypertension_stages
```

4. Открыть:

- `http://localhost:8080/feed/` — лента;
- `http://localhost:8080/draft` — черновик/добавление;
- `http://localhost:8080/stages` — плитка.

## Ровно 3 GET-маршрута

```go
router.GET("/feed/*serviceID", hypertensionHandler.ShowHypertensionFeed)
router.GET("/draft", hypertensionHandler.ShowHypertensionDraft)
router.GET("/stages", hypertensionHandler.ShowHypertensionGrid)
```

### Лента по ID и next=true

- `/feed/2` — карточка с ID 2;
- `/feed/2?next=true` — следующая опубликованная карточка после ID 2;
- `/feed/` — переход из нижней панели без ID, открывает первую опубликованную карточку.

### Серверная фильтрация

Пример:

- `/stages?sbp=150` → 1 стадия;
- `/stages?sbp=170` → 2 стадия;
- `/stages?sbp=190` → 3 стадия.

Значение `sbp` сохраняется в input после GET-запроса.

## Статусы

В единственной коллекции `hypertensionServices` есть услуги со статусами:

- `черновик` — ровно одна карточка, она показывается на `/draft`;
- `опубликован` — показываются в ленте и плитке;
- `удален` — в интерфейс не попадает.

## Minio

В каждой модели есть **два отдельных поля**:

```go
ImageURL string
VideoURL string
```

Оба URL ведут в Minio, например:

```text
http://localhost:9000/hypertension-media/stage1.png
http://localhost:9000/hypertension-media/stage1.mp4
```

Эти URL используются во всех трех шаблонах. В плитке `VideoURL` также выводится в HTML как `data-video-url`, поэтому его видно во вкладке Network → Response без JavaScript.

## Где считается количество лайков

Количество лайков не хранится отдельным числом. В модели лежит вложенная коллекция `Likes` с `UserID`, а контроллер `BuildHypertensionCardView` вычисляет `LikesCount` через `len(...)`.

## Цвета и элементы дизайна из макета

В `resources/styles/hypertension.css` явно указаны три ключевых цвета:

- `#B6001A` — красный;
- `#FBF2F3` — светло-розовый;
- `#333333` — графитовый текст.

Также реализованы форма скругленных карточек, pill-навигация снизу и `:hover` для карточек/кнопок.
