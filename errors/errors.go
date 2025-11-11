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
)
