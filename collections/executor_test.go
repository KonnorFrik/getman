package collections

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/core"
	"github.com/KonnorFrik/getman/types"
	"github.com/KonnorFrik/getman/variables"
)

func TestUnitExecuteCollection_EmptyCollection(t *testing.T) {
	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name:  "Empty Collection",
		Items: []*types.RequestItem{},
	}

	result, err := executor.ExecuteCollection(collection, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.CollectionName != "Empty Collection" {
		t.Errorf("expected collection name 'Empty Collection', got %s", result.CollectionName)
	}

	if result.Statistics.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Statistics.Total)
	}
}

func TestUnitExecuteCollection_SingleRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Test Request",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
		},
	}

	result, err := executor.ExecuteCollection(collection, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Statistics.Total)
	}

	if result.Statistics.Success != 1 {
		t.Errorf("expected success 1, got %d", result.Statistics.Success)
	}

	if len(result.Requests) != 1 {
		t.Fatalf("expected 1 request, got %d", len(result.Requests))
	}

	if result.Requests[0].Response == nil {
		t.Fatal("expected response to be set")
	}

	if result.Requests[0].Response.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", result.Requests[0].Response.StatusCode)
	}
}

func TestUnitExecuteCollection_MultipleRequests(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Request 1",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
			{
				Name: "Request 2",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
		},
	}

	result, err := executor.ExecuteCollection(collection, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Statistics.Total)
	}

	if result.Statistics.Success != 2 {
		t.Errorf("expected success 2, got %d", result.Statistics.Success)
	}

	if len(result.Requests) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(result.Requests))
	}
}

func TestUnitExecuteCollection_WithVariables(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users" {
			t.Errorf("expected path /api/users, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	store.SetEnv("baseUrl", server.URL)
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Test Request",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    "{{baseUrl}}/api/users",
				},
			},
		},
	}

	result, err := executor.ExecuteCollection(collection, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Statistics.Total)
	}

	if result.Requests[0].Response == nil {
		t.Fatal("expected response to be set")
	}
}

func TestUnitExecuteCollection_RequestError(t *testing.T) {
	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Test Request",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    "http://invalid-url-that-does-not-exist:9999",
				},
			},
		},
	}

	result, err := executor.ExecuteCollection(collection, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Statistics.Total)
	}

	if result.Statistics.Failed != 1 {
		t.Errorf("expected failed 1, got %d", result.Statistics.Failed)
	}

	if result.Requests[0].Error == "" {
		t.Fatal("expected error to be set")
	}
}

func TestUnitExecuteCollection_Statistics(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Request 1",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
			{
				Name: "Request 2",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
		},
	}

	result, err := executor.ExecuteCollection(collection, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics == nil {
		t.Fatal("expected statistics to be set")
	}

	if result.Statistics.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Statistics.Total)
	}

	if result.Statistics.Success != 2 {
		t.Errorf("expected success 2, got %d", result.Statistics.Success)
	}

	if result.Statistics.Failed != 0 {
		t.Errorf("expected failed 0, got %d", result.Statistics.Failed)
	}

	if result.Statistics.AvgTime <= 0 {
		t.Errorf("expected avg time to be positive, got %v", result.Statistics.AvgTime)
	}

	if result.Statistics.MinTime <= 0 {
		t.Errorf("expected min time to be positive, got %v", result.Statistics.MinTime)
	}

	if result.Statistics.MaxTime <= 0 {
		t.Errorf("expected max time to be positive, got %v", result.Statistics.MaxTime)
	}
}

func TestUnitExecuteCollectionSelective_AllItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Request 1",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
			{
				Name: "Request 2",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
		},
	}

	result, err := executor.ExecuteCollectionSelective(collection, "test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Statistics.Total)
	}
}

func TestUnitExecuteCollectionSelective_SpecificItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Request 1",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
			{
				Name: "Request 2",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
		},
	}

	result, err := executor.ExecuteCollectionSelective(collection, "test", []string{"Request 1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Statistics.Total)
	}

	if result.Requests[0].Request.URL != server.URL {
		t.Errorf("expected URL %s, got %s", server.URL, result.Requests[0].Request.URL)
	}
}

func TestUnitExecuteCollectionSelective_NonExistentItems(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Request 1",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
		},
	}

	result, err := executor.ExecuteCollectionSelective(collection, "test", []string{"NonExistent"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Statistics.Total)
	}
}

