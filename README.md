# getman

CLI-утилита и Go-библиотека для работы с HTTP-запросами, коллекциями,
окружениями и path placeholders. Формат коллекций совместим с Postman Collection v2.1.

## Установка

```bash
go install github.com/KonnorFrik/getman/cmd/getman@latest
```

Или собрать локально:

```bash
make build
# либо
go build -o getman ./cmd/getman
```

## Быстрый старт CLI

Все данные хранятся в каталоге из флага `--dir` (по умолчанию `.getman`).
Перед началом работы каталог нужно инициализировать.

```bash
# Создать хранилище
./getman init

# Создать окружение и задать переменные
./getman env init dev
./getman env set dev baseUrl=https://api.example.com token=secret

# Создать коллекцию, привязанную к окружению dev
./getman collection new api -e dev -d "Example API"

# Добавить запрос с переменными и path placeholder
./getman req new api -n get-user -X GET --url '{{baseUrl}}/users/:id'

# Задать значение path placeholder
./getman placeholder set id=42

# Выполнить всю коллекцию
./getman collection run api

# Посмотреть сохранённые сущности
./getman collection list
./getman req list api
./getman req info api get-user
./getman env list
./getman env vars dev
./getman placeholder list

# Импортировать коллекцию Postman
./getman import --postman ./postman_example.json
```

Для остановки выполнения при первой ошибке используйте `--fail-fast`.
Конкретный каталог хранилища задаётся флагом `--dir`, например
`./getman --dir ./my-getman collection list`.

### Работа с path placeholders

Path placeholders — это пары ключ-значение для сегментов URL вида `:id`:

```text
https://api.example.com/users/:id
```

- при импорте коллекции Postman найденные placeholders автоматически добавляются в хранилище;
- значение задаётся командой `getman placeholder set id=42`;
- список — `getman placeholder list`;
- удаление — `getman placeholder unset id`;
- если значение не задано, запрос не выполняется и возвращает ошибку
  `path placeholder not found`.

## Использование как библиотеки

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

	if err := client.LoadEnvironment("production"); err != nil {
		log.Fatal(err)
	}

	req, err := getman.NewRequestBuilder().
		Method("GET").
		URL("{{baseUrl}}/users/:id").
		Header("Accept", "application/json").
		AuthBearer("{{token}}").
		Build()
	if err != nil {
		log.Fatal(err)
	}

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

- Выполнение HTTP-запросов с подстановкой переменных `{{variable}}`
- Окружения и глобальные переменные с приоритетом окружения
- Path placeholders вида `:id`
- Коллекции запросов: последовательный и выборочный запуск
- Импорт и экспорт Postman Collection v2.1
- История выполнения и логи
- Форматированный вывод с цветовой индикацией
- CLI на базе Cobra

## Документация

- [Спецификация](specification.md)
- [Ограничения MVP](mvp_limits.md)
- [pkg.go.dev](https://pkg.go.dev/github.com/KonnorFrik/getman)
