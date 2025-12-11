package core

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnitRequestBuilder_Method(t *testing.T) {
	builder := NewRequestBuilder()
	builder.Method("GET")

	req, err := builder.URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("expected method 'GET', got %s", req.Method)
	}
}

func TestUnitRequestBuilder_URL(t *testing.T) {
	builder := NewRequestBuilder()
	builder.URL("http://example.com")

	req, err := builder.Method("GET").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.URL != "http://example.com" {
		t.Errorf("expected URL 'http://example.com', got %s", req.URL)
	}
}

func TestUnitRequestBuilder_Header(t *testing.T) {
	builder := NewRequestBuilder()
	builder.Header("Content-Type", "application/json")

	req, err := builder.Method("GET").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Headers["Content-Type"] != "application/json" {
		t.Errorf("expected header 'Content-Type' to be 'application/json', got %s", req.Headers["Content-Type"])
	}
}

func TestUnitRequestBuilder_Headers(t *testing.T) {
	builder := NewRequestBuilder()
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	builder.Headers(headers)

	req, err := builder.Method("GET").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Headers["Content-Type"] != "application/json" {
		t.Errorf("expected header 'Content-Type' to be 'application/json', got %s", req.Headers["Content-Type"])
	}

	if req.Headers["Accept"] != "application/json" {
		t.Errorf("expected header 'Accept' to be 'application/json', got %s", req.Headers["Accept"])
	}
}

func TestUnitRequestBuilder_BodyJSON(t *testing.T) {
	builder := NewRequestBuilder()
	data := map[string]string{
		"key": "value",
	}
	builder.BodyJSON(data)

	req, err := builder.Method("POST").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Body == nil {
		t.Fatal("expected body to be set")
	}

	if req.Body.Type != "json" {
		t.Errorf("expected body type 'json', got %s", req.Body.Type)
	}

	if req.Body.ContentType != "application/json" {
		t.Errorf("expected content type 'application/json', got %s", req.Body.ContentType)
	}

	var result map[string]string
	if err := json.Unmarshal(req.Body.Content, &result); err != nil {
		t.Fatalf("unexpected error unmarshaling body: %v", err)
	}

	if result["key"] != "value" {
		t.Errorf("expected body key 'value', got %s", result["key"])
	}
}

func TestUnitRequestBuilder_BodyXML(t *testing.T) {
	builder := NewRequestBuilder()
	builder.BodyXML("test data")

	req, err := builder.Method("POST").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Body == nil {
		t.Fatal("expected body to be set")
	}

	if req.Body.Type != "xml" {
		t.Errorf("expected body type 'xml', got %s", req.Body.Type)
	}

	if req.Body.ContentType != "application/xml" {
		t.Errorf("expected content type 'application/xml', got %s", req.Body.ContentType)
	}
}

func TestUnitRequestBuilder_BodyBinary(t *testing.T) {
	builder := NewRequestBuilder()
	builder.BodyBinary([]byte{1, 2, 3, 4}, "application/octet-stream")

	req, err := builder.Method("POST").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Body == nil {
		t.Fatal("expected body to be set")
	}

	if req.Body.Type != "binary" {
		t.Errorf("expected body type 'binary', got %s", req.Body.Type)
	}

	if req.Body.ContentType != "application/octet-stream" {
		t.Errorf("expected content type 'application/octet-stream', got %s", req.Body.ContentType)
	}
}

func TestUnitRequestBuilder_AuthBasic(t *testing.T) {
	builder := NewRequestBuilder()
	builder.AuthBasic("username", "password")

	req, err := builder.Method("GET").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Auth == nil {
		t.Fatal("expected auth to be set")
	}

	if req.Auth.Type != "basic" {
		t.Errorf("expected auth type 'basic', got %s", req.Auth.Type)
	}

	if req.Auth.Username != "username" {
		t.Errorf("expected username 'username', got %s", req.Auth.Username)
	}

	if req.Auth.Password != "password" {
		t.Errorf("expected password 'password', got %s", req.Auth.Password)
	}
}

func TestUnitRequestBuilder_AuthBearer(t *testing.T) {
	builder := NewRequestBuilder()
	builder.AuthBearer("token123")

	req, err := builder.Method("GET").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Auth == nil {
		t.Fatal("expected auth to be set")
	}

	if req.Auth.Type != "bearer" {
		t.Errorf("expected auth type 'bearer', got %s", req.Auth.Type)
	}

	if req.Auth.Token != "token123" {
		t.Errorf("expected token 'token123', got %s", req.Auth.Token)
	}
}

func TestUnitRequestBuilder_AuthAPIKey(t *testing.T) {
	builder := NewRequestBuilder()
	builder.AuthAPIKey("X-API-Key", "apikey123", "header")

	req, err := builder.Method("GET").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Auth == nil {
		t.Fatal("expected auth to be set")
	}

	if req.Auth.Type != "apikey" {
		t.Errorf("expected auth type 'apikey', got %s", req.Auth.Type)
	}

	if req.Auth.APIKey != "apikey123" {
		t.Errorf("expected API key 'apikey123', got %s", req.Auth.APIKey)
	}

	if req.Auth.KeyName != "X-API-Key" {
		t.Errorf("expected key name 'X-API-Key', got %s", req.Auth.KeyName)
	}

	if req.Auth.Location != "header" {
		t.Errorf("expected location 'header', got %s", req.Auth.Location)
	}
}

