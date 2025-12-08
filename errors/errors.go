/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package errors

import "errors"

var (
	ErrEnvironmentNotFound = errors.New("environment not found")
	ErrCollectionNotFound  = errors.New("collection not found")
	ErrVariableNotFound    = errors.New("variable not found")
	ErrInvalidRequest      = errors.New("invalid request")
	ErrInvalidURL          = errors.New("invalid URL")
	ErrRequestFailed       = errors.New("request failed")
	ErrStorageError        = errors.New("storage error")

	ErrInvalidArgument     = errors.New("invalid argument")
)
