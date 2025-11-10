# Детализированная спецификация библиотеки getman

## 1. Архитектура библиотеки

### 1.1. Компоненты

Библиотека состоит из следующих основных компонентов:

1. **Core** - построение и выполнение HTTP запросов
   - RequestBuilder - построитель запросов с fluent API
   - HTTPClient - клиент для выполнения запросов
   - VariableResolver - разрешение переменных в запросах

2. **Variables** - система переменных и окружений
   - Environment - управление окружениями
   - VariableStore - хранилище переменных с приоритетами

3. **Collections** - управление коллекциями запросов
   - Collection - структура коллекции
   - CollectionExecutor - выполнение коллекций

4. **Storage** - управление файлами и директориями
   - FileStorage - работа с файлами коллекций и окружений
   - HistoryStorage - сохранение истории выполнения
   - LogStorage - сохранение логов

5. **Importer** - импорт из Postman Collection v2.1
   - PostmanImporter - парсинг и конвертация формата Postman

6. **Formatter** - форматирование и визуализация
   - ResponseFormatter - форматирование ответов
   - ResultFormatter - форматирование результатов выполнения

### 1.2. Структура данных

#### Request (HTTP запрос)
```go
type Request struct {
    Method  string            `json:"method"`
    URL     string            `json:"url"`
    Headers map[string]string `json:"headers,omitempty"`
    Body    *RequestBody      `json:"body,omitempty"`
    Auth    *Auth             `json:"auth,omitempty"`
    Timeout *Timeout          `json:"timeout,omitempty"`
    Cookies *CookieSettings   `json:"cookies,omitempty"`
}

type RequestBody struct {
    Type        string `json:"type"`
    Content     []byte `json:"content"`
    ContentType string `json:"content_type,omitempty"`
}

type Auth struct {
    Type     string `json:"type"`
    Username string `json:"username,omitempty"`
    Password string `json:"password,omitempty"`
    Token    string `json:"token,omitempty"`
    APIKey   string `json:"api_key,omitempty"`
    KeyName  string `json:"key_name,omitempty"`
    Location string `json:"location,omitempty"`
}

type Timeout struct {
    Connect time.Duration `json:"connect"`
    Read    time.Duration `json:"read"`
}

type CookieSettings struct {
    AutoManage bool `json:"auto_manage"`
}
```

#### Response (HTTP ответ)
```go
type Response struct {
    StatusCode int               `json:"status_code"`
    Status     string            `json:"status"`
    Headers    map[string][]string `json:"headers"`
    Body       []byte            `json:"body"`
    Duration   time.Duration     `json:"duration"`
    Size       int64             `json:"size"`
}
```

#### Environment (Окружение)
```go
type Environment struct {
    Name      string            `json:"name"`
    Variables map[string]string `json:"variables"`
}
```

#### Collection (Коллекция)
```go
type Collection struct {
    Name        string         `json:"name"`
    Description string         `json:"description,omitempty"`
    Items       []*RequestItem `json:"items"`
}

type RequestItem struct {
    Name    string   `json:"name"`
    Request *Request `json:"request"`
}
```

#### RequestExecution (Результат выполнения запроса)
```go
type RequestExecution struct {
    Request   *Request      `json:"request"`
    Response  *Response     `json:"response,omitempty"`
    Error     string        `json:"error,omitempty"`
    Duration  time.Duration `json:"duration"`
    Timestamp time.Time     `json:"timestamp"`
}
```

#### ExecutionResult (Результат выполнения коллекции)
```go
type ExecutionResult struct {
    CollectionName string              `json:"collection_name"`
    Environment    string              `json:"environment"`
    StartTime      time.Time           `json:"start_time"`
    EndTime        time.Time           `json:"end_time"`
    TotalDuration  time.Duration       `json:"total_duration"`
    Requests       []*RequestExecution `json:"requests"`
    Statistics     *Statistics         `json:"statistics"`
}

type Statistics struct {
    Total    int           `json:"total"`
    Success  int           `json:"success"`
    Failed   int           `json:"failed"`
    AvgTime  time.Duration `json:"avg_time"`
    MinTime  time.Duration `json:"min_time"`
    MaxTime  time.Duration `json:"max_time"`
}
```

## 2. Система переменных

### 2.1. Области видимости

