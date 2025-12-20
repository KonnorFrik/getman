/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
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
	"github.com/KonnorFrik/getman/environment"
	"github.com/KonnorFrik/getman/formatter"
	"github.com/KonnorFrik/getman/importer"
	"github.com/KonnorFrik/getman/storage"
	"github.com/KonnorFrik/getman/types"
)

// Client represents the main client for interacting with the getman library.
// It provides methods for managing collections, environments, executing requests,
// and handling storage operations.
type Client struct {
	storage            *storage.FileStorage
	historyStorage     *storage.HistoryStorage
	logStorage         *storage.LogStorage
	httpClient         *core.HTTPClient
	collectionExecutor *collections.CollectionExecutor
	variableResolver *core.VariableResolver
	config           *Config
}

const globalEnvName = "global"

// NewClient creates a new Client instance with the specified base path for storage.
// It initializes all required components including file storage, history storage,
// log storage, HTTP client, and variable resolver.
func NewClient(basePath string) (*Client, error) {
	var (
		client Client
		globalEnv = environment.NewEnvironment(globalEnvName)
	)

	fileStorage, err := storage.NewFileStorage(basePath)

	if err != nil {
		return nil, err
	}

	client.storage = fileStorage
	config := DefaultConfig()
	configPath := fileStorage.ConfigPath()

	if _, err := os.Stat(configPath); err == nil {
		loadedConfig, err := LoadConfig(configPath)

		if err == nil {
			config = loadedConfig
		}
	}

	historyStorage := storage.NewHistoryStorage(fileStorage)
	logStorage := storage.NewLogStorage(fileStorage)
	variableResolver, err := core.NewVariableResolver(globalEnv, nil)

	if err != nil {
		return nil, err
	}

	client.LoadGlobalEnvironment()

	// if err != nil {
	// 	variableResolver.SetGlobal(globalEnv)
	// }

	connectTimeout := config.Defaults.Timeout.Connect
	readTimeout := config.Defaults.Timeout.Read
	autoManageCookies := config.Defaults.Cookies.AutoManage
	httpClient := core.NewHTTPClient(connectTimeout, readTimeout, autoManageCookies)
	collectionExecutor := collections.NewCollectionExecutor(httpClient, variableResolver)

	client.historyStorage = historyStorage
	client.logStorage = logStorage
	client.variableResolver = variableResolver
	client.httpClient = httpClient
	client.collectionExecutor = collectionExecutor
	client.config = config

	return &client, nil
}

// NewClientWithConfig creates a new Client instance using configuration from the specified file.
func NewClientWithConfig(configPath string) (*Client, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	return NewClient(config.Storage.BasePath)
}

// NewClientWithDefaults creates a Client instance with default paths (~/.getman).
func NewClientWithDefaults() (*Client, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	basePath := filepath.Join(homeDir, ".getman")
	return NewClient(basePath)
}

// LoadLocalEnvironment loads a local environment by name from storage.
func (c *Client) LoadLocalEnvironment(name string) error {
	filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", name))
	env, err := environment.NewEnvironmentFromFile(filePath)

	if err != nil {
		return fmt.Errorf("%w: %s", ErrEnvironmentNotFound, name)
	}

	// c.env = env
	c.variableResolver.SetLocal(env)
	return nil
}

// LoadGlobalEnvironment loads the global environment from storage.
func (c *Client) LoadGlobalEnvironment() error {
	filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", globalEnvName))
	env, err := environment.NewEnvironmentFromFile(filePath)

	if err != nil {
		return fmt.Errorf("%w: %s", ErrEnvironmentNotFound, globalEnvName)
	}

	// c.globalEnv = env
	c.variableResolver.SetGlobal(env)
	return nil
}

// SaveEnvironments saves both local and global environments to storage.
func (c *Client) SaveEnvironments() error {
	var (
		localEnv = c.variableResolver.GetLocal()
		globalEnv = c.variableResolver.GetGlobal()
	)

	if localEnv != nil {
		filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", localEnv.Name))
		err := localEnv.Save(filePath)

		if err != nil {
			return err
		}
	}

	if globalEnv != nil {
		filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", globalEnv.Name))
		err := globalEnv.Save(filePath)

		if err != nil {
			return err
		}
	}

	return nil
}

