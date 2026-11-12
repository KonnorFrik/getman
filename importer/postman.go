/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package importer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
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

// ExportToPostman exports a collection to a Postman Collection v2.1 JSON file.
func ExportToPostman(collection *collections.Collection, filePath string) error {
	if collection == nil {
		return fmt.Errorf("collection is nil")
	}

	postmanCollection := PostmanCollection{
		Info: PostmanInfo{
			Name:        collection.Name,
			Description: collection.Description,
			Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item: make([]PostmanItem, 0, len(collection.Items)),
	}

	for _, item := range collection.Items {
		if item == nil || item.Request == nil {
			continue
		}

		postmanCollection.Item = append(postmanCollection.Item, PostmanItem{
			Name:    item.Name,
			Request: convertRequestToPostman(item.Request),
		})
	}

	data, err := json.MarshalIndent(postmanCollection, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal Postman collection: %w", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write Postman collection file: %w", err)
	}

	return nil
}

func convertRequestToPostman(req *types.Request) *PostmanRequest {
	postmanReq := &PostmanRequest{
		Method: req.Method,
		URL: PostmanURL{
			Raw: req.URL,
		},
		Header: make([]PostmanHeader, 0, len(req.Headers)),
	}

	headerKeys := make([]string, 0, len(req.Headers))
	for key := range req.Headers {
		headerKeys = append(headerKeys, key)
	}
	sort.Strings(headerKeys)

	for _, key := range headerKeys {
		postmanReq.Header = append(postmanReq.Header, PostmanHeader{
			Key:   key,
			Value: req.Headers[key],
		})
	}

	postmanReq.Body = convertBodyToPostman(req.Body)
	postmanReq.Auth = convertAuthToPostman(req.Auth)

	if req.Body != nil && req.Body.ContentType != "" && !hasHeader(req.Headers, "Content-Type") {
		postmanReq.Header = append(postmanReq.Header, PostmanHeader{
			Key:   "Content-Type",
			Value: req.Body.ContentType,
		})
	}

	return postmanReq
}

func hasHeader(headers map[string]string, name string) bool {
	for key := range headers {
		if strings.EqualFold(key, name) {
			return true
		}
	}

	return false
}

func convertBodyToPostman(body *types.RequestBody) *PostmanBody {
	if body == nil {
		return nil
	}

	switch strings.ToLower(body.Type) {
	case "formdata":
		return &PostmanBody{
			Mode:     "formdata",
			Formdata: queryToFormData(string(body.Content)),
		}
	case "urlencoded":
		return &PostmanBody{
			Mode:       "urlencoded",
			Urlencoded: queryToURLEncoded(string(body.Content)),
		}
	default:
		return &PostmanBody{
			Mode: "raw",
			Raw:  string(body.Content),
		}
	}
}

func queryToFormData(raw string) []PostmanFormData {
	values, err := url.ParseQuery(raw)
	if err != nil {
		return nil
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([]PostmanFormData, 0, len(keys))
	for _, key := range keys {
		for _, value := range values[key] {
			result = append(result, PostmanFormData{Key: key, Value: value})
		}
	}

	return result
}

func queryToURLEncoded(raw string) []PostmanURLEncoded {
	values, err := url.ParseQuery(raw)
	if err != nil {
		return nil
	}

	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	result := make([]PostmanURLEncoded, 0, len(keys))
	for _, key := range keys {
		for _, value := range values[key] {
			result = append(result, PostmanURLEncoded{Key: key, Value: value})
		}
	}

	return result
}

func convertAuthToPostman(auth *types.Auth) *PostmanAuth {
	if auth == nil {
		return nil
	}

	switch strings.ToLower(auth.Type) {
	case "basic":
		return &PostmanAuth{
			Type: "basic",
			Basic: []PostmanAuthField{
				{Key: "username", Value: auth.Username, Type: "string"},
				{Key: "password", Value: auth.Password, Type: "string"},
			},
		}
	case "bearer":
		return &PostmanAuth{
			Type: "bearer",
			Bearer: []PostmanAuthField{
				{Key: "token", Value: auth.Token, Type: "string"},
			},
		}
	case "apikey":
		location := auth.Location
		if location == "" {
			location = "header"
		}
		return &PostmanAuth{
			Type: "apikey",
			Apikey: []PostmanAuthField{
				{Key: "key", Value: auth.KeyName, Type: "string"},
				{Key: "value", Value: auth.APIKey, Type: "string"},
				{Key: "in", Value: location, Type: "string"},
			},
		}
	default:
		return nil
	}
}

func convertPostmanRequest(postmanReq *PostmanRequest) *types.Request {
	req := &types.Request{
		Method:  strings.ToUpper(postmanReq.Method),
		URL:     postmanURLString(postmanReq.URL),
		Headers: make(map[string]string),
	}

	contentType := ""

	for _, header := range postmanReq.Header {
		req.Headers[header.Key] = header.Value

		if strings.EqualFold(header.Key, "Content-Type") {
			contentType = header.Value
		}
	}

	if postmanReq.Body != nil {
		req.Body = convertPostmanBody(postmanReq.Body, contentType)
	}

	if postmanReq.Auth != nil {
		req.Auth = convertPostmanAuth(postmanReq.Auth)
	}

	return req
}

func postmanURLString(postmanURL PostmanURL) string {
	if postmanURL.Raw != "" {
		return postmanURL.Raw
	}

	raw := strings.Join(postmanURL.Host, ".")

	if postmanURL.Protocol != "" {
		raw = postmanURL.Protocol + "://" + raw
	}

	if len(postmanURL.Path) > 0 {
		raw += "/" + strings.Join(postmanURL.Path, "/")
	}

	if len(postmanURL.Query) > 0 {
		query := make([]string, 0, len(postmanURL.Query))
		for _, q := range postmanURL.Query {
			query = append(query, url.QueryEscape(q.Key)+"="+url.QueryEscape(q.Value))
		}
		raw += "?" + strings.Join(query, "&")
	}

	return raw
}

func convertPostmanBody(postmanBody *PostmanBody, contentType string) *types.RequestBody {
	switch strings.ToLower(postmanBody.Mode) {
	case "raw":
		if contentType == "" {
			contentType = inferContentType(postmanBody.Raw)
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

func inferContentType(raw string) string {
	trimmed := strings.TrimSpace(raw)

	switch {
	case strings.HasPrefix(trimmed, "{"), strings.HasPrefix(trimmed, "["):
		return "application/json"
	case strings.HasPrefix(trimmed, "<"):
		return "application/xml"
	default:
		return "text/plain"
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