1. **Environment (окружение)** - приоритет выше
   - Загружается из JSON файла в `~/.getman/environments/{name}.json`
   - Активное окружение устанавливается через `LoadEnvironment()`

2. **Global (глобальная)** - стандартное окружение
   - Используется, если не загружено другое окружение
   - Устанавливается через `SetGlobalVariable()`

### 2.2. Правила подстановки

- **Синтаксис**: `{{variable}}`
- **Подстановка**: Происходит во всех частях запроса:
  - URL (включая path и query параметры)
  - Headers (ключи и значения)
  - Body (содержимое тела запроса)
- **Валидация**: Проверка наличия всех переменных перед выполнением запроса
- **Приоритет**: При конфликте имен переменная из активного окружения перезаписывает глобальную
- **Типы**: Все переменные - строки

### 2.3. Примеры использования переменных

```
URL: {{baseUrl}}/users/{{userId}}?token={{apiToken}}
Header: Authorization: Bearer {{token}}
Body: {"name": "{{userName}}", "email": "{{userEmail}}"}
```

## 3. Структура директорий и файлов

### 3.1. Постоянное хранилище (~/.getman/)

```
~/.getman/
├── collections/          # JSON файлы коллекций
│   └── {name}.json
├── environments/         # JSON файлы окружений
│   └── {name}.json
├── history/              # История выполнения запросов
│   └── {timestamp}.json  # Файлы истории с временными метками
├── logs/                 # Логи выполнения
│   └── {timestamp}.json  # Файлы логов с временными метками
└── config.yaml          # Конфигурация библиотеки
```

**Формат timestamp**: `DD_MM_YY_HH_MM_SS` (например: `01_12_25_22_55_39`)

### 3.2. Форматы файлов

#### Environment JSON
```json
{
  "name": "production",
  "variables": {
    "baseUrl": "https://api.example.com",
    "token": "abc123xyz",
    "apiKey": "key_12345"
  }
}
```

#### Collection JSON
Базируется на Postman Collection v2.1, упрощенная версия без скриптов и тестов.

Пример структуры:
```json
{
  "name": "My API Collection",
  "description": "Collection description",
  "items": [
    {
      "name": "Get Users",
      "request": {
        "method": "GET",
        "url": "{{baseUrl}}/users",
        "headers": {
          "Accept": "application/json"
        },
        "auth": {
          "type": "bearer",
          "token": "{{token}}"
        }
      }
    }
  ]
}
```

#### Config YAML
```yaml
storage:
  base_path: ~/.getman

defaults:
  timeout:
    connect: 30s
    read: 30s
  cookies:
    auto_manage: true

logging:
  level: info
  format: text
```

## 4. API библиотеки

### 4.1. Основные типы

```go
package getman

import (
    "net/http"
    "time"
)

type Client struct {
    storage    *Storage
    env        *Environment
    globalEnv  *Environment
    config     *Config
    httpClient *http.Client
}

type Storage struct {
    basePath string
}

type Config struct {
    Storage  StorageConfig  `yaml:"storage"`
    Defaults DefaultsConfig `yaml:"defaults"`
    Logging  LoggingConfig  `yaml:"logging"`
}

type StorageConfig struct {
    BasePath string `yaml:"base_path"`
}

type DefaultsConfig struct {
    Timeout TimeoutConfig `yaml:"timeout"`
    Cookies CookiesConfig `yaml:"cookies"`
}

type TimeoutConfig struct {
    Connect time.Duration `yaml:"connect"`
    Read    time.Duration `yaml:"read"`
}

type CookiesConfig struct {
    AutoManage bool `yaml:"auto_manage"`
}

type LoggingConfig struct {
    Level  string `yaml:"level"`
    Format string `yaml:"format"`
}

type LogEntry struct {
    Time    time.Time `json:"time"`
    Level   string    `json:"level"`
    Message string    `json:"message"`
}
```

### 4.2. Инициализация

```go
// NewClient создает новый клиент с базовым путем
func NewClient(basePath string) (*Client, error)

// NewClientWithConfig создает новый клиент с конфигурацией из файла
func NewClientWithConfig(configPath string) (*Client, error)

// NewClientWithDefaults создает клиента с путями по умолчанию (~/.getman)
func NewClientWithDefaults() (*Client, error)
```

### 4.3. Управление окружениями