func TestUnitRequestBuilder_Timeout(t *testing.T) {
	builder := NewRequestBuilder()
	builder.Timeout(10*time.Second, 30*time.Second)

	req, err := builder.Method("GET").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Timeout == nil {
		t.Fatal("expected timeout to be set")
	}

	if req.Timeout.Connect != 10*time.Second {
		t.Errorf("expected connect timeout 10s, got %v", req.Timeout.Connect)
	}

	if req.Timeout.Read != 30*time.Second {
		t.Errorf("expected read timeout 30s, got %v", req.Timeout.Read)
	}
}

func TestUnitRequestBuilder_CookiesAutoManage(t *testing.T) {
	builder := NewRequestBuilder()
	builder.CookiesAutoManage(true)

	req, err := builder.Method("GET").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Cookies == nil {
		t.Fatal("expected cookies to be set")
	}

	if req.Cookies.AutoManage != true {
		t.Errorf("expected auto manage to be true, got %v", req.Cookies.AutoManage)
	}
}

func TestUnitRequestBuilder_Build_Valid(t *testing.T) {
	builder := NewRequestBuilder()
	req, err := builder.
		Method("GET").
		URL("http://example.com").
		Header("Accept", "application/json").
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("expected method 'GET', got %s", req.Method)
	}

	if req.URL != "http://example.com" {
		t.Errorf("expected URL 'http://example.com', got %s", req.URL)
	}

	if req.Headers["Accept"] != "application/json" {
		t.Errorf("expected header 'Accept' to be 'application/json', got %s", req.Headers["Accept"])
	}
}

func TestUnitRequestBuilder_Build_MissingMethod(t *testing.T) {
	builder := NewRequestBuilder()
	_, err := builder.URL("http://example.com").Build()

	if err == nil {
		t.Fatal("expected error for missing method")
	}
}

func TestUnitRequestBuilder_Build_MissingURL(t *testing.T) {
	builder := NewRequestBuilder()
	_, err := builder.Method("GET").Build()

	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestUnitRequestBuilder_Chaining(t *testing.T) {
	req, err := NewRequestBuilder().
		Method("POST").
		URL("http://example.com").
		Header("Content-Type", "application/json").
		BodyJSON(map[string]string{"key": "value"}).
		AuthBearer("token123").
		Timeout(10*time.Second, 30*time.Second).
		Build()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Method != "POST" {
		t.Errorf("expected method 'POST', got %s", req.Method)
	}

	if req.URL != "http://example.com" {
		t.Errorf("expected URL 'http://example.com', got %s", req.URL)
	}

	if req.Headers["Content-Type"] != "application/json" {
		t.Errorf("expected header 'Content-Type' to be 'application/json', got %s", req.Headers["Content-Type"])
	}

	if req.Body == nil {
		t.Fatal("expected body to be set")
	}

	if req.Auth == nil {
		t.Fatal("expected auth to be set")
	}

	if req.Timeout == nil {
		t.Fatal("expected timeout to be set")
	}
}

func TestUnitRequestBuilder_MultipleHeaders(t *testing.T) {
	builder := NewRequestBuilder()
	builder.Header("Header1", "Value1")
	builder.Header("Header2", "Value2")

	req, err := builder.Method("GET").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(req.Headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(req.Headers))
	}

	if req.Headers["Header1"] != "Value1" {
		t.Errorf("expected header 'Header1' to be 'Value1', got %s", req.Headers["Header1"])
	}

	if req.Headers["Header2"] != "Value2" {
		t.Errorf("expected header 'Header2' to be 'Value2', got %s", req.Headers["Header2"])
	}
}

func TestUnitRequestBuilder_OverrideHeader(t *testing.T) {
	builder := NewRequestBuilder()
	builder.Header("Content-Type", "application/json")
	builder.Header("Content-Type", "application/xml")

	req, err := builder.Method("POST").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Headers["Content-Type"] != "application/xml" {
		t.Errorf("expected header 'Content-Type' to be 'application/xml', got %s", req.Headers["Content-Type"])
	}
}

func TestUnitRequestBuilder_BodyJSON_InvalidData(t *testing.T) {
	builder := NewRequestBuilder()
	invalidData := func() {}
	builder.BodyJSON(invalidData)

	req, err := builder.Method("POST").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Body != nil {
		t.Error("expected body to be nil for invalid JSON data")
	}
}

func TestUnitRequestBuilder_BodyXML_ValidData(t *testing.T) {
	builder := NewRequestBuilder()
	builder.BodyXML("<root>test</root>")

	req, err := builder.Method("POST").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Body == nil {
		t.Fatal("expected body to be set")
	}

	if req.Body.Type != "xml" {
		t.Errorf("expected body type 'xml', got %s", req.Body.Type)
	}
}

func TestUnitRequestBuilder_AuthAPIKey_QueryLocation(t *testing.T) {
	builder := NewRequestBuilder()
	builder.AuthAPIKey("api_key", "apikey123", "query")

	req, err := builder.Method("GET").URL("http://example.com").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if req.Auth == nil {
		t.Fatal("expected auth to be set")
	}

	if req.Auth.Location != "query" {
		t.Errorf("expected location 'query', got %s", req.Auth.Location)
	}
}
