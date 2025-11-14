package getman

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KonnorFrik/getman/testutil"
	"github.com/KonnorFrik/getman/types"
)

func TestUnitNewClient_Valid(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client == nil {
		t.Fatal("expected client to be created")
	}
}

func TestUnitNewClientWithDefaults(t *testing.T) {
	client, err := NewClientWithDefaults()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if client == nil {
		t.Fatal("expected client to be created")
	}
}

func TestUnitLoadEnvironment(t *testing.T) {
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
		"key1": "value1",
	})

	if err := client.SaveEnvironment(env); err != nil {
		t.Fatalf("unexpected error saving environment: %v", err)
	}

	if err := client.LoadEnvironment("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	currentEnv := client.GetCurrentEnvironment()
	if currentEnv == nil {
		t.Fatal("expected environment to be loaded")
	}

	if currentEnv.Name != "test" {
		t.Errorf("expected environment name 'test', got %s", currentEnv.Name)
	}
}

func TestUnitLoadEnvironment_NotFound(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = client.LoadEnvironment("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent environment")
	}
}

func TestUnitSaveEnvironment(t *testing.T) {
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
		"key1": "value1",
	})

	if err := client.SaveEnvironment(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitListEnvironments(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env1 := testutil.CreateTestEnvironment("env1", map[string]string{})
	env2 := testutil.CreateTestEnvironment("env2", map[string]string{})

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
		t.Errorf("expected 2 environments, got %d: %v", len(environments), environments)
	}
}

func TestUnitDeleteEnvironment(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	env := testutil.CreateTestEnvironment("test", map[string]string{})

	if err := client.SaveEnvironment(env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.DeleteEnvironment("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitSetGlobalVariable(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	client.SetGlobalVariable("testKey", "testValue")

	value, ok := client.GetGlobalVariable("testKey")
	if !ok {
		t.Fatal("expected variable to be found")
	}

	if value != "testValue" {
		t.Errorf("expected value 'testValue', got %s", value)
	}
}

func TestUnitGetVariable(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	client.SetGlobalVariable("testKey", "testValue")

	value, ok := client.GetVariable("testKey")
	if !ok {
		t.Fatal("expected variable to be found")
	}

	if value != "testValue" {
		t.Errorf("expected value 'testValue', got %s", value)
	}
}

func TestUnitResolveVariables(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	client.SetGlobalVariable("testVar", "testValue")

	result, err := client.ResolveVariables("{{testVar}}")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "testValue" {
		t.Errorf("expected result 'testValue', got %s", result)
	}
}

func TestUnitLoadCollection(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collection := testutil.CreateTestCollection("test", []*types.RequestItem{
		{
			Name: "Test Request",
			Request: &types.Request{
				Method: "GET",
				URL:    "http://example.com",
			},
		},
	})

	if err := client.SaveCollection(collection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	loaded, err := client.LoadCollection("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if loaded.Name != "test" {
		t.Errorf("expected collection name 'test', got %s", loaded.Name)
	}
}

func TestUnitLoadCollection_NotFound(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = client.LoadCollection("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent collection")
	}
}

func TestUnitSaveCollection(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collection := testutil.CreateTestCollection("test", []*types.RequestItem{
		{
			Name: "Test Request",
			Request: &types.Request{
				Method: "GET",
				URL:    "http://example.com",
			},
		},
	})

	if err := client.SaveCollection(collection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitListCollections(t *testing.T) {
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
}

func TestUnitDeleteCollection(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	collection := testutil.CreateTestCollection("test", []*types.RequestItem{})

	if err := client.SaveCollection(collection); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := client.DeleteCollection("test"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitValidateRequest_Valid(t *testing.T) {
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
		Method: "GET",
		URL:    "http://example.com",
	}

	if err := client.ValidateRequest(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUnitValidateRequest_Invalid(t *testing.T) {
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
		Method: "",
		URL:    "http://example.com",
	}

	err = client.ValidateRequest(req)
	if err == nil {
		t.Fatal("expected error for invalid request")
	}
}

func TestUnitExecuteRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer server.Close()

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
		Method: "GET",
		URL:    server.URL,
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

func TestUnitGetHistory(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := &types.ExecutionResult{
		CollectionName: "Test Collection",
		Environment:    "test",
		StartTime:      time.Now(),
		EndTime:        time.Now(),
		TotalDuration:  time.Second,
		Requests:       []*types.RequestExecution{},
		Statistics: &types.Statistics{
			Total:   0,
			Success: 0,
			Failed:  0,
			AvgTime: 0,
			MinTime: 0,
			MaxTime: 0,
		},
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

func TestUnitGetConfig(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	config := client.GetConfig()
	if config == nil {
		t.Fatal("expected config to be set")
	}
}

func TestUnitUpdateConfig(t *testing.T) {
	dir, err := testutil.CreateTempDir()
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer testutil.CleanupTempDir(dir)

	client, err := NewClient(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	config := DefaultConfig()
	config.Defaults.Timeout.Connect = 60 * time.Second

	if err := client.UpdateConfig(config); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	updatedConfig := client.GetConfig()
	if updatedConfig.Defaults.Timeout.Connect != 60*time.Second {
		t.Errorf("expected connect timeout 60s, got %v", updatedConfig.Defaults.Timeout.Connect)
	}
}

