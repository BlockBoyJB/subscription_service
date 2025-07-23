# subscription_service

Сервис для агрегации данных об онлайн-подписках пользователей

### Prerequisites

- Docker, Docker Compose
- or Golang 1.24 + postgresql
- make (for running Makefile) (optional)

### Getting started

* Добавить репозиторий к себе
* Создать .env файл в директории с проектом и заполнить информацией из .env.example

### Usage

Запустить сервис можно с помощью `make compose-up` (или `docker-compose up -d --build`)
или `make run` (при наличии go1.24 и локально развернутого postgresql)  
Тесты доступны по команде `make tests`

Документация доступна по адресу `http://localhost:8000/swagger/index.html`

### Примеры запросов

#### Создание

`request`

```shell
curl -X 'POST' \
  'http://localhost:8000/api/v1/subscription' \
  -H 'Content-Type: application/json' \
  -d '{ \
	"service_name": "Yandex", \
	"user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", \
	"price": 600, \
	"start_date": "07-2025" \
}'
```

`response`  
`200`

#### Поиск всех

`request`

```shell
curl -X 'GET' \
  'http://localhost:8000/api/v1/subscription/all'
```

`response`

```json
[
  {
    "id": 1,
    "service_name": "Yandex",
    "price": 600,
    "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a",
    "start_date": "07-2025",
    "end_date": null
  }
]
```

#### Поиск по id

`request`

```shell
curl -X 'GET' \
  'http://localhost:8000/api/v1/subscription/1'
```

`response`

```json
{
  "id": 1,
  "service_name": "Yandex",
  "price": 600,
  "user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a",
  "start_date": "07-2025",
  "end_date": null
}
```

#### Обновление

`request`

```shell
curl -X 'PUT' \
  'http://localhost:8000/api/v1/subscription/1' \
  -H 'Content-Type: application/json' \
  -d '{ \
	"service_name": "Yandex", \
	"user_id": "6114696a-d069-4fad-a3ed-f27c13651c3a", \
	"price": 400, \
	"start_date": "08-2025", \
	"end_date": "09-2025" \
}'
```

`response`  
`200`

#### Удаление

`request`

```shell
curl -X 'DELETE' \
  'http://localhost:8000/api/v1/subscription/1'
```

`response`  
`200`

#### Подсчет суммарной стоимости

Параметры service_name и user_id - опциональные, start и end - обязательные

`request`

```shell
curl -X 'GET' \
  'http://localhost:8000/api/v1/subscription/price?service_name=Yandex&user_id=6114696a-d069-4fad-a3ed-f27c13651c3a&start=01-2025&end=12-2025'
```

`response`

```json
{
  "price": 400
}
```