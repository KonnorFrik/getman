/*
Copyright © 2025 Шелковский Сергей (Shelkovskiy Sergey) <konnor.frik666@gmail.com>
*/
package core

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/KonnorFrik/getman/errors"
	"github.com/KonnorFrik/getman/types"
)

// HTTPClient provides functionality for executing HTTP requests.
type HTTPClient struct {
	client     *http.Client
	autoManage bool
}

// NewHTTPClient creates a new HTTPClient with the specified timeouts and cookie management settings.
func NewHTTPClient(connectTimeout, readTimeout time.Duration, autoManageCookies bool) *HTTPClient {
	transport := &http.Transport{
		ResponseHeaderTimeout: connectTimeout,
	}

	var jar http.CookieJar
	if autoManageCookies {
		jar = &cookieJarImpl{
			cookies: make(map[string][]*http.Cookie),
		}
	}

	client := &http.Client{
		Transport:     transport,
		Jar:           jar,
		Timeout:       readTimeout,
		CheckRedirect: nil,
	}

	return &HTTPClient{
		client:     client,
		autoManage: autoManageCookies,
	}
}

type cookieJarImpl struct {
	mu      sync.RWMutex
	cookies map[string][]*http.Cookie
}

func (j *cookieJarImpl) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.mu.Lock()
	defer j.mu.Unlock()

	if len(cookies) == 0 {
		delete(j.cookies, u.Host)
		return
	}

	j.cookies[u.Host] = append([]*http.Cookie(nil), cookies...)
}

func (j *cookieJarImpl) Cookies(u *url.URL) []*http.Cookie {
	j.mu.RLock()
	defer j.mu.RUnlock()
	return append([]*http.Cookie(nil), j.cookies[u.Host]...)
}

// Execute performs an HTTP request and returns the response.
func (hc *HTTPClient) Execute(req *types.Request) (*types.Response, error) {
	startTime := time.Now()

	httpReq, err := hc.buildHTTPRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP request: %w", err)
	}

	httpResp, err := hc.clientForRequest(req).Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrRequestFailed, err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	duration := time.Since(startTime)

	headers := make(map[string][]string)
	for k, v := range httpResp.Header {
		headers[k] = v
	}

	response := &types.Response{
		StatusCode: httpResp.StatusCode,
		Status:     httpResp.Status,
		Headers:    headers,
		Body:       body,
		Duration:   duration,
		Size:       int64(len(body)),
	}

	return response, nil
}

func (hc *HTTPClient) clientForRequest(req *types.Request) *http.Client {
	client := *hc.client

	if req.Timeout != nil {
		if req.Timeout.Read > 0 {
			client.Timeout = req.Timeout.Read
		}

		if req.Timeout.Connect > 0 {
			if transport, ok := hc.client.Transport.(*http.Transport); ok {
				cloned := transport.Clone()
				cloned.ResponseHeaderTimeout = req.Timeout.Connect
				client.Transport = cloned
			}
		}
	}

	if req.Cookies != nil {
		if req.Cookies.AutoManage {
			if client.Jar == nil {
				client.Jar = &cookieJarImpl{cookies: make(map[string][]*http.Cookie)}
			}
		} else {
			client.Jar = nil
		}
	}

	return &client
}

func (hc *HTTPClient) buildHTTPRequest(req *types.Request) (*http.Request, error) {
	var bodyReader io.Reader
	if req.Body != nil && len(req.Body.Content) > 0 {
		bodyReader = bytes.NewReader(req.Body.Content)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errors.ErrInvalidURL, err)
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	if req.Body != nil && req.Body.ContentType != "" {
		httpReq.Header.Set("Content-Type", req.Body.ContentType)
	}

	if req.Auth != nil {
		hc.applyAuth(httpReq, req.Auth)
	}

	return httpReq, nil
}

func (hc *HTTPClient) applyAuth(httpReq *http.Request, auth *types.Auth) {
	switch strings.ToLower(auth.Type) {
	case authTypeBasic:
		httpReq.SetBasicAuth(auth.Username, auth.Password)
	case authTypeBearer:
		httpReq.Header.Set("Authorization", "Bearer "+auth.Token)
	case authTypeApiKey:
		if strings.ToLower(auth.Location) == "header" {
			httpReq.Header.Set(auth.KeyName, auth.APIKey)
		} else if strings.ToLower(auth.Location) == "query" {
			u, err := url.Parse(httpReq.URL.String())
			if err == nil {
				q := u.Query()
				q.Set(auth.KeyName, auth.APIKey)
				u.RawQuery = q.Encode()
				httpReq.URL = u
			}
		}
	}
}
