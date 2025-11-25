package getman

import (
	"net/http"
	"testing"

	"github.com/KonnorFrik/getman/testutil"
	"github.com/KonnorFrik/getman/types"
)

func TestIntegrationExecuteRequest(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := &types.Request{
		Method: http.MethodGet,
		URL:    testutil.GetServerURL() + "/health",
	}

	execution, err := client.ExecuteRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if execution.Response == nil {
		t.Fatal("expected response to be set")
	}

	if execution.Response.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", execution.Response.StatusCode)
	}
}

func TestIntegrationExecuteRequest_WithVariables(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	client.SetGlobalVariable("baseUrl", testutil.GetServerURL())

	req := &types.Request{
		Method: http.MethodGet,
		URL:    "{{baseUrl}}/health",
	}

	execution, err := client.ExecuteRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if execution.Response == nil {
		t.Fatal("expected response to be set")
	}
}

func TestIntegrationExecuteCollection_FullFlow(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collection := testutil.CreateTestCollection("Test Collection", []*types.RequestItem{
		{
			Name: "Test Request",
			Request: &types.Request{
				Method: http.MethodGet,
				URL:    testutil.GetServerURL() + "/health",
			},
		},
	})

	if err := client.SaveCollection(collection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := client.ExecuteCollection("Test Collection")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Statistics.Total)
	}
}

func TestIntegrationExecuteCollection_WithEnvironment(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env := testutil.CreateTestEnvironment("test", map[string]string{
		"baseUrl": testutil.GetServerURL(),
	})

	if err := client.SaveEnvironment(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.LoadLocalEnvironment("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collection := testutil.CreateTestCollection("Test Collection", []*types.RequestItem{
		{
			Name: "Test Request",
			Request: &types.Request{
				Method: http.MethodGet,
				URL:    "{{baseUrl}}/health",
			},
		},
	})

	if err := client.SaveCollection(collection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := client.ExecuteCollection("Test Collection")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Statistics.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Statistics.Total)
	}
}

func TestIntegrationExecuteCollection_History(t *testing.T) {
	_, err := testutil.StartTestServer()
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer testutil.StopTestServer()

	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collection := testutil.CreateTestCollection("Test Collection", []*types.RequestItem{
		{
			Name: "Test Request",
			Request: &types.Request{
				Method: http.MethodGet,
				URL:    testutil.GetServerURL() + "/health",
			},
		},
	})

	if err := client.SaveCollection(collection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := client.ExecuteCollection("Test Collection")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.SaveHistory(result); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	history, err := client.GetHistory(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(history) == 0 {
		t.Error("expected history to contain items")
	}
}

func TestIntegrationVariableResolution(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	client.SetGlobalVariable("globalVar", "globalValue")

	env := testutil.CreateTestEnvironment("test", map[string]string{
		"envVar": "envValue",
	})

	if err := client.SaveEnvironment(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.LoadGlobalEnvironment(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	globalResult, err := client.ResolveVariables("{{globalVar}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if globalResult != "globalValue" {
		t.Errorf("expected 'globalValue', got %s", globalResult)
	}

	envResult, err := client.ResolveVariables("{{envVar}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if envResult != "envValue" {
		t.Errorf("expected 'envValue', got %s", envResult)
	}
}

func TestIntegrationEnvironmentManagement(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env1 := testutil.CreateTestEnvironment("env1", map[string]string{"key1": "value1"})
	env2 := testutil.CreateTestEnvironment("env2", map[string]string{"key2": "value2"})

	if err := client.SaveEnvironment(env1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.SaveEnvironment(env2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	environments, err := client.ListEnvironments()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(environments) != 2 {
		t.Errorf("expected 2 environments, got %d", len(environments))
	}

	if err := client.LoadLocalEnvironment("env1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	currentEnv := client.GetCurrentEnvironment()
	if currentEnv == nil {
		t.Fatal("expected environment to be loaded")
	}

	if currentEnv.Name != "env1" {
		t.Errorf("expected environment name 'env1', got %s", currentEnv.Name)
	}

	if err := client.DeleteEnvironment("env1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	environments, err = client.ListEnvironments()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(environments) != 1 {
		t.Errorf("expected 1 environment, got %d", len(environments))
	}
}

func TestIntegrationCollectionManagement(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collection1 := testutil.CreateTestCollection("collection1", []*types.RequestItem{})
	collection2 := testutil.CreateTestCollection("collection2", []*types.RequestItem{})

	if err := client.SaveCollection(collection1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.SaveCollection(collection2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collections, err := client.ListCollections()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(collections) != 2 {
		t.Errorf("expected 2 collections, got %d", len(collections))
	}

	loaded, err := client.LoadCollection("collection1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if loaded.Name != "collection1" {
		t.Errorf("expected collection name 'collection1', got %s", loaded.Name)
	}

	if err := client.DeleteCollection("collection1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collections, err = client.ListCollections()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(collections) != 1 {
		t.Errorf("expected 1 collection, got %d", len(collections))
	}
}

