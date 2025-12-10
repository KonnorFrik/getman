/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>

*/
package importer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/KonnorFrik/getman/collections"
	"github.com/KonnorFrik/getman/types"
)

// PostmanCollection represents a Postman collection structure.
type PostmanCollection struct {
	Info PostmanInfo   `json:"info"`
	Item []PostmanItem `json:"item"`
}

// PostmanInfo contains metadata about a Postman collection.
type PostmanInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
}

// PostmanItem represents an item in a Postman collection (can be a request or a folder).
type PostmanItem struct {
	Name     string          `json:"name"`
	Request  *PostmanRequest `json:"request,omitempty"`
	Item     []PostmanItem   `json:"item,omitempty"`
	Response []interface{}   `json:"response,omitempty"`
}

// PostmanRequest represents a request in a Postman collection.
type PostmanRequest struct {
	Method string          `json:"method"`
	Header []PostmanHeader `json:"header"`
	Body   *PostmanBody    `json:"body,omitempty"`
	URL    PostmanURL      `json:"url"`
	Auth   *PostmanAuth    `json:"auth,omitempty"`
}

// PostmanHeader represents a header in a Postman request.
type PostmanHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PostmanBody represents the body of a Postman request.
type PostmanBody struct {
	Mode       string              `json:"mode"`
	Raw        string              `json:"raw,omitempty"`
	Formdata   []PostmanFormData   `json:"formdata,omitempty"`
	Urlencoded []PostmanURLEncoded `json:"urlencoded,omitempty"`
}

// PostmanFormData represents a form data field in a Postman request.
type PostmanFormData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

// PostmanURLEncoded represents a URL-encoded form field in a Postman request.
type PostmanURLEncoded struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PostmanURL represents a URL structure in a Postman request.
type PostmanURL struct {
	Raw      string         `json:"raw"`
	Protocol string         `json:"protocol,omitempty"`
	Host     []string       `json:"host,omitempty"`
	Path     []string       `json:"path,omitempty"`
	Query    []PostmanQuery `json:"query,omitempty"`
}

// PostmanQuery represents a query parameter in a Postman URL.
type PostmanQuery struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// PostmanAuth represents authentication settings in a Postman request.
type PostmanAuth struct {
	Type   string             `json:"type"`
	Basic  []PostmanAuthField `json:"basic,omitempty"`
	Bearer []PostmanAuthField `json:"bearer,omitempty"`
	Apikey []PostmanAuthField `json:"apikey,omitempty"`
}

// PostmanAuthField represents a field in Postman authentication configuration.
type PostmanAuthField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

// ImportFromPostman imports a Postman collection from a JSON file and converts it to a Collection.
func ImportFromPostman(filePath string) (*collections.Collection, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Postman collection file: %w", err)
	}

	var postmanCollection PostmanCollection
	if err := json.Unmarshal(data, &postmanCollection); err != nil {
		return nil, fmt.Errorf("failed to parse Postman collection: %w", err)
	}

	collection := &collections.Collection{
		Name:        postmanCollection.Info.Name,
		Description: postmanCollection.Info.Description,
		Items:       []*types.RequestItem{},
	}

	for _, item := range postmanCollection.Item {
		items := convertPostmanItems(item)
		collection.Items = append(collection.Items, items...)
	}

	return collection, nil
}

func convertPostmanItems(item PostmanItem) []*types.RequestItem {
	var items []*types.RequestItem

	if item.Request != nil {
		req := convertPostmanRequest(item.Request)
		items = append(items, &types.RequestItem{
			Name:    item.Name,
			Request: req,
		})
	}

	for _, subItem := range item.Item {
		subItems := convertPostmanItems(subItem)
		items = append(items, subItems...)
	}

	return items
}

func convertPostmanRequest(postmanReq *PostmanRequest) *types.Request {
	req := &types.Request{
		Method:  strings.ToUpper(postmanReq.Method),
		URL:     postmanReq.URL.Raw,
		Headers: make(map[string]string),
	}

	for _, header := range postmanReq.Header {
		req.Headers[header.Key] = header.Value
	}

	if postmanReq.Body != nil {
		req.Body = convertPostmanBody(postmanReq.Body)
	}

	if postmanReq.Auth != nil {
		req.Auth = convertPostmanAuth(postmanReq.Auth)
	}

	return req
}

func convertPostmanBody(postmanBody *PostmanBody) *types.RequestBody {
	switch strings.ToLower(postmanBody.Mode) {
	case "raw":
		contentType := "text/plain"
		if strings.Contains(strings.ToLower(postmanBody.Raw), "json") {
			contentType = "application/json"
		} else if strings.Contains(strings.ToLower(postmanBody.Raw), "xml") {
			contentType = "application/xml"
		}
		return &types.RequestBody{
			Type:        "raw",
			Content:     []byte(postmanBody.Raw),
			ContentType: contentType,
		}
	case "formdata":
		var parts []string
		for _, field := range postmanBody.Formdata {
			parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(field.Key), url.QueryEscape(field.Value)))
		}
		body := strings.Join(parts, "&")
		return &types.RequestBody{
			Type:        "formdata",
			Content:     []byte(body),
			ContentType: "multipart/form-data",
		}
	case "urlencoded":
		var parts []string
		for _, field := range postmanBody.Urlencoded {
			parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(field.Key), url.QueryEscape(field.Value)))
		}
		body := strings.Join(parts, "&")
		return &types.RequestBody{
			Type:        "urlencoded",
			Content:     []byte(body),
			ContentType: "application/x-www-form-urlencoded",
		}
	default:
		return nil
	}
}

func convertPostmanAuth(postmanAuth *PostmanAuth) *types.Auth {
	authType := strings.ToLower(postmanAuth.Type)

	switch authType {
	case "basic":
		var username, password string
		for _, field := range postmanAuth.Basic {
			switch field.Key {
			case "username":
				username = field.Value
			case "password":
				password = field.Value
			}
		}
		return &types.Auth{
			Type:     "basic",
			Username: username,
			Password: password,
		}
	case "bearer":
		var token string
		for _, field := range postmanAuth.Bearer {
			if field.Key == "token" {
				token = field.Value
			}
		}
		return &types.Auth{
			Type:  "bearer",
			Token: token,
		}
	case "apikey":
		var keyName, keyValue, location string
		for _, field := range postmanAuth.Apikey {
			switch field.Key {
			case "key":
				keyName = field.Value
			case "value":
				keyValue = field.Value
			case "in":
				switch field.Value {
				case "header":
					location = "header"
				case "query":
					location = "query"
				}
			}
		}
		return &types.Auth{
			Type:     "apikey",
			APIKey:   keyValue,
			KeyName:  keyName,
			Location: location,
		}
	default:
		return nil
	}
}
