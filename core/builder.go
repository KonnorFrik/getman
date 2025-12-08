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

type RequestBuilder struct {
	method  string
	url     string
	headers map[string]string
	body    *types.RequestBody
	auth    *types.Auth
	timeout *types.Timeout
	cookies *types.CookieSettings
}

func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		headers: make(map[string]string),
	}
}

func (b *RequestBuilder) Method(method string) *RequestBuilder {
	b.method = method
	return b
}

func (b *RequestBuilder) URL(url string) *RequestBuilder {
	b.url = url
	return b
}

func (b *RequestBuilder) Header(key, value string) *RequestBuilder {
	if b.headers == nil {
		b.headers = make(map[string]string)
	}
	b.headers[key] = value
	return b
}

func (b *RequestBuilder) Headers(headers map[string]string) *RequestBuilder {
	if b.headers == nil {
		b.headers = make(map[string]string)
	}
	for k, v := range headers {
		b.headers[k] = v
	}
	return b
}

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

func (b *RequestBuilder) BodyRaw(data []byte, contentType string) *RequestBuilder {
	b.body = &types.RequestBody{
		Type:        "raw",
		Content:     data,
		ContentType: contentType,
	}
	return b
}

func (b *RequestBuilder) BodyBinary(data []byte, contentType string) *RequestBuilder {
	b.body = &types.RequestBody{
		Type:        "binary",
		Content:     data,
		ContentType: contentType,
	}
	return b
}

func (b *RequestBuilder) AuthBasic(username, password string) *RequestBuilder {
	b.auth = &types.Auth{
		Type:     "basic",
		Username: username,
		Password: password,
	}
	return b
}

func (b *RequestBuilder) AuthBearer(token string) *RequestBuilder {
	b.auth = &types.Auth{
		Type:  "bearer",
		Token: token,
	}
	return b
}

func (b *RequestBuilder) AuthAPIKey(keyName, keyValue, location string) *RequestBuilder {
	b.auth = &types.Auth{
		Type:     "apikey",
		APIKey:   keyValue,
		KeyName:  keyName,
		Location: location,
	}
	return b
}

func (b *RequestBuilder) Timeout(connect, read time.Duration) *RequestBuilder {
	b.timeout = &types.Timeout{
		Connect: connect,
		Read:    read,
	}
	return b
}

func (b *RequestBuilder) CookiesAutoManage(autoManage bool) *RequestBuilder {
	b.cookies = &types.CookieSettings{
		AutoManage: autoManage,
	}
	return b
}

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
