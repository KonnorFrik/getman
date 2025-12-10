/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package core

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/KonnorFrik/getman/types"
)

// RequestBuilder provides a fluent interface for building HTTP requests.
type RequestBuilder struct {
	method  string
	url     string
	headers map[string]string
	body    *types.RequestBody
	auth    *types.Auth
	timeout *types.Timeout
	cookies *types.CookieSettings
}

// NewRequestBuilder creates a new RequestBuilder instance.
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		headers: make(map[string]string),
	}
}

// Method sets the HTTP method for the request.
func (b *RequestBuilder) Method(method string) *RequestBuilder {
	b.method = method
	return b
}

// URL sets the request URL.
func (b *RequestBuilder) URL(url string) *RequestBuilder {
	b.url = url
	return b
}

// Header adds a single header to the request.
func (b *RequestBuilder) Header(key, value string) *RequestBuilder {
	if b.headers == nil {
		b.headers = make(map[string]string)
	}
	b.headers[key] = value
	return b
}

// Headers sets multiple headers for the request.
func (b *RequestBuilder) Headers(headers map[string]string) *RequestBuilder {
	if b.headers == nil {
		b.headers = make(map[string]string)
	}
	for k, v := range headers {
		b.headers[k] = v
	}
	return b
}

// BodyJSON sets the request body as JSON.
func (b *RequestBuilder) BodyJSON(data interface{}) *RequestBuilder {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return b
	}

	b.body = &types.RequestBody{
		Type:        "json",
		Content:     jsonData,
		ContentType: "application/json",
	}
	return b
}

// BodyXML sets the request body as XML.
func (b *RequestBuilder) BodyXML(data string) *RequestBuilder {
	xmlData, err := xml.Marshal(data)
	if err != nil {
		return b
	}

	b.body = &types.RequestBody{
		Type:        "xml",
		Content:     xmlData,
		ContentType: "application/xml",
	}
	return b
}

// BodyRaw sets the request body as raw bytes with the specified content type.
func (b *RequestBuilder) BodyRaw(data []byte, contentType string) *RequestBuilder {
	b.body = &types.RequestBody{
		Type:        "raw",
		Content:     data,
		ContentType: contentType,
	}
	return b
}

// BodyBinary sets the request body as binary data with the specified content type.
func (b *RequestBuilder) BodyBinary(data []byte, contentType string) *RequestBuilder {
	b.body = &types.RequestBody{
		Type:        "binary",
		Content:     data,
		ContentType: contentType,
	}
	return b
}

// AuthBasic sets Basic authentication credentials.
func (b *RequestBuilder) AuthBasic(username, password string) *RequestBuilder {
	b.auth = &types.Auth{
		Type:     "basic",
		Username: username,
		Password: password,
	}
	return b
}

// AuthBearer sets Bearer token authentication.
func (b *RequestBuilder) AuthBearer(token string) *RequestBuilder {
	b.auth = &types.Auth{
		Type:  "bearer",
		Token: token,
	}
	return b
}

// AuthAPIKey sets API key authentication with the specified key name, value, and location.
func (b *RequestBuilder) AuthAPIKey(keyName, keyValue, location string) *RequestBuilder {
	b.auth = &types.Auth{
		Type:     "apikey",
		APIKey:   keyValue,
		KeyName:  keyName,
		Location: location,
	}
	return b
}

// Timeout sets connection and read timeouts for the request.
func (b *RequestBuilder) Timeout(connect, read time.Duration) *RequestBuilder {
	b.timeout = &types.Timeout{
		Connect: connect,
		Read:    read,
	}
	return b
}

// CookiesAutoManage enables or disables automatic cookie management.
func (b *RequestBuilder) CookiesAutoManage(autoManage bool) *RequestBuilder {
	b.cookies = &types.CookieSettings{
		AutoManage: autoManage,
	}
	return b
}

// Build constructs and returns the final Request object.
func (b *RequestBuilder) Build() (*types.Request, error) {
	if b.method == "" {
		return nil, fmt.Errorf("method is required")
	}

	if b.url == "" {
		return nil, fmt.Errorf("url is required")
	}

	req := &types.Request{
		Method:  b.method,
		URL:     b.url,
		Headers: b.headers,
		Body:    b.body,
		Auth:    b.auth,
		Timeout: b.timeout,
		Cookies: b.cookies,
	}

	return req, nil
}