```go
// LoadEnvironment загружает окружение по имени
func (c *Client) LoadEnvironment(name string) error

// SaveEnvironment сохраняет окружение в файл
func (c *Client) SaveEnvironment(env *Environment) error

// ListEnvironments возвращает список всех доступных окружений
func (c *Client) ListEnvironments() ([]string, error)

// DeleteEnvironment удаляет окружение
func (c *Client) DeleteEnvironment(name string) error

// GetCurrentEnvironment возвращает текущее активное окружение
func (c *Client) GetCurrentEnvironment() *Environment

// SetGlobalVariable устанавливает глобальную переменную
func (c *Client) SetGlobalVariable(key, value string)

// GetGlobalVariable получает значение глобальной переменной
func (c *Client) GetGlobalVariable(key string) (string, bool)

// GetVariable получает значение переменной (с учетом приоритета: окружение > глобальная)
func (c *Client) GetVariable(key string) (string, bool)

// ResolveVariables разрешает все переменные в строке
func (c *Client) ResolveVariables(template string) (string, error)
```

### 4.4. Управление коллекциями

```go
// LoadCollection загружает коллекцию по имени
func (c *Client) LoadCollection(name string) (*Collection, error)

// SaveCollection сохраняет коллекцию в файл
func (c *Client) SaveCollection(collection *Collection) error

// ListCollections возвращает список всех доступных коллекций
func (c *Client) ListCollections() ([]string, error)

// DeleteCollection удаляет коллекцию
func (c *Client) DeleteCollection(name string) error

// ImportFromPostman импортирует коллекцию из файла Postman Collection v2.1
func (c *Client) ImportFromPostman(filePath string) (*Collection, error)

// ExportToPostman экспортирует коллекцию в формат Postman Collection v2.1
func (c *Client) ExportToPostman(collection *Collection, filePath string) error
```

### 4.5. Построение запросов (Fluent API)

```go
// NewRequestBuilder создает новый построитель запросов
func NewRequestBuilder() *RequestBuilder

type RequestBuilder struct {
    method  string
    url     string
    headers map[string]string
    body    *RequestBody
    auth    *Auth
    timeout *Timeout
    cookies *CookieSettings
}

// Method устанавливает HTTP метод
func (b *RequestBuilder) Method(method string) *RequestBuilder

// URL устанавливает URL (с поддержкой переменных {{variable}})
func (b *RequestBuilder) URL(url string) *RequestBuilder

// Header добавляет заголовок (или перезаписывает существующий)
func (b *RequestBuilder) Header(key, value string) *RequestBuilder

// Headers устанавливает несколько заголовков сразу
func (b *RequestBuilder) Headers(headers map[string]string) *RequestBuilder

// BodyJSON устанавливает тело запроса в формате JSON
func (b *RequestBuilder) BodyJSON(data interface{}) *RequestBuilder

// BodyXML устанавливает тело запроса в формате XML
func (b *RequestBuilder) BodyXML(data string) *RequestBuilder

// BodyRaw устанавливает тело запроса как raw данные
func (b *RequestBuilder) BodyRaw(data []byte, contentType string) *RequestBuilder

// BodyBinary устанавливает тело запроса как binary данные
func (b *RequestBuilder) BodyBinary(data []byte, contentType string) *RequestBuilder

// AuthBasic устанавливает Basic Authentication
func (b *RequestBuilder) AuthBasic(username, password string) *RequestBuilder

// AuthBearer устанавливает Bearer Token authentication
func (b *RequestBuilder) AuthBearer(token string) *RequestBuilder

// AuthAPIKey устанавливает API Key authentication
// location: "header" или "query"
func (b *RequestBuilder) AuthAPIKey(keyName, keyValue, location string) *RequestBuilder

// Timeout устанавливает таймауты подключения и чтения
func (b *RequestBuilder) Timeout(connect, read time.Duration) *RequestBuilder

// CookiesAutoManage включает/выключает автоматическое управление cookies
func (b *RequestBuilder) CookiesAutoManage(autoManage bool) *RequestBuilder

// Build создает Request из построителя
func (b *RequestBuilder) Build() (*Request, error)
```

### 4.6. Выполнение запросов

