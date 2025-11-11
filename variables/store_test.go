package variables

import (
	"fmt"
	"sync"
	"testing"
)

func TestUnitNewVariableStore(t *testing.T) {
	store := NewVariableStore()
	if store == nil {
		t.Fatal("expected variable store to be created")
	}
}

func TestUnitSetGlobal(t *testing.T) {
	store := NewVariableStore()
	store.SetGlobal("testKey", "testValue")

	value, ok := store.GetGlobal("testKey")
	if !ok {
		t.Fatal("expected key to be found")
	}

	if value != "testValue" {
		t.Errorf("expected value 'testValue', got %s", value)
	}
}

func TestUnitGetGlobal(t *testing.T) {
	store := NewVariableStore()
	store.SetGlobal("testKey", "testValue")

	value, ok := store.GetGlobal("testKey")
	if !ok {
		t.Fatal("expected key to be found")
	}

	if value != "testValue" {
		t.Errorf("expected value 'testValue', got %s", value)
	}
}

func TestUnitGetGlobal_NotFound(t *testing.T) {
	store := NewVariableStore()

	_, ok := store.GetGlobal("nonexistent")
	if ok {
		t.Fatal("expected key to not be found")
	}
}

func TestUnitSetEnv(t *testing.T) {
	store := NewVariableStore()
	store.SetEnv("testKey", "testValue")

	value, ok := store.GetEnv("testKey")
	if !ok {
		t.Fatal("expected key to be found")
	}

	if value != "testValue" {
		t.Errorf("expected value 'testValue', got %s", value)
	}
}

func TestUnitGetEnv(t *testing.T) {
	store := NewVariableStore()
	store.SetEnv("testKey", "testValue")

	value, ok := store.GetEnv("testKey")
	if !ok {
		t.Fatal("expected key to be found")
	}

	if value != "testValue" {
		t.Errorf("expected value 'testValue', got %s", value)
	}
}

func TestUnitGetEnv_NotFound(t *testing.T) {
	store := NewVariableStore()

	_, ok := store.GetEnv("nonexistent")
	if ok {
		t.Fatal("expected key to not be found")
	}
}

func TestUnitGet_EnvPriority(t *testing.T) {
	store := NewVariableStore()
	store.SetGlobal("testKey", "globalValue")
	store.SetEnv("testKey", "envValue")

	value, ok := store.Get("testKey")
	if !ok {
		t.Fatal("expected key to be found")
	}

	if value != "envValue" {
		t.Errorf("expected env value 'envValue', got %s", value)
	}
}

func TestUnitGet_GlobalFallback(t *testing.T) {
	store := NewVariableStore()
	store.SetGlobal("testKey", "globalValue")

	value, ok := store.Get("testKey")
	if !ok {
		t.Fatal("expected key to be found")
	}

	if value != "globalValue" {
		t.Errorf("expected global value 'globalValue', got %s", value)
	}
}

func TestUnitDeleteGlobal(t *testing.T) {
	store := NewVariableStore()
	store.SetGlobal("testKey", "testValue")
	store.DeleteGlobal("testKey")

	_, ok := store.GetGlobal("testKey")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestUnitDeleteEnv(t *testing.T) {
	store := NewVariableStore()
	store.SetEnv("testKey", "testValue")
	store.DeleteEnv("testKey")

	_, ok := store.GetEnv("testKey")
	if ok {
		t.Fatal("expected key to be deleted")
	}
}

func TestUnitClearGlobal(t *testing.T) {
	store := NewVariableStore()
	store.SetGlobal("key1", "value1")
	store.SetGlobal("key2", "value2")
	store.ClearGlobal()

	_, ok := store.GetGlobal("key1")
	if ok {
		t.Fatal("expected key1 to be cleared")
	}

	_, ok = store.GetGlobal("key2")
	if ok {
		t.Fatal("expected key2 to be cleared")
	}
}

func TestUnitClearEnv(t *testing.T) {
	store := NewVariableStore()
	store.SetEnv("key1", "value1")
	store.SetEnv("key2", "value2")
	store.ClearEnv()

	_, ok := store.GetEnv("key1")
	if ok {
		t.Fatal("expected key1 to be cleared")
	}

	_, ok = store.GetEnv("key2")
	if ok {
		t.Fatal("expected key2 to be cleared")
	}
}

func TestUnitSetEnvVars(t *testing.T) {
	store := NewVariableStore()
	vars := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	store.SetEnvVars(vars)

	value1, ok := store.GetEnv("key1")
	if !ok {
		t.Fatal("expected key1 to be found")
	}
	if value1 != "value1" {
		t.Errorf("expected value1 'value1', got %s", value1)
	}

	value2, ok := store.GetEnv("key2")
	if !ok {
		t.Fatal("expected key2 to be found")
	}
	if value2 != "value2" {
		t.Errorf("expected value2 'value2', got %s", value2)
	}
}

func TestUnitSetEnvVars_Overwrite(t *testing.T) {
	store := NewVariableStore()
	store.SetEnv("key1", "oldValue")

	vars := map[string]string{
		"key1": "newValue",
	}
	store.SetEnvVars(vars)

	value, ok := store.GetEnv("key1")
	if !ok {
		t.Fatal("expected key1 to be found")
	}
	if value != "newValue" {
		t.Errorf("expected value 'newValue', got %s", value)
	}
}

func TestUnitGetAll(t *testing.T) {
	store := NewVariableStore()
	store.SetGlobal("globalKey", "globalValue")
	store.SetEnv("envKey", "envValue")

	all := store.GetAll()

	if all["globalKey"] != "globalValue" {
		t.Errorf("expected globalKey 'globalValue', got %s", all["globalKey"])
	}

	if all["envKey"] != "envValue" {
		t.Errorf("expected envKey 'envValue', got %s", all["envKey"])
	}

	if len(all) != 2 {
		t.Errorf("expected 2 variables, got %d", len(all))
	}
}

func TestUnitGetAll_EnvOverridesGlobal(t *testing.T) {
	store := NewVariableStore()
	store.SetGlobal("testKey", "globalValue")
	store.SetEnv("testKey", "envValue")

	all := store.GetAll()

	if all["testKey"] != "envValue" {
		t.Errorf("expected testKey 'envValue', got %s", all["testKey"])
	}

	if len(all) != 1 {
		t.Errorf("expected 1 variable, got %d", len(all))
	}
}

func TestUnitConcurrentAccess(t *testing.T) {
	store := NewVariableStore()
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	wg.Add(numGoroutines * 2)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key%d_%d", id, j)
				value := fmt.Sprintf("value%d_%d", id, j)
				store.SetGlobal(key, value)
			}
		}(i)

		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("key%d_%d", id, j)
				store.GetGlobal(key)
			}
		}(i)
	}

	wg.Wait()
}

func TestUnitGet_NotFound(t *testing.T) {
	store := NewVariableStore()

	_, ok := store.Get("nonexistent")
	if ok {
		t.Fatal("expected key to not be found")
	}
}

func TestUnitSetEnvVars_EmptyMap(t *testing.T) {
	store := NewVariableStore()
	store.SetEnv("key1", "value1")
	store.SetEnvVars(map[string]string{})

	_, ok := store.GetEnv("key1")
	if ok {
		t.Fatal("expected key1 to be cleared")
	}
}