func TestUnitResolveRequest_URL(t *testing.T) {
	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	store.SetEnv("baseUrl", "http://example.com")
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	req := &types.Request{
		Method: http.MethodGet,
		URL:    "{{baseUrl}}/api/users",
	}

	resolvedReq, err := executor.resolveRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolvedReq.URL != "http://example.com/api/users" {
		t.Errorf("expected URL 'http://example.com/api/users', got %s", resolvedReq.URL)
	}
}

func TestUnitResolveRequest_Headers(t *testing.T) {
	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	store.SetEnv("headerName", "Authorization")
	store.SetEnv("headerValue", "Bearer token123")
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	req := &types.Request{
		Method: http.MethodGet,
		URL:    "http://example.com",
		Headers: map[string]string{
			"{{headerName}}": "{{headerValue}}",
		},
	}

	resolvedReq, err := executor.resolveRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolvedReq.Headers["Authorization"] != "Bearer token123" {
		t.Errorf("expected header 'Authorization' to be 'Bearer token123', got %s", resolvedReq.Headers["Authorization"])
	}
}

func TestUnitResolveRequest_Body(t *testing.T) {
	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	store.SetEnv("userName", "Test User")
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	req := &types.Request{
		Method: http.MethodPost,
		URL:    "http://example.com",
		Body: &types.RequestBody{
			Type:        "json",
			Content:     []byte(`{"name": "{{userName}}"}`),
			ContentType: "application/json",
		},
	}

	resolvedReq, err := executor.resolveRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(resolvedReq.Body.Content) != `{"name": "Test User"}` {
		t.Errorf("expected body '{\"name\": \"Test User\"}', got %s", string(resolvedReq.Body.Content))
	}
}

func TestUnitResolveRequest_Auth(t *testing.T) {
	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	store.SetEnv("token", "testtoken123")
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	req := &types.Request{
		Method: http.MethodGet,
		URL:    "http://example.com",
		Auth: &types.Auth{
			Type:  "bearer",
			Token: "{{token}}",
		},
	}

	resolvedReq, err := executor.resolveRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolvedReq.Auth.Token != "testtoken123" {
		t.Errorf("expected token 'testtoken123', got %s", resolvedReq.Auth.Token)
	}
}

func TestUnitResolveRequest_VariableNotFound(t *testing.T) {
	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	req := &types.Request{
		Method: http.MethodGet,
		URL:    "{{nonexistent}}",
	}

	_, err := executor.resolveRequest(req)
	if err == nil {
		t.Fatal("expected error for nonexistent variable")
	}
}

func TestUnitExecuteCollection_StatusCodeFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error"))
	}))
	defer server.Close()

	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Test Request",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    server.URL,
				},
			},
		},
	}

	result, err := executor.ExecuteCollection(collection, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Success != 0 {
		t.Errorf("expected success 0, got %d", result.Statistics.Success)
	}

	if result.Statistics.Failed != 1 {
		t.Errorf("expected failed 1, got %d", result.Statistics.Failed)
	}
}

func TestUnitExecuteCollection_VariableResolutionError(t *testing.T) {
	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Test Request",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    "{{nonexistent}}",
				},
			},
		},
	}

	result, err := executor.ExecuteCollection(collection, "test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Failed != 1 {
		t.Errorf("expected failed 1, got %d", result.Statistics.Failed)
	}

	if result.Requests[0].Error == "" {
		t.Fatal("expected error to be set")
	}
}

func TestUnitResolveRequest_BasicAuth(t *testing.T) {
	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	store.SetEnv("username", "testuser")
	store.SetEnv("password", "testpass")
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	req := &types.Request{
		Method: http.MethodGet,
		URL:    "http://example.com",
		Auth: &types.Auth{
			Type:     "basic",
			Username: "{{username}}",
			Password: "{{password}}",
		},
	}

	resolvedReq, err := executor.resolveRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolvedReq.Auth.Username != "testuser" {
		t.Errorf("expected username 'testuser', got %s", resolvedReq.Auth.Username)
	}

	if resolvedReq.Auth.Password != "testpass" {
		t.Errorf("expected password 'testpass', got %s", resolvedReq.Auth.Password)
	}
}

func TestUnitResolveRequest_APIKeyAuth(t *testing.T) {
	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	store.SetEnv("apikey", "testapikey123")
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	req := &types.Request{
		Method: http.MethodGet,
		URL:    "http://example.com",
		Auth: &types.Auth{
			Type:     "apikey",
			APIKey:   "{{apikey}}",
			KeyName:  "X-API-Key",
			Location: "header",
		},
	}

	resolvedReq, err := executor.resolveRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resolvedReq.Auth.APIKey != "testapikey123" {
		t.Errorf("expected API key 'testapikey123', got %s", resolvedReq.Auth.APIKey)
	}
}
