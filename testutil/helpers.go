package testutil

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/KonnorFrik/getman/collections"
	"github.com/KonnorFrik/getman/environment"
	"github.com/KonnorFrik/getman/types"
)

func CreateTempDir() (string, error) {
	return os.MkdirTemp("", "getman_test_*")
}

func CleanupTempDir(dir string) error {
	return os.RemoveAll(dir)
}

func CreateTestCollection(name string, items []*types.RequestItem) *collections.Collection {
	return &collections.Collection{
		Name:        name,
		Description: "Test collection",
		Items:       items,
	}
}

func CreateTestEnvironment(name string, variables map[string]string) *environment.Environment {
	return &environment.Environment{
		Name:      name,
		Variables: variables,
	}
}

func CreateTestRequest(method, url string) *types.Request {
	return &types.Request{
		Method:  method,
		URL:     url,
		Headers: make(map[string]string),
	}
}

func CreateTestRequestWithHeaders(method, url string, headers map[string]string) *types.Request {
	return &types.Request{
		Method:  method,
		URL:     url,
		Headers: headers,
	}
}

func CreateTestRequestWithBody(method, url string, body *types.RequestBody) *types.Request {
	return &types.Request{
		Method:  method,
		URL:     url,
		Headers: make(map[string]string),
		Body:    body,
	}
}

func CreateTestRequestWithAuth(method, url string, auth *types.Auth) *types.Request {
	return &types.Request{
		Method:  method,
		URL:     url,
		Headers: make(map[string]string),
		Auth:    auth,
	}
}

func CreateTestResponse(statusCode int, body []byte) *types.Response {
	return &types.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Headers:    make(map[string][]string),
		Body:       body,
		Duration:   time.Millisecond * 100,
		Size:       int64(len(body)),
	}
}

func CreateTestJSONBody(content string) *types.RequestBody {
	return &types.RequestBody{
		Type:        "json",
		Content:     []byte(content),
		ContentType: "application/json",
	}
}

func CreateTestXMLBody(content string) *types.RequestBody {
	return &types.RequestBody{
		Type:        "xml",
		Content:     []byte(content),
		ContentType: "application/xml",
	}
}

func CreateTestRawBody(content string, contentType string) *types.RequestBody {
	return &types.RequestBody{
		Type:        "raw",
		Content:     []byte(content),
		ContentType: contentType,
	}
}

func CreateTestBasicAuth(username, password string) *types.Auth {
	return &types.Auth{
		Type:     "basic",
		Username: username,
		Password: password,
	}
}

func CreateTestBearerAuth(token string) *types.Auth {
	return &types.Auth{
		Type:  "bearer",
		Token: token,
	}
}

func CreateTestAPIKeyAuth(keyName, keyValue, location string) *types.Auth {
	return &types.Auth{
		Type:     "apikey",
		APIKey:   keyValue,
		KeyName:  keyName,
		Location: location,
	}
}

func WriteTestFile(dir, filename, content string) error {
	filePath := filepath.Join(dir, filename)
	return os.WriteFile(filePath, []byte(content), 0644)
}

func ReadTestFile(dir, filename string) (string, error) {
	filePath := filepath.Join(dir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
