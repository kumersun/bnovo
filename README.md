# Guest Service

Микросервис для работы с гостями, реализуюший API для CRUD операций.

## Запуск

docker compose up --build

## API

### POST http://localhost/guest - создать гостя
Request:
```js
{
    "name": "Иван",
    "surname": "Петров",
    "phone": "+70000000000",
    "email": "ivan.petrov@ya.ru",
    "country": "RU"
}
```

### GET http://localhost/guest - получить гостей

### GET http://localhost/guest/{id} - получить гостя

### PUT http://localhost/guest/{id} - обновить гостя
Request:
```js
{
    "name": "Ivan",
    "surname": "Petrov",
    "phone": "+70000000000",
    "email": "ivan.petrov@ya.ru",
    "country": ""
}
```

### DELETE http://localhost/guest/{id} - удалить гостя
Response:
```js
{
    "data": [
        1
    ],
    "error": ""
}
```