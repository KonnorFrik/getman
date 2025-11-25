package core

import (
	"net/http"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/testutil/http_server"
	"github.com/KonnorFrik/getman/types"
)

func TestIntegrationExecute_GET(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    http_server.GetServerURL() + "/health",
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestIntegrationExecute_POST(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodPost,
		URL:    http_server.GetServerURL() + "/echo",
		Body: &types.RequestBody{
			Type:        "raw",
			Content:     []byte("test body"),
			ContentType: "text/plain",
		},
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestIntegrationExecute_WithHeaders(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    http_server.GetServerURL() + "/headers",
		Headers: map[string]string{
			"X-Custom-Header": "test-value",
		},
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestIntegrationExecute_WithBody(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodPost,
		URL:    http_server.GetServerURL() + "/body",
		Body: &types.RequestBody{
			Type:        "raw",
			Content:     []byte("test body content"),
			ContentType: "text/plain",
		},
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestIntegrationExecute_WithBasicAuth(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    http_server.GetServerURL() + "/auth/basic",
		Auth: &types.Auth{
			Type:     "basic",
			Username: "testuser",
			Password: "testpass",
		},
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestIntegrationExecute_WithBearerAuth(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    http_server.GetServerURL() + "/auth/bearer",
		Auth: &types.Auth{
			Type:  "bearer",
			Token: "testtoken",
		},
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestIntegrationExecute_WithAPIKeyAuth(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    http_server.GetServerURL() + "/auth/apikey",
		Auth: &types.Auth{
			Type:     "apikey",
			APIKey:   "testapikey",
			KeyName:  "X-API-Key",
			Location: "header",
		},
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestIntegrationExecute_WithCookies(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, true)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    http_server.GetServerURL() + "/cookies",
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestIntegrationExecute_StatusCode(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)

	tests := []struct {
		statusCode int
		url        string
	}{
		{200, "/status/200"},
		{404, "/status/404"},
		{500, "/status/500"},
	}

	for _, tt := range tests {
		req := &types.Request{
			Method: http.MethodGet,
			URL:    http_server.GetServerURL() + tt.url,
		}

		resp, err := client.Execute(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.StatusCode != tt.statusCode {
			t.Errorf("expected status code %d, got %d", tt.statusCode, resp.StatusCode)
		}
	}
}

func TestIntegrationExecute_ResponseHeaders(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    http_server.GetServerURL() + "/headers",
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Headers) == 0 {
		t.Error("expected response to have headers")
	}
}

func TestIntegrationExecute_ResponseBody(t *testing.T) {
	_, err := http_server.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer http_server.StopTestServer()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    http_server.GetServerURL() + "/health",
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resp.Body) == 0 {
		t.Error("expected response to have body")
	}
}