// SaveEnvironment saves the given environment to storage.
func (c *Client) SaveEnvironment(env *Environment) error {
	if env == nil {
		return fmt.Errorf("%w: 'env' is nil", ErrInvalidArgument)
	}

	filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", env.Name))
	err := env.Save(filePath)

	if err != nil {
		return err
	}

	// globalFilePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", c.globalEnv.Name))
	// err = c.globalEnv.Save(globalFilePath)
	//
	// if err != nil {
	// 	return err
	// }

	return nil
}

// ListEnvironments returns a list of all available environment names.
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

// DeleteEnvironment deletes an environment by name.
func (c *Client) DeleteEnvironment(name string) error {
	filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", name))

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("%w: %s", ErrEnvironmentNotFound, name)
	}

	return nil
}

// GetCurrentEnvironment returns the currently loaded local environment.
func (c *Client) GetCurrentEnvironment() *environment.Environment {
	return c.variableResolver.GetLocal()
}

// SetGlobalVariable sets a variable in the global environment.
func (c *Client) SetGlobalVariable(key, value string) {
	env := c.variableResolver.GetGlobal()

	if env != nil {
		env.Set(key, value)
	}
}

// GetGlobalVariable retrieves a variable value from the global environment.
func (c *Client) GetGlobalVariable(key string) (string, bool) {
	env := c.variableResolver.GetGlobal()
	var (
		result string
		ok     bool
	)

	if env != nil {
		result, ok = env.Get(key)
	}

	return result, ok
}

// SetVariable sets a variable in the current local environment.
func (c *Client) SetVariable(key, value string) {
	env := c.variableResolver.GetLocal()

	if env != nil {
		env.Set(key, value)
	}
}

// GetVariable retrieves a variable value from the current local environment
// If key not exists, retrieves a variable value from the global environment.
func (c *Client) GetVariable(key string) (string, bool) {
	env := c.variableResolver.GetLocal()
	var (
		result string
		ok     bool
	)

	if env != nil {
		result, ok = env.Get(key)
	}

	if !ok {
		env = c.variableResolver.GetGlobal()

		if env != nil {
			result, ok = env.Get(key)
		}
	}

	return result, ok
}

// ResolveVariables resolves variables in the given template string using the current environment.
func (c *Client) ResolveVariables(template string) (string, error) {
	return c.variableResolver.Resolve(template)
}

// LoadCollection loads a collection by name from storage.
func (c *Client) LoadCollection(name string) (*collections.Collection, error) {
	filePath := collections.GetCollectionPath(c.storage, name)
	collection, err := collections.LoadCollectionFromFile(filePath)

	fmt.Printf("LoadCollection: load: %s\n", name)

	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCollectionNotFound, name)
	}

	fmt.Printf("LoadCollection: loaded: %+v\n", collection)

	if collection.EnvName != "" {
		err = c.LoadLocalEnvironment(collection.EnvName)

		if err != nil {
			return collection, fmt.Errorf("%w: %s", ErrEnvironmentNotFound, collection.EnvName)
		}
	}

	return collection, nil
}

// SaveCollection saves a collection to storage.
func (c *Client) SaveCollection(collection *collections.Collection) error {
	filePath := collections.GetCollectionPath(c.storage, collection.Name)
	return collections.SaveCollectionToFile(collection, filePath)
}

// ListCollections returns a list of all available collection names.
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

// DeleteCollection deletes a collection by name.
func (c *Client) DeleteCollection(name string) error {
	filePath := collections.GetCollectionPath(c.storage, name)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("%w: %s", ErrCollectionNotFound, name)
	}
	return nil
}

// ImportFromPostman imports a Postman collection from the specified file path.
func (c *Client) ImportFromPostman(filePath string) (*collections.Collection, error) {
	return importer.ImportFromPostman(filePath)
}

// ExportToPostman exports a collection to a Postman-compatible JSON file.
func (c *Client) ExportToPostman(collection *collections.Collection, filePath string) error {
	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal collection: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write Postman collection file: %w", err)
	}

	return nil
}

