/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package fixture

import (
	"net/http"
	"time"

	"github.com/KonnorFrik/getman/collections"
	"github.com/KonnorFrik/getman/environment"
	"github.com/KonnorFrik/getman/types"
)

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

func GetTestPostmanCollectionJSON() string {
	return `{
		"info": {
			"name": "Test Collection",
			"description": "Test collection description",
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
		},
		"item": [
			{
				"name": "Get Users",
				"request": {
					"method": "GET",
					"header": [
						{
							"key": "Accept",
							"value": "application/json"
						}
					],
					"url": {
						"raw": "{{baseUrl}}/users",
						"host": ["{{baseUrl}}"],
						"path": ["users"]
					},
					"auth": {
						"type": "bearer",
						"bearer": [
							{
								"key": "token",
								"value": "{{token}}",
								"type": "string"
							}
						]
					}
				}
			},
			{
				"name": "Create User",
				"request": {
					"method": "POST",
					"header": [
						{
							"key": "Content-Type",
							"value": "application/json"
						}
					],
					"body": {
						"mode": "raw",
						"raw": "{\"name\": \"{{userName}}\", \"email\": \"{{userEmail}}\"}"
					},
					"url": {
						"raw": "{{baseUrl}}/users",
						"host": ["{{baseUrl}}"],
						"path": ["users"]
					}
				}
			}
		]
	}`
}

func GetTestEnvironmentJSON() string {
	return `{
		"name": "test",
		"variables": {
			"baseUrl": "http://localhost:8080",
			"token": "testtoken123",
			"userName": "Test User",
			"userEmail": "test@example.com"
		}
	}`
}

func GetTestCollectionJSON() string {
	return `{
		"name": "Test Collection",
		"description": "Test collection description",
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
			},
			{
				"name": "Create User",
				"request": {
					"method": "POST",
					"url": "{{baseUrl}}/users",
					"headers": {
						"Content-Type": "application/json"
					},
					"body": {
						"type": "json",
						"content": "{\"name\": \"{{userName}}\", \"email\": \"{{userEmail}}\"}",
						"content_type": "application/json"
					}
				}
			}
		]
	}`
}

func GetTestConfigYAML() string {
	return `storage:
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
`
}

func GetTestRequestItem(name string, req *types.Request) *types.RequestItem {
	return &types.RequestItem{
		Name:    name,
		Request: req,
	}
}

func GetTestCollectionWithItems(name string, items []*types.RequestItem) *collections.Collection {
	return &collections.Collection{
		Name:        name,
		Description: "Test collection",
		Items:       items,
	}
}

func GetTestEnvironmentWithVars(name string, vars map[string]string) *environment.Environment {
	return &environment.Environment{
		Name:      name,
		Variables: vars,
	}
}
