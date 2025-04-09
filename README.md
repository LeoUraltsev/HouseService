# House Service

## Запуск приложения:

Создать файл `.env` и заполнить его по анологии с `.env.example`

```bash
docker build -t house-service:local .
docker compose up -d
```

После запуска базы данных выполнить миграции:

```bash
make up
```

Task List на фикс:

6. Подписка на отправку уведомлений о изменении статуса
7. Разобраться с ci/cd
8. Дописать сервисы