// ExecuteRequest executes a single HTTP request and returns the execution result.
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

// ExecuteCollection executes all requests in a collection by name.
func (c *Client) ExecuteCollection(collectionName string) (*types.ExecutionResult, error) {
	collection, err := c.LoadCollection(collectionName)

	if err != nil {
		return nil, err
	}

	localEnvName := ""
	localEnv := c.variableResolver.GetLocal()

	if localEnv != nil {
		localEnvName = localEnv.Name
	}

	return c.collectionExecutor.ExecuteCollection(collection, localEnvName)
}

func (c *Client) ExecuteCollectionAsync(collectionName string) <-chan *types.RequestExecution {
	collection, err := c.LoadCollection(collectionName)

	if err != nil {
		return nil
	}

	localEnvName := ""
	localEnv := c.variableResolver.GetLocal()

	if localEnv != nil {
		localEnvName = localEnv.Name
	}

	return c.collectionExecutor.ExecuteCollectionAsync(collection, localEnvName)
}

// ExecuteCollectionSelective executes only the specified requests from a collection.
func (c *Client) ExecuteCollectionSelective(collectionName string, itemNames []string) (*types.ExecutionResult, error) {
	collection, err := c.LoadCollection(collectionName)

	if err != nil {
		return nil, err
	}

	localEnvName := ""
	localEnv := c.variableResolver.GetLocal()

	if localEnv != nil {
		localEnvName = localEnv.Name
	}

	return c.collectionExecutor.ExecuteCollectionSelective(collection, localEnvName, itemNames)
}

// ValidateRequest validates a request before execution, checking method, URL, and variables.
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

// GetHistory retrieves request execution history up to the specified limit.
func (c *Client) GetHistory(limit int) ([]*types.RequestExecution, error) {
	return c.historyStorage.GetHistory(limit)
}

// GetLastExecution retrieves the most recent execution result.
func (c *Client) GetLastExecution() (*types.ExecutionResult, error) {
	return c.historyStorage.GetLast()
}

// GetLogs retrieves the most recent log entries as JSON bytes.
func (c *Client) GetLogs() ([]byte, error) {
	return c.logStorage.GetLast()
}

// ClearHistory removes all stored execution history.
func (c *Client) ClearHistory() error {
	return c.historyStorage.Clear()
}

// SaveHistory saves an execution result to history storage.
func (c *Client) SaveHistory(result *types.ExecutionResult) error {
	return c.historyStorage.Save(result)
}

// SaveLogs saves log entries to storage.
func (c *Client) SaveLogs(logs []types.LogEntry) error {
	return c.logStorage.Save(logs)
}

// GetConfig returns the current client configuration.
func (c *Client) GetConfig() *Config {
	return c.config
}

// UpdateConfig updates the client configuration and saves it to storage.
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

// NewRequestBuilder creates a new RequestBuilder instance for constructing HTTP requests.
func NewRequestBuilder() *core.RequestBuilder {
	return core.NewRequestBuilder()
}

// FormatResponse formats a response as a string for display.
func FormatResponse(resp *types.Response) string {
	return formatter.FormatResponse(resp)
}

// FormatRequest formats a request as a string for display.
func FormatRequest(req *types.Request) string {
	return formatter.FormatRequest(req)
}

// FormatExecutionResult formats an execution result as a string for display.
func FormatExecutionResult(result *types.ExecutionResult) string {
	return formatter.FormatExecutionResult(result)
}

// FormatStatistics formats statistics as a string for display.
func FormatStatistics(stats *types.Statistics) string {
	return formatter.FormatStatistics(stats)
}

// PrintResponse prints a formatted response to stdout.
func PrintResponse(resp *types.Response) {
	formatter.PrintResponse(resp)
}

// PrintRequest prints a formatted request to stdout.
func PrintRequest(req *types.Request) {
	formatter.PrintRequest(req)
}

func PrintExecutionResult(result *types.ExecutionResult) {
	formatter.PrintExecutionResult(result)
}

func PrintStatistics(stats *types.Statistics) {
	formatter.PrintStatistics(stats)
}
