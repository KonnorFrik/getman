package getman

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/KonnorFrik/getman/collections"
	"github.com/KonnorFrik/getman/core"
	"github.com/KonnorFrik/getman/formatter"
	"github.com/KonnorFrik/getman/importer"
	"github.com/KonnorFrik/getman/storage"
	"github.com/KonnorFrik/getman/types"
	"github.com/KonnorFrik/getman/variables"
)

type Client struct {
	storage            *storage.FileStorage
	historyStorage     *storage.HistoryStorage
	logStorage         *storage.LogStorage
	variableStore      *variables.VariableStore
	variableResolver   *core.VariableResolver
	httpClient         *core.HTTPClient
	collectionExecutor *collections.CollectionExecutor
	env                *types.Environment
	// globalEnv          *types.Environment
	config *Config
}

func NewClient(basePath string) (*Client, error) {
	fileStorage, err := storage.NewFileStorage(basePath)
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()
	configPath := fileStorage.ConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		loadedConfig, err := LoadConfig(configPath)
		if err == nil {
			config = loadedConfig
		}
	}

	expandedBasePath, _ := storage.ExpandPath(config.Storage.BasePath)
	if expandedBasePath != "" {
		fileStorage, _ = storage.NewFileStorage(expandedBasePath)
	}

	historyStorage := storage.NewHistoryStorage(fileStorage)
	logStorage := storage.NewLogStorage(fileStorage)
	variableStore := variables.NewVariableStore()
	variableResolver := core.NewVariableResolver(variableStore)

	connectTimeout := config.Defaults.Timeout.Connect
	readTimeout := config.Defaults.Timeout.Read
	autoManageCookies := config.Defaults.Cookies.AutoManage
	httpClient := core.NewHTTPClient(connectTimeout, readTimeout, autoManageCookies)
	collectionExecutor := collections.NewCollectionExecutor(httpClient, variableResolver)

	return &Client{
		storage:            fileStorage,
		historyStorage:     historyStorage,
		logStorage:         logStorage,
		variableStore:      variableStore,
		variableResolver:   variableResolver,
		httpClient:         httpClient,
		collectionExecutor: collectionExecutor,
		config:             config,
	}, nil
}

// NewClientWithConfig создает новый клиент с конфигурацией из файла
func NewClientWithConfig(configPath string) (*Client, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return NewClient(config.Storage.BasePath)
}

// NewClientWithDefaults создает клиента с путями по умолчанию (~/.getman)
func NewClientWithDefaults() (*Client, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	basePath := filepath.Join(homeDir, ".getman")
	return NewClient(basePath)
}

func (c *Client) LoadEnvironment(name string) error {
	filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", name))
	env, err := variables.LoadEnvironmentFromFile(filePath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrEnvironmentNotFound, name)
	}

	c.env = env
	c.variableStore.SetEnvVars(env.Variables)
	return nil
}

func (c *Client) SaveEnvironment(env *types.Environment) error {
	filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", env.Name))
	return variables.SaveEnvironmentToFile(env, filePath)
}

func (c *Client) ListEnvironments() ([]string, error) {
	dir := c.storage.EnvironmentsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read environments directory: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) > 5 && name[len(name)-5:] == ".json" {
			names = append(names, name[:len(name)-5])
		}
	}

	return names, nil
}

func (c *Client) DeleteEnvironment(name string) error {
	filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", name))
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("%w: %s", ErrEnvironmentNotFound, name)
	}
	return nil
}

func (c *Client) GetCurrentEnvironment() *types.Environment {
	return c.env
}

func (c *Client) SetGlobalVariable(key, value string) {
	c.variableStore.SetGlobal(key, value)
}

func (c *Client) GetGlobalVariable(key string) (string, bool) {
	return c.variableStore.GetGlobal(key)
}

func (c *Client) GetVariable(key string) (string, bool) {
	return c.variableStore.Get(key)
}

func (c *Client) ResolveVariables(template string) (string, error) {
	return c.variableResolver.Resolve(template)
}

func (c *Client) LoadCollection(name string) (*types.Collection, error) {
	filePath := collections.GetCollectionPath(c.storage, name)
	collection, err := collections.LoadCollectionFromFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCollectionNotFound, name)
	}
	return collection, nil
}

func (c *Client) SaveCollection(collection *types.Collection) error {
	filePath := collections.GetCollectionPath(c.storage, collection.Name)
	return collections.SaveCollectionToFile(collection, filePath)
}

func (c *Client) ListCollections() ([]string, error) {
	dir := c.storage.CollectionsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read collections directory: %w", err)
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if len(name) > 5 && name[len(name)-5:] == ".json" {
			names = append(names, name[:len(name)-5])
		}
	}

	return names, nil
}

func (c *Client) DeleteCollection(name string) error {
	filePath := collections.GetCollectionPath(c.storage, name)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("%w: %s", ErrCollectionNotFound, name)
	}
	return nil
}

func (c *Client) ImportFromPostman(filePath string) (*types.Collection, error) {
	return importer.ImportFromPostman(filePath)
}

func (c *Client) ExportToPostman(collection *types.Collection, filePath string) error {
	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal collection: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write Postman collection file: %w", err)
	}

	return nil
}

func (c *Client) ExecuteRequest(req *types.Request) (*types.RequestExecution, error) {
	if err := c.ValidateRequest(req); err != nil {
		return nil, err
	}

	resolvedReq, err := c.resolveRequest(req)
	if err != nil {
		return nil, err
	}

	startTime := time.Now()
	response, err := c.httpClient.Execute(resolvedReq)
	duration := time.Since(startTime)

	execution := &types.RequestExecution{
		Request:   resolvedReq,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		execution.Error = err.Error()
	} else {
		execution.Response = response
	}

	return execution, nil
}

