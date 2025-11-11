package collections

import (
	"net/http"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/core"
	"github.com/KonnorFrik/getman/testutil"
	"github.com/KonnorFrik/getman/types"
	"github.com/KonnorFrik/getman/variables"
)

func TestIntegrationExecuteCollection_SingleRequest(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

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
					URL:    testutil.GetServerURL() + "/health",
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
}

func TestIntegrationExecuteCollection_MultipleRequests(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

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
					URL:    testutil.GetServerURL() + "/health",
				},
			},
			{
				Name: "Request 2",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    testutil.GetServerURL() + "/health",
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
}

func TestIntegrationExecuteCollection_WithVariables(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

	httpClient := core.NewHTTPClient(10*time.Second, 30*time.Second, false)
	store := variables.NewVariableStore()
	store.SetEnv("baseUrl", testutil.GetServerURL())
	resolver := core.NewVariableResolver(store)
	executor := NewCollectionExecutor(httpClient, resolver)

	collection := &types.Collection{
		Name: "Test Collection",
		Items: []*types.RequestItem{
			{
				Name: "Test Request",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    "{{baseUrl}}/health",
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
}

func TestIntegrationExecuteCollection_WithAuth(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

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
					URL:    testutil.GetServerURL() + "/auth/basic",
					Auth: &types.Auth{
						Type:     "basic",
						Username: "testuser",
						Password: "testpass",
					},
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
}

func TestIntegrationExecuteCollection_Statistics(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

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
					URL:    testutil.GetServerURL() + "/health",
				},
			},
			{
				Name: "Request 2",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    testutil.GetServerURL() + "/health",
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

	if result.Statistics.AvgTime <= 0 {
		t.Error("expected avg time to be positive")
	}
}

func TestIntegrationExecuteCollectionSelective(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

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
					URL:    testutil.GetServerURL() + "/health",
				},
			},
			{
				Name: "Request 2",
				Request: &types.Request{
					Method: http.MethodGet,
					URL:    testutil.GetServerURL() + "/health",
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
}

