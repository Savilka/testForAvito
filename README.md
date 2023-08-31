# Тестовое задание

Полное описание API находится в `api.json`

### Запуск

Для запуска нужно создать .env файл по примеру .env.example . Также нужно в папке config создать файл cfg.yaml по
примеру cfg.example.yaml.
Дальше нужно вызвать следующие команды:

docker compose -f .\docker-compose.yml build

docker compose -f .\docker-compose.yml up

Api будет доступно по адресу localhost:8000

### Краткое описание

GET v1/user/add \
Добавление нового пользователя. Парметры не нужны. В ответ приходит id

GET v1/user/{id}/getSegments  \
Метод для получения сегментов пользователя.

Пример: localhost:8000/v1/user/1/getSegments

PUT v1/user/{id}/addSegments  \
Метод для добавления пользователя в группы.

Тело запроса:

```json
{
  "new_segments": [
    {
      "name": "aaa",
      "ttl": 100
    }
  ],
  "old_segments": [
    {
      "name": "ccc"
    }
  ]
}
```

Пример: localhost:8000/v1/user/1/addSegments

POST v1/segment/add  \
Метод для создания нового сегмента. Принимает опциональный параметр percent

Тело запроса:

```json
{
  "slug": "ccc",
  "percent": 50
}
```

Пример: localhost:8000/v1/segment/add

DELETE v1/segment/add  \
Метод для удаления сегмента

Тело запроса:

```json
{
  "slug": "ccc"
}
```

Пример: localhost:8000/v1/segment/delete