func (c *Client) ExecuteCollection(collectionName string) (*types.ExecutionResult, error) {
	collection, err := c.LoadCollection(collectionName)
	if err != nil {
		return nil, err
	}

	envName := ""
	if c.env != nil {
		envName = c.env.Name
	}

	return c.collectionExecutor.ExecuteCollection(collection, envName)
}

func (c *Client) ExecuteCollectionSelective(collectionName string, itemNames []string) (*types.ExecutionResult, error) {
	collection, err := c.LoadCollection(collectionName)
	if err != nil {
		return nil, err
	}

	envName := ""
	if c.env != nil {
		envName = c.env.Name
	}

	return c.collectionExecutor.ExecuteCollectionSelective(collection, envName, itemNames)
}

func (c *Client) ValidateRequest(req *types.Request) error {
	if req.Method == "" {
		return fmt.Errorf("%w: method is required", ErrInvalidRequest)
	}

	if req.URL == "" {
		return fmt.Errorf("%w: URL is required", ErrInvalidRequest)
	}

	if _, err := url.Parse(req.URL); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	if err := c.variableResolver.ValidateVariables(req.URL); err != nil {
		return err
	}

	if err := c.variableResolver.ValidateVariablesInMap(req.Headers); err != nil {
		return err
	}

	if req.Body != nil && len(req.Body.Content) > 0 {
		if err := c.variableResolver.ValidateVariables(string(req.Body.Content)); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) GetHistory(limit int) ([]*types.RequestExecution, error) {
	return c.historyStorage.GetHistory(limit)
}

func (c *Client) GetLastExecution() (*types.ExecutionResult, error) {
	return c.historyStorage.GetLast()
}

func (c *Client) GetLogs() ([]byte, error) {
	return c.logStorage.GetLast()
}

func (c *Client) ClearHistory() error {
	return c.historyStorage.Clear()
}

func (c *Client) SaveHistory(result *types.ExecutionResult) error {
	return c.historyStorage.Save(result)
}

func (c *Client) SaveLogs(logs []types.LogEntry) error {
	return c.logStorage.Save(logs)
}

func (c *Client) GetConfig() *Config {
	return c.config
}

func (c *Client) UpdateConfig(config *Config) error {
	if err := validateConfig(config); err != nil {
		return err
	}

	c.config = config
	configPath := c.storage.ConfigPath()
	return SaveConfig(config, configPath)
}

func (c *Client) resolveRequest(req *types.Request) (*types.Request, error) {
	resolvedURL, err := c.variableResolver.Resolve(req.URL)
	if err != nil {
		return nil, err
	}

	resolvedHeaders, err := c.variableResolver.ResolveMap(req.Headers)
	if err != nil {
		return nil, err
	}

	resolvedReq := &types.Request{
		Method:  req.Method,
		URL:     resolvedURL,
		Headers: resolvedHeaders,
		Body:    req.Body,
		Auth:    req.Auth,
		Timeout: req.Timeout,
		Cookies: req.Cookies,
	}

	if req.Body != nil && len(req.Body.Content) > 0 {
		resolvedBodyContent, err := c.variableResolver.Resolve(string(req.Body.Content))
		if err != nil {
			return nil, err
		}
		resolvedReq.Body = &types.RequestBody{
			Type:        req.Body.Type,
			Content:     []byte(resolvedBodyContent),
			ContentType: req.Body.ContentType,
		}
	}

	if req.Auth != nil {
		resolvedAuth := &types.Auth{
			Type:     req.Auth.Type,
			Username: req.Auth.Username,
			Password: req.Auth.Password,
			Token:    req.Auth.Token,
			APIKey:   req.Auth.APIKey,
			KeyName:  req.Auth.KeyName,
			Location: req.Auth.Location,
		}

		if req.Auth.Username != "" {
			resolvedAuth.Username, err = c.variableResolver.Resolve(req.Auth.Username)
			if err != nil {
				return nil, err
			}
		}
		if req.Auth.Password != "" {
			resolvedAuth.Password, err = c.variableResolver.Resolve(req.Auth.Password)
			if err != nil {
				return nil, err
			}
		}
		if req.Auth.Token != "" {
			resolvedAuth.Token, err = c.variableResolver.Resolve(req.Auth.Token)
			if err != nil {
				return nil, err
			}
		}
		if req.Auth.APIKey != "" {
			resolvedAuth.APIKey, err = c.variableResolver.Resolve(req.Auth.APIKey)
			if err != nil {
				return nil, err
			}
		}

		resolvedReq.Auth = resolvedAuth
	}

	return resolvedReq, nil
}

func NewRequestBuilder() *core.RequestBuilder {
	return core.NewRequestBuilder()
}

func FormatResponse(resp *types.Response) string {
	return formatter.FormatResponse(resp)
}

func FormatRequest(req *types.Request) string {
	return formatter.FormatRequest(req)
}

func FormatExecutionResult(result *types.ExecutionResult) string {
	return formatter.FormatExecutionResult(result)
}

func FormatStatistics(stats *types.Statistics) string {
	return formatter.FormatStatistics(stats)
}

func PrintResponse(resp *types.Response) {
	formatter.PrintResponse(resp)
}

func PrintRequest(req *types.Request) {
	formatter.PrintRequest(req)
}

func PrintExecutionResult(result *types.ExecutionResult) {
	formatter.PrintExecutionResult(result)
}

func PrintStatistics(stats *types.Statistics) {
	formatter.PrintStatistics(stats)
}
