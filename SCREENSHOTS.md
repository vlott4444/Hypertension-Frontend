# Порядок показа ЛР6 (16 скриншотов)

1. `docker compose stop backend`; `http://localhost:3000/feed`: гость и mock.
2. `http://localhost:3000/stages`: mock и фильтры.
3. `docker compose start backend`; `/feed` — реальная карточка.
4. `/stages` — опубликованный список Go.
5. `/stages/4` — карточка по ID Go.
6. `/login` войти под `test_user / test_pass`, `/add` — создать и опубликовать карточку.
7. В Adminer/PostgreSQL изменить `title` опубликованной карточки. Обновить `/stages`.
8. Поменять `systolic_bp` той же карточки, обновить диапазон фильтра во фронтенде.
9. После входа Redux DevTools показывает true; после F5 false, но Application → Local Storage содержит `hypertension.jwt`.
10. В Insomnia GET `http://localhost:8080/api/hypertension_stages` с `Authorization: Bearer <JWT>`; is_owner=true у собственных карточек.
11. Код `frontend/src/store/store.ts`, `authSlice.ts`, `filterSlice.ts`.
12. Код `frontend/src/api/http.ts`, `frontend/src/api/generated/HypertensionStagesApi.ts`, `stagesSlice.ts`.
13. Код `frontend/openapi/lab4.contract.json`, `frontend/scripts/generate-api.mjs`.
14. Код `frontend/src/store/stagesSlice.ts` и `feedSlice.ts` с вызовами тематических методов сгенерированного клиента.
15. Network → Fetch/XHR: URL `/api/hypertension_stages?min_sbp=...` на порту 3000; показать `frontend/vite.config.ts`.
16. Network → Img/Media: URL `/s3/hypertension-media/stage1.png` на порту 3000; MinIO Console c bucket/object; `frontend/src/components/StageMedia.tsx`.

Скриншоты нужно сделать самим во время работающей демонстрации. Они не включены в архив.
