/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package getman

import (
	client "github.com/KonnorFrik/getman/client"
	"github.com/KonnorFrik/getman/core"
	"github.com/KonnorFrik/getman/types"
)

// Core types.
type (
	Request          = types.Request
	RequestBody      = types.RequestBody
	Auth             = types.Auth
	Timeout          = types.Timeout
	CookieSettings   = types.CookieSettings
	Response         = types.Response
	RequestItem      = types.RequestItem
	RequestExecution = types.RequestExecution
	ExecutionResult  = types.ExecutionResult
	Statistics       = types.Statistics
	LogEntry         = types.LogEntry
)

// Collection and environment types.
type (
	Collection  = client.Collection
	Environment = client.Environment
)

// Client and configuration types.
type (
	Client         = client.Client
	Config         = client.Config
	StorageConfig  = client.StorageConfig
	DefaultsConfig = client.DefaultsConfig
	TimeoutConfig  = client.TimeoutConfig
	CookiesConfig  = client.CookiesConfig
	LoggingConfig  = client.LoggingConfig
	RequestBuilder = core.RequestBuilder
)

// Errors returned by the library.
var (
	ErrEnvironmentNotFound = client.ErrEnvironmentNotFound
	ErrCollectionNotFound  = client.ErrCollectionNotFound
	ErrVariableNotFound    = client.ErrVariableNotFound
	ErrInvalidRequest      = client.ErrInvalidRequest
	ErrInvalidURL          = client.ErrInvalidURL
	ErrRequestFailed       = client.ErrRequestFailed
	ErrStorageError        = client.ErrStorageError
	ErrInvalidArgument     = client.ErrInvalidArgument
	ErrPlaceholderNotFound = client.ErrPlaceholderNotFound
)

// NewClient creates a new Client instance with the specified base path for storage.
func NewClient(basePath string) (*Client, error) {
	return client.NewClient(basePath)
}

// NewClientWithConfig creates a new Client instance using configuration from the specified file.
func NewClientWithConfig(configPath string) (*Client, error) {
	return client.NewClientWithConfig(configPath)
}

// NewClientWithDefaults creates a Client instance with default paths (~/.getman).
func NewClientWithDefaults() (*Client, error) {
	return client.NewClientWithDefaults()
}

// NewRequestBuilder creates a new RequestBuilder instance for constructing HTTP requests.
func NewRequestBuilder() *RequestBuilder {
	return core.NewRequestBuilder()
}

// LoadConfig loads configuration from a YAML file.
func LoadConfig(configPath string) (*Config, error) {
	return client.LoadConfig(configPath)
}

// SaveConfig saves configuration to a YAML file.
func SaveConfig(config *Config, configPath string) error {
	return client.SaveConfig(config, configPath)
}

// DefaultConfig returns a configuration with default values.
func DefaultConfig() *Config {
	return client.DefaultConfig()
}

// FormatResponse formats a response as a string for display.
func FormatResponse(resp *Response) string {
	return client.FormatResponse(resp)
}

// FormatRequest formats a request as a string for display.
func FormatRequest(req *Request) string {
	return client.FormatRequest(req)
}

// FormatExecutionResult formats an execution result as a string for display.
func FormatExecutionResult(result *ExecutionResult) string {
	return client.FormatExecutionResult(result)
}

// FormatStatistics formats statistics as a string for display.
func FormatStatistics(stats *Statistics) string {
	return client.FormatStatistics(stats)
}

// PrintResponse prints a formatted response to stdout.
func PrintResponse(resp *Response) {
	client.PrintResponse(resp)
}

// PrintRequest prints a formatted request to stdout.
func PrintRequest(req *Request) {
	client.PrintRequest(req)
}

// PrintExecutionResult prints a formatted execution result to stdout.
func PrintExecutionResult(result *ExecutionResult) {
	client.PrintExecutionResult(result)
}

// PrintRequestExecution prints a formatted request execution to stdout.
func PrintRequestExecution(req *RequestExecution, index int) {
	client.PrintRequestExecution(req, index)
}

// PrintStatistics prints formatted statistics to stdout.
func PrintStatistics(stats *Statistics) {
	client.PrintStatistics(stats)
}
