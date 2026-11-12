/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package getman

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
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
	variableResolver   *core.VariableResolver
	placeholders       *environment.PlaceholderStore
	config             *Config
}

const globalEnvName = "global"

// NewClient creates a new Client instance with the specified base path for storage.
// It initializes all required components including file storage, history storage,
// log storage, HTTP client, and variable resolver.
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

	return newClient(fileStorage, config)
}

// NewClientWithConfig creates a new Client instance using configuration from the specified file.
func NewClientWithConfig(configPath string) (*Client, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	fileStorage, err := storage.NewFileStorage(config.Storage.BasePath)
	if err != nil {
		return nil, err
	}

	return newClient(fileStorage, config)
}

func newClient(fileStorage *storage.FileStorage, config *Config) (*Client, error) {
	globalEnv := environment.NewEnvironment(globalEnvName)
	variableResolver, err := core.NewVariableResolver(globalEnv, nil)
	if err != nil {
		return nil, err
	}

	httpClient := core.NewHTTPClient(
		config.Defaults.Timeout.Connect,
		config.Defaults.Timeout.Read,
		config.Defaults.Cookies.AutoManage,
	)

	placeholderStore := environment.NewPlaceholderStore(fileStorage.PlaceholdersPath())
	_ = placeholderStore.Load()

	client := &Client{
		storage:            fileStorage,
		historyStorage:     storage.NewHistoryStorage(fileStorage),
		logStorage:         storage.NewLogStorage(fileStorage),
		httpClient:         httpClient,
		collectionExecutor: collections.NewCollectionExecutor(httpClient, variableResolver),
		variableResolver:   variableResolver,
		placeholders:       placeholderStore,
		config:             config,
	}

	client.collectionExecutor.SetPathResolver(client.resolvePathPlaceholders)
	_ = client.LoadGlobalEnvironment()

	return client, nil
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

// LoadEnvironment loads a local environment by name and makes it active.
// It is an alias for LoadLocalEnvironment.
func (c *Client) LoadEnvironment(name string) error {
	return c.LoadLocalEnvironment(name)
}

// LoadLocalEnvironment loads a local environment by name from storage.
func (c *Client) LoadLocalEnvironment(name string) error {
	filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", name))
	env, err := environment.NewEnvironmentFromFile(filePath)

	if err != nil {
		return fmt.Errorf("%w: %s", ErrEnvironmentNotFound, name)
	}

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

	c.variableResolver.SetGlobal(env)
	return nil
}

// SaveEnvironments saves both local and global environments to storage.
func (c *Client) SaveEnvironments() error {
	var (
		localEnv  = c.variableResolver.GetLocal()
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

// SaveEnvironment saves the given environment to storage and makes it the active local environment.
func (c *Client) SaveEnvironment(env *Environment) error {
	if env == nil {
		return fmt.Errorf("%w: 'env' is nil", ErrInvalidArgument)
	}

	if err := c.saveEnvironmentFile(env); err != nil {
		return err
	}

	c.variableResolver.SetLocal(env)
	return nil
}

func (c *Client) saveEnvironmentFile(env *Environment) error {
	filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", env.Name))
	return env.Save(filePath)
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

	sort.Strings(names)
	return names, nil
}

// DeleteEnvironment deletes an environment by name.
func (c *Client) DeleteEnvironment(name string) error {
	filePath := filepath.Join(c.storage.EnvironmentsDir(), fmt.Sprintf("%s.json", name))

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("%w: %s", ErrEnvironmentNotFound, name)
	}

	if local := c.variableResolver.GetLocal(); local != nil && local.Name == name {
		c.variableResolver.SetLocal(nil)
	}

	return nil
}

// GetCurrentEnvironment returns the currently loaded local environment.
func (c *Client) GetCurrentEnvironment() *environment.Environment {
	return c.variableResolver.GetLocal()
}

// SetGlobalVariable sets a variable in the global environment and persists it.
func (c *Client) SetGlobalVariable(key, value string) {
	env := c.variableResolver.GetGlobal()

	if env != nil {
		env.Set(key, value)
		_ = c.saveEnvironmentFile(env)
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

// SetVariable sets a variable in the current local environment and persists it.
func (c *Client) SetVariable(key, value string) {
	env := c.variableResolver.GetLocal()

	if env != nil {
		env.Set(key, value)
		_ = c.saveEnvironmentFile(env)
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

	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrCollectionNotFound, name)
	}

	if collection.EnvName != "" {
		err = c.LoadLocalEnvironment(collection.EnvName)

		if err != nil {
			return collection, fmt.Errorf("%w: %s", ErrEnvironmentNotFound, collection.EnvName)
		}
	}

	return collection, nil
}

// SaveCollection saves a collection to storage.
// Path placeholders found in request URLs are registered in placeholders storage.
func (c *Client) SaveCollection(collection *collections.Collection) error {
	if collection == nil {
		return fmt.Errorf("%w: 'collection' is nil", ErrInvalidArgument)
	}

	filePath := collections.GetCollectionPath(c.storage, collection.Name)
	if err := collections.SaveCollectionToFile(collection, filePath); err != nil {
		return err
	}

	foundPlaceholders := false

	for _, item := range collection.Items {
		if item == nil || item.Request == nil {
			continue
		}

		for _, name := range environment.ExtractPathPlaceholders(item.Request.URL) {
			c.placeholders.Ensure(name)
			foundPlaceholders = true
		}
	}

	if foundPlaceholders {
		return c.placeholders.Save()
	}

	return nil
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

	sort.Strings(names)
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
// Path placeholders found in imported requests are registered in placeholders storage.
func (c *Client) ImportFromPostman(filePath string) (*collections.Collection, error) {
	collection, err := importer.ImportFromPostman(filePath)
	if err != nil {
		return nil, err
	}

	foundPlaceholders := false

	for _, item := range collection.Items {
		if item == nil || item.Request == nil {
			continue
		}

		for _, name := range environment.ExtractPathPlaceholders(item.Request.URL) {
			c.placeholders.Ensure(name)
			foundPlaceholders = true
		}
	}

	if foundPlaceholders {
		if err := c.placeholders.Save(); err != nil {
			return collection, err
		}
	}

	return collection, nil
}

// ExportToPostman exports a collection to a Postman Collection v2.1 JSON file.
func (c *Client) ExportToPostman(collection *collections.Collection, filePath string) error {
	if collection == nil {
		return fmt.Errorf("%w: 'collection' is nil", ErrInvalidArgument)
	}

	return importer.ExportToPostman(collection, filePath)
}

// SetPlaceholder sets a path placeholder value and persists it.
func (c *Client) SetPlaceholder(key, value string) error {
	c.placeholders.Set(key, value)
	return c.placeholders.Save()
}

// GetPlaceholder returns a path placeholder value.
func (c *Client) GetPlaceholder(key string) (string, bool) {
	return c.placeholders.Get(key)
}

// DeletePlaceholder removes a path placeholder and persists the change.
func (c *Client) DeletePlaceholder(key string) error {
	c.placeholders.Delete(key)
	return c.placeholders.Save()
}

// ListPlaceholders returns a copy of all path placeholders.
func (c *Client) ListPlaceholders() map[string]string {
	return c.placeholders.Values()
}

// ResolvePathPlaceholders replaces path placeholders in a URL with stored values.
func (c *Client) ResolvePathPlaceholders(rawURL string) (string, error) {
	return c.placeholders.ResolvePath(rawURL)
}

func (c *Client) resolvePathPlaceholders(rawURL string) (string, error) {
	return c.placeholders.ResolvePath(rawURL)
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

	resolvedURL, err = c.placeholders.ResolvePath(resolvedURL)
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

func PrintRequestExecution(req *types.RequestExecution, index int) {
	formatter.PrintRequestExecution(req, index)
}

func PrintStatistics(stats *types.Statistics) {
	formatter.PrintStatistics(stats)
}
