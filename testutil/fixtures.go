package testutil

import (
	"github.com/KonnorFrik/getman/collections"
	"github.com/KonnorFrik/getman/environment"
	"github.com/KonnorFrik/getman/types"
)

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
