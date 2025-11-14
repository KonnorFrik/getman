package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/types"
)

func TestUnitNewHTTPClient(t *testing.T) {
	tests := []struct {
		name              string
		connectTimeout    time.Duration
		readTimeout       time.Duration
		autoManageCookies bool
	}{
		{
			name:              "with cookies",
			connectTimeout:    10 * time.Second,
			readTimeout:       30 * time.Second,
			autoManageCookies: true,
		},
		{
			name:              "without cookies",
			connectTimeout:    5 * time.Second,
			readTimeout:       15 * time.Second,
			autoManageCookies: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewHTTPClient(tt.connectTimeout, tt.readTimeout, tt.autoManageCookies)
			if client == nil {
				t.Fatal("expected client to be created")
			}
			if client.client == nil {
				t.Fatal("expected http client to be created")
			}
			if tt.autoManageCookies && client.client.Jar == nil {
				t.Fatal("expected cookie jar to be set when autoManageCookies is true")
			}
			if !tt.autoManageCookies && client.client.Jar != nil {
				t.Fatal("expected cookie jar to be nil when autoManageCookies is false")
			}
		})
	}
}

func TestUnitExecute_GET(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    server.URL,
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	if string(resp.Body) != "OK" {
		t.Errorf("expected body 'OK', got %s", string(resp.Body))
	}
}

func TestUnitExecute_POST(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Created"))
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodPost,
		URL:    server.URL,
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status code 201, got %d", resp.StatusCode)
	}
}

func TestUnitExecute_PUT(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Updated"))
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodPut,
		URL:    server.URL,
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}
}

func TestUnitExecute_DELETE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodDelete,
		URL:    server.URL,
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("expected status code 204, got %d", resp.StatusCode)
	}
}

func TestUnitExecute_WithHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Custom-Header") != "test-value" {
			t.Errorf("expected X-Custom-Header to be 'test-value', got %s", r.Header.Get("X-Custom-Header"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    server.URL,
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

func TestUnitExecute_WithBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := make([]byte, 1024)
		n, _ := r.Body.Read(body)
		if string(body[:n]) != "test body" {
			t.Errorf("expected body 'test body', got %s", string(body[:n]))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodPost,
		URL:    server.URL,
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

func TestUnitExecute_WithJSONBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type to be 'application/json', got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodPost,
		URL:    server.URL,
		Body: &types.RequestBody{
			Type:        "json",
			Content:     []byte(`{"key": "value"}`),
			ContentType: "application/json",
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

func TestUnitExecute_WithXMLBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/xml" {
			t.Errorf("expected Content-Type to be 'application/xml', got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodPost,
		URL:    server.URL,
		Body: &types.RequestBody{
			Type:        "xml",
			Content:     []byte("<root>test</root>"),
			ContentType: "application/xml",
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

func TestUnitExecute_WithBasicAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("expected basic auth to be present")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if username != "testuser" || password != "testpass" {
			t.Errorf("expected username 'testuser' and password 'testpass', got %s:%s", username, password)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    server.URL,
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

func TestUnitExecute_WithBearerAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer testtoken" {
			t.Errorf("expected Authorization header to be 'Bearer testtoken', got %s", authHeader)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    server.URL,
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

func TestUnitExecute_WithAPIKeyAuth_Header(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey != "testapikey" {
			t.Errorf("expected X-API-Key header to be 'testapikey', got %s", apiKey)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    server.URL,
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

func TestUnitExecute_WithAPIKeyAuth_Query(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.URL.Query().Get("api_key")
		if apiKey != "testapikey" {
			t.Errorf("expected api_key query parameter to be 'testapikey', got %s", apiKey)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    server.URL,
		Auth: &types.Auth{
			Type:     "apikey",
			APIKey:   "testapikey",
			KeyName:  "api_key",
			Location: "query",
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

func TestUnitExecute_WithCookies(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := &http.Cookie{
			Name:  "testcookie",
			Value: "testvalue",
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, true)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    server.URL,
	}

	resp, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp.StatusCode)
	}

	req2 := &types.Request{
		Method: http.MethodGet,
		URL:    server.URL,
	}

	resp2, err := client.Execute(req2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp2.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", resp2.StatusCode)
	}
}

func TestUnitExecute_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1000 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10 * time.Millisecond, 10 * time.Millisecond, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    server.URL,
	}

	_, err := client.Execute(req)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestUnitExecute_InvalidURL(t *testing.T) {
	client := NewHTTPClient(10*time.Second, 30*time.Second, false)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    "invalid-url",
	}

	_, err := client.Execute(req)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestUnitCookieJar(t *testing.T) {
	var cookieValue string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cookieValue == "" {
			cookie := &http.Cookie{
				Name:  "testcookie",
				Value: "testvalue",
			}
			http.SetCookie(w, cookie)
			cookieValue = "testvalue"
		} else {
			cookie, err := r.Cookie("testcookie")
			if err != nil {
				t.Errorf("expected cookie to be present, got error: %v", err)
			} else if cookie.Value != "testvalue" {
				t.Errorf("expected cookie value 'testvalue', got %s", cookie.Value)
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewHTTPClient(10*time.Second, 30*time.Second, true)
	req := &types.Request{
		Method: http.MethodGet,
		URL:    server.URL,
	}

	_, err := client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.Execute(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
