# Сервис аналитики запросов

Сервис разработан как тестовое задание для стажировки в RWB (Wildberries & Russ)

## Реализованные дополнительные пункты

- **Динамический стоп-лист:** Реализация API для управления стоп-листом, добавление и удаление нежелательных слов «на лету» без перезапуска сервиса.

- **Мониторинг:** Добавление базовых метрик в формате Prometheus для
наблюдения за состоянием сервиса.

- **DX:** Настроенный процесс быстрого локального запуска для проверяющих docker-compose, поднимающий сервис, брокер сообщений(RabbitMQ), базу данныех(Redis), метрики(Prometheus).

## Инструкция по локальному запуску проекта и примеры запросов к API

1) Склонируйте репозиторий
```bash
git clone https://github.com/IvanDrf/analyse-search-requests

cd analyse-search-requests
```

2) Создайте .env
```bash
touch .env
cat example.env > .env
```

3) Создайте конфиг в config/config.yaml
```bash
touch config/config.yaml
```

4) Напишите данные в конфиг (config/config.yaml)
```yaml
app:
  host: 0.0.0.0 # 0.0.0.0 если в docker-compose, localhost если локально
  port: 8080 #должен совпадать с APP_PORT в .env
  request_time: 10s
  search_duration: 5m
  search_interval: 5  

  logger_level: "DEBUG"
  add_source: false

metrics:
  host: 0.0.0.0 # 0.0.0.0 если в docker-compose, localhost если локально
  port: 2112

database:
  host: database # database если в docker-compose, localhost если локально
  port: 6379 # должен совпадать с REDIS_PORT в .env
  user: ""
  password: ""
  database: 1

  duplicate_time: 7m
  bad_word_key: "bad_word_key"

broker:
  host: broker # broker если в docker-compose, localhost если локально
  port: 5672 # должен совпадать с RABBITMQ_PORT в .env
  user: user # пользователь должен совпадать с RABBITMQ_USER в .env
  password: "123" # пароль должен совпадать с RABBITMQ_PASSWORD в .env
  queue: searches_queue
  workers: 10
```

5) Запустите контейнеры
```bash
docker-compose up -d 
```

6) Можно проверить, что контейнеры запущены
```bash
docker ps
```

7) Можно использовать сервис

## Контракт данных
- **Сообщение из брокера**
    ```golang
    type Message struct {
	    ID            uuid.UUID `json:"message_id"`
	    SearchMessage string    `json:"search_message"`
	    Date          time.Time `json:"search_time"`
    }
    ```

- Почему я выбрал такой контракт
1) ID - UUID необходимо, чтобы идентифицировать сообщение и проверять, что бы не записывать дважды одинаковое сообщение. Если отправитель случайно отправил одно сообщение два раза в очередь, сервис засчитает только одно.
<br></br>
2) SearchMessage - string сам поисковой запрос
<br></br>
3) Date - время поиска сообщения, важно, чтобы правильно оценивать дату сообщения. Например, сервис упал, а сообщения в очереди остались, поднимаем сервис и видим там сообщения с датой на час раньше, тогда мы не будем учитывать его при поиске топа самых популярны запросов, если бы времени не было, а использовали бы ```time.Now().UTC```, то учитывали бы старые сообщения.

- **Запросы к API**

1) Добавление стоп-слова
```bash
POST /api/v1/bad
```
Тело запроса
```json
{
    "bad_word": "bad"
}
```
Ответ ```http.StatusNoContent```

2) Удаление стоп-слова
```bash
DELETE /api/v1/bad
```
Тело запроса
```json
{
    "bad_word": "bad"
}
```
Ответ: ```http.StatusNoContent```

3) Получение топ N самых популярных запросов
```bash
GET /api/v1/searches?limit=N
```
Запрос: ```GET /api/v1/searches?limit=3```

Тело ответа
```json
[
    {
        "search_message": "python", // запрос
        "amount": 8 // количество таких запросов
    },
    {
        "search_message": "golang",
        "amount": 4
    },
    {
        "search_message": "c++",
        "amount": 2
    }
]
```

## Архитектура
- База данных
Я выбрал **Redis**, потому что по ТЗ количество оппераций чтения в 10-50 раз больше, чем сообщений в брокере, то есть нужно будет очень часто обращаться к бд. Redis для этого идеально подходит, так как хранит все данные в оперативной памяти, что гораздо быстрее, чем дисковая память PostgreSQL.

- Брокер сообщений
Я выбрал **RabbitMQ**, потому что работал с ним и он хорошо подходит для данной задачи. 

- Хранение топа
Я храню топ сообщений в **Sorted Set(zset)** в Redis, потому что эта структура данныз сразу сортирует входные данные по score, в этом случае score это количество таких запросовю

- Хранение стоп-листа
Я храню стоп-лис в Set, эта структура данных отлично подходит, потому что хранит уникальные значения и поиск значений в ней O(1)

- Бизнес-логика
    - Изначально я хотел хранить данные в PostgreSQL. Создать таблицу ```searches```, в ней 4 поля ```search, amount, date, status```. Поле ```status``` отвечает за то, внесено ли слово стоп лист или нет. Создать таблицу ```bad_words``` и в ней хранить стоп слова и при вставке/удалении нового слова с помощью триггера обновлять статус поисковых запроов, а при ```SELECT``` запросе топ N популярных запросов с помощью триггера удалять сначала старые запросы и только потом искать новые. Но я понял, что это будет слишком долго и накладно для PostgreSQL в условиях highload, поэтому начал искать другие базы данных, Redis идеально подошла.

    - Еще одна задача, с которой я сталкнулся, это проверка сообщений-дубликатов, если в брокер случайно отправили одно сообщение два раза, то нужно учитывать сообщение первый раз, а второе сообщение игнорировать, для этого я использовал ключ в Redis ```duplicate:uuid```, который показывал, учли уже это сообщение или нет, если ключ есть, значит уже посчитали.

### Структура проекта
Подробное описание структуры проекта - https://habr.com/ru/articles/911018/

```bash
├── cmd
│   └── main.go
├── config
│   ├── config.example.yaml
│   └── config.yaml
├── docker-compose.yaml
├── Dockerfile
├── example.env
├── go.mod
├── go.sum
├── internal
│   ├── app
│   │   ├── app.go
│   │   └── fabric.go
│   ├── config
│   │   ├── app.go
│   │   ├── broker.go
│   │   ├── config.go
│   │   ├── database.go
│   │   └── metrics.go
│   ├── domain
│   │   ├── models
│   │   │   ├── bad_word.go
│   │   │   ├── error.go
│   │   │   └── message.go
│   │   ├── ports
│   │   │   ├── messaging
│   │   │   │   └── consumer.go
│   │   │   ├── repo
│   │   │   │   └── message.go
│   │   │   └── service
│   │   │       └── message.go
│   │   └── rules
│   │       ├── bad_word.go
│   │       ├── limit.go
│   │       └── message.go
│   ├── infrastructure
│   │   ├── adapters
│   │   │   ├── logger.go
│   │   │   └── metrics.go
│   │   ├── messaging
│   │   │   └── rabbitmq
│   │   │       ├── connect.go
│   │   │       └── consumer.go
│   │   ├── persistence
│   │   │   └── redis
│   │   │       ├── connect.go
│   │   │       └── message_repo.go
│   │   └── service
│   │       └── message.go
│   └── interfaces
│       └── http
│           ├── handlers.go
│           ├── metrics_server.go
│           ├── middleware
│           │   └── metrics.go
│           ├── search_server.go
│           └── utils.go
├── LICENSE
├── prometheus.yaml
└── README.md
```
