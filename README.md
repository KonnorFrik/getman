# getman

CLI программа для работы с HTTP запросами, коллекциями и окружениями в Go.

## Установка

```bash
go get github.com/KonnorFrik/getman@cli_cobra_latest
```

# TODO: rewrite readme for cli 

## Быстрый старт

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/KonnorFrik/getman"
)

func main() {
    client, err := getman.NewClientWithDefaults()
    if err != nil {
        log.Fatal(err)
    }

    err = client.LoadEnvironment("production")
    if err != nil {
        log.Fatal(err)
    }

    req := getman.NewRequestBuilder().
        Method("GET").
        URL("{{baseUrl}}/users").
        Header("Accept", "application/json").
        AuthBearer("{{token}}").
        Build()

    result, err := client.ExecuteRequest(req)
    if err != nil {
        log.Fatal(err)
    }

    if result.Error != "" {
        fmt.Printf("Error: %s\n", result.Error)
        return
    }

    getman.PrintResponse(result.Response)
    fmt.Printf("Duration: %v\n", result.Duration)
}
```

## Основные возможности

- Выполнение HTTP запросов с поддержкой переменных
- Управление окружениями и переменными
- Работа с коллекциями запросов
- Импорт из Postman Collection v2.1
- История выполнения запросов
- Форматирование и визуализация результатов

## Документация

Полная документация доступна на [pkg.go.dev](https://pkg.go.dev/github.com/KonnorFrik/getman).

