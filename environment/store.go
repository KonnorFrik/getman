package environment
//
// import (
// 	"sync"
// )
//
// type VariableStore struct {
// 	mu         sync.RWMutex
// 	globalVars map[string]string
// 	envVars    map[string]string
// }
//
// func NewVariableStore() *VariableStore {
// 	return &VariableStore{
// 		globalVars: make(map[string]string),
// 		envVars:    make(map[string]string),
// 	}
// }
//
// func (vs *VariableStore) SetGlobal(key, value string) {
// 	vs.mu.Lock()
// 	defer vs.mu.Unlock()
// 	vs.globalVars[key] = value
// }
//
// func (vs *VariableStore) GetGlobal(key string) (string, bool) {
// 	vs.mu.RLock()
// 	defer vs.mu.RUnlock()
// 	value, ok := vs.globalVars[key]
// 	return value, ok
// }
//
// func (vs *VariableStore) SetEnv(key, value string) {
// 	vs.mu.Lock()
// 	defer vs.mu.Unlock()
// 	vs.envVars[key] = value
// }
//
// func (vs *VariableStore) GetEnv(key string) (string, bool) {
// 	vs.mu.RLock()
// 	defer vs.mu.RUnlock()
// 	value, ok := vs.envVars[key]
// 	return value, ok
// }
//
// func (vs *VariableStore) Get(key string) (string, bool) {
// 	vs.mu.RLock()
// 	defer vs.mu.RUnlock()
//
// 	if value, ok := vs.envVars[key]; ok {
// 		return value, true
// 	}
//
// 	if value, ok := vs.globalVars[key]; ok {
// 		return value, true
// 	}
//
// 	return "", false
// }
//
// func (vs *VariableStore) DeleteGlobal(key string) {
// 	vs.mu.Lock()
// 	defer vs.mu.Unlock()
// 	delete(vs.globalVars, key)
// }
//
// func (vs *VariableStore) DeleteEnv(key string) {
// 	vs.mu.Lock()
// 	defer vs.mu.Unlock()
// 	delete(vs.envVars, key)
// }
//
// func (vs *VariableStore) ClearGlobal() {
// 	vs.mu.Lock()
// 	defer vs.mu.Unlock()
// 	vs.globalVars = make(map[string]string)
// }
//
// func (vs *VariableStore) ClearEnv() {
// 	vs.mu.Lock()
// 	defer vs.mu.Unlock()
// 	vs.envVars = make(map[string]string)
// }
//
// func (vs *VariableStore) SetEnvVars(vars map[string]string) {
// 	vs.mu.Lock()
// 	defer vs.mu.Unlock()
// 	vs.envVars = make(map[string]string)
// 	for k, v := range vars {
// 		vs.envVars[k] = v
// 	}
// }
//
// func (vs *VariableStore) GetAll() map[string]string {
// 	vs.mu.RLock()
// 	defer vs.mu.RUnlock()
//
// 	result := make(map[string]string)
//
// 	for k, v := range vs.globalVars {
// 		result[k] = v
// 	}
//
// 	for k, v := range vs.envVars {
// 		result[k] = v
// 	}
//
// 	return result
// }
