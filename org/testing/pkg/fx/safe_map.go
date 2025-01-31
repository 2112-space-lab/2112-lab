package fx

import "sync"

// SafeMap a map with mutex lock, for thread safe usage
// TODO: replace with sync.Map once implementation with generics is available
type SafeMap[K comparable, V interface{}] struct {
	m       sync.RWMutex
	safeMap map[K]V
}

// NewSafeMap returns a new SafeMap map with mutex lock, for thread safe usage
func NewSafeMap[K comparable, V interface{}]() SafeMap[K, V] {
	return SafeMap[K, V]{
		safeMap: make(map[K]V),
		m:       sync.RWMutex{},
	}
}

// Read reads key from map, returns value and boolean indicating if value exists
func (sm *SafeMap[K, V]) Read(key K) (V, bool) {
	sm.m.RLock()
	defer sm.m.RUnlock()
	val, ok := sm.safeMap[key]
	return val, ok
}

// Write writes a key value pair into map
func (sm *SafeMap[K, V]) Write(key K, val V) {
	sm.m.Lock()
	defer sm.m.Unlock()
	sm.safeMap[key] = val
}

// Delete deletes a key value pair from map. noop if no value
func (sm *SafeMap[K, V]) Delete(key K) {
	sm.m.Lock()
	defer sm.m.Unlock()
	delete(sm.safeMap, key)
}

// Clear deletes all key value pairs from map
func (sm *SafeMap[K, V]) Clear() {
	sm.m.Lock()
	defer sm.m.Unlock()
	sm.safeMap = map[K]V{}
}

// Size returns number of elements in map
func (sm *SafeMap[K, V]) Size() int {
	sm.m.Lock()
	defer sm.m.Unlock()
	return len(sm.safeMap)
}

// CopyToSnapshot returns a new map with current values of safe map
func (sm *SafeMap[K, V]) CopyToSnapshot() map[K]V {
	sm.m.Lock()
	defer sm.m.Unlock()
	snapshot := make(map[K]V)
	for k, v := range sm.safeMap {
		snapshot[k] = v
	}
	return snapshot
}