```go
// ExecuteRequest выполняет один HTTP запрос
func (c *Client) ExecuteRequest(req *Request) (*RequestExecution, error)

// ExecuteCollection выполняет все запросы из коллекции последовательно
func (c *Client) ExecuteCollection(collectionName string) (*ExecutionResult, error)

// ExecuteCollectionSelective выполняет выборочные запросы из коллекции
func (c *Client) ExecuteCollectionSelective(collectionName string, itemNames []string) (*ExecutionResult, error)

// ValidateRequest проверяет запрос перед выполнением (переменные, формат)
func (c *Client) ValidateRequest(req *Request) error
```

### 4.7. История и логирование

```go
// GetHistory возвращает историю выполнения запросов из постоянного хранилища
// Читает файлы из ~/.getman/history/ и возвращает последние limit записей
func (c *Client) GetHistory(limit int) ([]*RequestExecution, error)

// GetLastExecution возвращает результат последнего выполнения из постоянного хранилища
// Читает последний файл из ~/.getman/history/
func (c *Client) GetLastExecution() (*ExecutionResult, error)

// GetLogs возвращает логи последнего выполнения из постоянного хранилища
// Читает последний файл из ~/.getman/logs/
func (c *Client) GetLogs() ([]byte, error)

// ClearHistory очищает историю выполнения (удаляет все файлы из ~/.getman/history/)
func (c *Client) ClearHistory() error

// SaveHistory сохраняет историю в файл в постоянном хранилище
// Создает файл с timestamp в имени: ~/.getman/history/{timestamp}.json
func (c *Client) SaveHistory(result *ExecutionResult) error

// SaveLogs сохраняет логи в файл в постоянном хранилище
// Создает файл с timestamp в имени: ~/.getman/logs/{timestamp}.json
func (c *Client) SaveLogs(logs []LogEntry) error
```

### 4.8. Визуализация и форматирование

```go
// FormatResponse форматирует ответ в читаемый текст
func FormatResponse(resp *Response) string

// FormatRequest форматирует запрос в читаемый текст
func FormatRequest(req *Request) string

// FormatExecutionResult форматирует результат выполнения коллекции
func FormatExecutionResult(result *ExecutionResult) string

// FormatStatistics форматирует статистику выполнения
func FormatStatistics(stats *Statistics) string

// PrintResponse выводит ответ в консоль с цветовой индикацией
func PrintResponse(resp *Response)

// PrintRequest выводит запрос в консоль
func PrintRequest(req *Request)

// PrintExecutionResult выводит результат выполнения в консоль
func PrintExecutionResult(result *ExecutionResult)

// PrintStatistics выводит статистику в консоль
func PrintStatistics(stats *Statistics)
```

### 4.9. Конфигурация

```go
// LoadConfig загружает конфигурацию из файла
func LoadConfig(configPath string) (*Config, error)

// SaveConfig сохраняет конфигурацию в файл
func SaveConfig(config *Config, configPath string) error

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config

// GetConfig возвращает текущую конфигурацию клиента
func (c *Client) GetConfig() *Config

// UpdateConfig обновляет конфигурацию клиента
func (c *Client) UpdateConfig(config *Config) error
```

## 5. Примеры использования

### 5.1. Базовое использование

```go
package main

import (
    "fmt"
    "log"
    "time"
    
    "github.com/yourusername/getman"
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
        Timeout(30*time.Second, 30*time.Second).
        Build()

    result, err := client.ExecuteRequest(req)
    if err != nil {
        log.Fatal(err)
    }

    if result.Error != nil {
        fmt.Printf("Error: %v\n", result.Error)
        return
    }

    getman.PrintResponse(result.Response)
    fmt.Printf("Duration: %v\n", result.Duration)
}
```

### 5.2. Работа с коллекциями

```go
func main() {
    client, err := getman.NewClientWithDefaults()
    if err != nil {
        log.Fatal(err)
    }

    err = client.LoadEnvironment("production")
    if err != nil {
        log.Fatal(err)
    }

    result, err := client.ExecuteCollection("My API Collection")
    if err != nil {
        log.Fatal(err)
    }

    getman.PrintExecutionResult(result)
    fmt.Printf("\nStatistics:\n")
    getman.PrintStatistics(result.Statistics)
}
```

### 5.3. Выборочное выполнение запросов

```go
func main() {
    client, err := getman.NewClientWithDefaults()
    if err != nil {
        log.Fatal(err)
    }

    err = client.LoadEnvironment("production")
    if err != nil {
        log.Fatal(err)
    }

    itemNames := []string{"Get Users", "Create User"}
    result, err := client.ExecuteCollectionSelective("My API Collection", itemNames)
    if err != nil {
        log.Fatal(err)
    }

    getman.PrintExecutionResult(result)
}
```

