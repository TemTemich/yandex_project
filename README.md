# Project

Для запуска:
`make start_with_logs`

GUI доступен в папке frontend

## Примеры работы API:

#### Добавление выражения для вычисления:
---
Запрос

```bash
curl --location 'http://localhost:8080/expressions' \
--header 'Content-Type: application/json' \
--data '{
    "expression":"2+2*2/2+2-3-10*27123+2+10000"
}'
```
Пример ответа:

```json
{
    "id": "685da2a4-9ff8-4815-b954-297b8dac0c4e"
}

```

#### Узнать статус и результат по ID:
---
Запрос
```bash
curl --location 'http://localhost:8080/expressions/685da2a4-9ff8-4815-b954-297b8dac0c4e'
```

Пример ответа:
```json
{
    "id": "685da2a4-9ff8-4815-b954-297b8dac0c4e",
    "status": "done",
    "result": "-261225"
}

```

#### Получить данные по всем выражениям:
---
Запрос
```bash
curl --location 'http://localhost:8080/expressions'
```

Пример ответа:
```json
[
    {
        "id": "29ebb5c5-ddba-4f82-9a56-8b79ebdb41ba",
        "expression": "1+2",
        "status": "work",
        "result": ""
    },
    ...
    {
        "id": "a379abcd-c7f1-4956-8307-25510998ea24",
        "expression": "1+2",
        "status": "done",
        "result": "3"
    }
]

```



#### Посмотреть результат работы всех операций:
---
Запрос
```bash
curl --location 'http://localhost:8080/operation/all'
```

Пример ответа:
```json
[
    {
        "id": "c989284d-a51b-401a-b74a-05c3bd61b500",
        "operation": "+2",
        "duration": "5.264000",
        "result": "2"
    },
    {
        "id": "5ba947d3-d35c-44af-99a6-76f7f84b7bc9",
        "operation": "1",
        "duration": "8.908000",
        "result": "1"
    },
    {
        "id": "6f23131e-0c59-40a4-9c0f-985ebbd60dc2",
        "operation": "1",
        "duration": "3.114000",
        "result": "1"
    },
    ...
    {
        "id": "8b9fa032-e299-4088-9391-35559a5c0889",
        "operation": "+2*2/2",
        "duration": "22.331000",
        "result": "2"
    }
]

```
/////////