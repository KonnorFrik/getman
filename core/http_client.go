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
	"time"

	"github.com/KonnorFrik/getman/errors"
	"github.com/KonnorFrik/getman/types"
)

type HTTPClient struct {
	client     *http.Client
	autoManage bool
}

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
	cookies map[string][]*http.Cookie
}

func (j *cookieJarImpl) SetCookies(u *url.URL, cookies []*http.Cookie) {
	j.cookies[u.Host] = cookies
}

func (j *cookieJarImpl) Cookies(u *url.URL) []*http.Cookie {
	return j.cookies[u.Host]
}

func (hc *HTTPClient) Execute(req *types.Request) (*types.Response, error) {
	startTime := time.Now()

	httpReq, err := hc.buildHTTPRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP request: %w", err)
	}

	httpResp, err := hc.client.Do(httpReq)
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
	case "basic":
		httpReq.SetBasicAuth(auth.Username, auth.Password)
	case "bearer":
		httpReq.Header.Set("Authorization", "Bearer "+auth.Token)
	case "apikey":
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
