/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package errors

import "errors"

var (
	// ErrEnvironmentNotFound is returned when an environment is not found.
	ErrEnvironmentNotFound = errors.New("environment not found")
	// ErrCollectionNotFound is returned when a collection is not found.
	ErrCollectionNotFound  = errors.New("collection not found")
	// ErrVariableNotFound is returned when a variable is not found in the environment.
	ErrVariableNotFound    = errors.New("variable not found")
	// ErrInvalidRequest is returned when a request is invalid.
	ErrInvalidRequest      = errors.New("invalid request")
	// ErrInvalidURL is returned when a URL is invalid.
	ErrInvalidURL          = errors.New("invalid URL")
	// ErrRequestFailed is returned when an HTTP request fails.
	ErrRequestFailed       = errors.New("request failed")
	// ErrStorageError is returned when a storage operation fails.
	ErrStorageError        = errors.New("storage error")
	// ErrInvalidArgument is returned when an invalid argument is provided.
	ErrInvalidArgument     = errors.New("invalid argument")
)