### 5.4. Импорт из Postman

```go
func main() {
    client, err := getman.NewClientWithDefaults()
    if err != nil {
        log.Fatal(err)
    }

    collection, err := client.ImportFromPostman("./postman_collection.json")
    if err != nil {
        log.Fatal(err)
    }

    err = client.SaveCollection(collection)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Imported collection: %s\n", collection.Name)
}
```

### 5.5. Создание окружения

```go
func main() {
    client, err := getman.NewClientWithDefaults()
    if err != nil {
        log.Fatal(err)
    }

    env := &getman.Environment{
        Name: "staging",
        Variables: map[string]string{
            "baseUrl": "https://staging-api.example.com",
            "token": "staging_token_123",
            "apiKey": "staging_key_456",
        },
    }

    err = client.SaveEnvironment(env)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created environment: %s\n", env.Name)
}
```

### 5.6. Работа с переменными

```go
func main() {
    client, err := getman.NewClientWithDefaults()
    if err != nil {
        log.Fatal(err)
    }

    client.SetGlobalVariable("globalVar", "global_value")
    
    err = client.LoadEnvironment("production")
    if err != nil {
        log.Fatal(err)
    }

    value, ok := client.GetVariable("baseUrl")
    if ok {
        fmt.Printf("baseUrl: %s\n", value)
    }

    resolved, err := client.ResolveVariables("{{baseUrl}}/users/{{userId}}")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Resolved URL: %s\n", resolved)
}
```

### 5.7. Получение истории

```go
func main() {
    client, err := getman.NewClientWithDefaults()
    if err != nil {
        log.Fatal(err)
    }

    history, err := client.GetHistory(10)
    if err != nil {
        log.Fatal(err)
    }

    for i, exec := range history {
        fmt.Printf("%d. Request: %s %s\n", i+1, exec.Request.Method, exec.Request.URL)
        if exec.Error != nil {
            fmt.Printf("   Error: %v\n", exec.Error)
        } else {
            fmt.Printf("   Status: %d\n", exec.Response.StatusCode)
            fmt.Printf("   Duration: %v\n", exec.Duration)
        }
    }
}
```

## 6. Обработка ошибок

### 6.1. Типы ошибок

```go
var (
    ErrEnvironmentNotFound = errors.New("environment not found")
    ErrCollectionNotFound   = errors.New("collection not found")
    ErrVariableNotFound     = errors.New("variable not found")
    ErrInvalidRequest       = errors.New("invalid request")
    ErrInvalidURL           = errors.New("invalid URL")
    ErrRequestFailed        = errors.New("request failed")
    ErrStorageError         = errors.New("storage error")
)
```

### 6.2. Валидация

Все запросы валидируются перед выполнением:
- Проверка наличия всех переменных
- Проверка корректности URL
- Проверка обязательных полей (method, URL)
- Проверка формата тела запроса

## 7. Зависимости и библиотеки

### 7.1. Стандартная библиотека Go
- `net/http` - HTTP клиент
- `encoding/json` - работа с JSON
- `os`, `path/filepath` - работа с файловой системой
- `time` - работа со временем
- `fmt`, `strings` - форматирование и работа со строками

### 7.2. Рекомендуемые third-party библиотеки
- `github.com/fatih/color` - цветной вывод в консоль
- `github.com/google/uuid` - генерация UUID
- `gopkg.in/yaml.v3` - парсинг и работа с YAML

## 8. Ограничения MVP

На стадии MVP не реализуются:
- Пред/пост-скрипты
- Тесты и ассерты
- Динамические переменные
- Переменные из ответов предыдущих запросов
- Версионирование коллекций и окружений
- Параллельное выполнение запросов
- Retry механизм
- Плагины и расширения

## 9. Будущие улучшения

Возможные улучшения для следующих версий:
- Поддержка скриптов (JavaScript или Go)
- Система тестов и ассертов
- Динамические переменные (генерация значений)
- Извлечение переменных из ответов (JSONPath, regex)
- Параллельное выполнение запросов
- Retry механизм с настраиваемыми стратегиями
- Экспорт результатов в различные форматы (HTML, PDF)
- Веб-интерфейс для управления коллекциями
- Интеграция с CI/CD системами

