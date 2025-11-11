package importer

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/KonnorFrik/getman/types"
)

type PostmanCollection struct {
	Info PostmanInfo   `json:"info"`
	Item []PostmanItem `json:"item"`
}

type PostmanInfo struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Schema      string `json:"schema"`
}

type PostmanItem struct {
	Name     string          `json:"name"`
	Request  *PostmanRequest `json:"request,omitempty"`
	Item     []PostmanItem   `json:"item,omitempty"`
	Response []interface{}   `json:"response,omitempty"`
}

type PostmanRequest struct {
	Method string          `json:"method"`
	Header []PostmanHeader `json:"header"`
	Body   *PostmanBody    `json:"body,omitempty"`
	URL    PostmanURL      `json:"url"`
	Auth   *PostmanAuth    `json:"auth,omitempty"`
}

type PostmanHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PostmanBody struct {
	Mode       string              `json:"mode"`
	Raw        string              `json:"raw,omitempty"`
	Formdata   []PostmanFormData   `json:"formdata,omitempty"`
	Urlencoded []PostmanURLEncoded `json:"urlencoded,omitempty"`
}

type PostmanFormData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

type PostmanURLEncoded struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PostmanURL struct {
	Raw      string         `json:"raw"`
	Protocol string         `json:"protocol,omitempty"`
	Host     []string       `json:"host,omitempty"`
	Path     []string       `json:"path,omitempty"`
	Query    []PostmanQuery `json:"query,omitempty"`
}

type PostmanQuery struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PostmanAuth struct {
	Type   string             `json:"type"`
	Basic  []PostmanAuthField `json:"basic,omitempty"`
	Bearer []PostmanAuthField `json:"bearer,omitempty"`
	Apikey []PostmanAuthField `json:"apikey,omitempty"`
}

type PostmanAuthField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

func ImportFromPostman(filePath string) (*types.Collection, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Postman collection file: %w", err)
	}

	var postmanCollection PostmanCollection
	if err := json.Unmarshal(data, &postmanCollection); err != nil {
		return nil, fmt.Errorf("failed to parse Postman collection: %w", err)
	}

	collection := &types.Collection{
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
