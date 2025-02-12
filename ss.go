package sortedslice

import (
	"encoding/gob"
	"os"
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

// SortedSlice is a thread-safe, sorted slice that stores key-value pairs.
type SortedSlice[K constraints.Ordered, V any] struct {
	sync.RWMutex
	data []kv[K, V]
}

type kv[K constraints.Ordered, V any] struct {
	Key   K
	Value V
}

// New creates a new SortedSlice.
func New[K constraints.Ordered, V any]() *SortedSlice[K, V] {
	return &SortedSlice[K, V]{}
}

// Add adds a value to the slice associated with the given key.
func (ss *SortedSlice[K, V]) Add(key K, value V) {
	ss.Lock()
	defer ss.Unlock()

	index := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key >= key
	})

	if index < len(ss.data) && ss.data[index].Key == key {
		// Key already exists, replace the value
		ss.data[index].Value = value
	} else {
		// Insert new key-value pair
		ss.data = append(ss.data, kv[K, V]{})
		copy(ss.data[index+1:], ss.data[index:])
		ss.data[index] = kv[K, V]{key, value}
	}
}

// Set is an alias for Add.
func (ss *SortedSlice[K, V]) Set(key K, value V) {
	ss.Add(key, value)
}

// Store is an alias for Add.
func (ss *SortedSlice[K, V]) Store(key K, value V) {
	ss.Add(key, value)
}

// Get retrieves the value associated with the given key.
func (ss *SortedSlice[K, V]) Get(key K) (V, bool) {
	ss.RLock()
	defer ss.RUnlock()

	index := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key >= key
	})

	if index < len(ss.data) && ss.data[index].Key == key {
		return ss.data[index].Value, true
	}
	var zero V
	return zero, false
}

// Find is an alias for Get.
func (ss *SortedSlice[K, V]) Find(key K) (V, bool) {
	return ss.Get(key)
}

// Exist checks if a key exists in the slice.
func (ss *SortedSlice[K, V]) Exist(key K) bool {
	ss.RLock()
	defer ss.RUnlock()

	index := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key >= key
	})

	return index < len(ss.data) && ss.data[index].Key == key
}

// Delete removes a key-value pair from the slice.
func (ss *SortedSlice[K, V]) Delete(key K) (V, bool) {
	ss.Lock()
	defer ss.Unlock()

	index := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key >= key
	})

	if index < len(ss.data) && ss.data[index].Key == key {
		value := ss.data[index].Value
		ss.data = append(ss.data[:index], ss.data[index+1:]...)
		return value, true
	}
	var zero V
	return zero, false
}

// Remove is an alias for Delete.
func (ss *SortedSlice[K, V]) Remove(key K) (V, bool) {
	return ss.Delete(key)
}

// Len returns the number of elements in the slice.
func (ss *SortedSlice[K, V]) Len() int {
	ss.RLock()
	defer ss.RUnlock()

	return len(ss.data)
}

// Clear removes all elements from the slice.
func (ss *SortedSlice[K, V]) Clear() {
	ss.Lock()
	defer ss.Unlock()

	ss.data = nil
}

// Load reads the slice from a file using gob encoding.
func (ss *SortedSlice[K, V]) Load(filename string) error {
	ss.Lock()
	defer ss.Unlock()

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(&ss.data)
}

// Save writes the slice to a file using gob encoding.
func (ss *SortedSlice[K, V]) Save(filename string) error {
	ss.RLock()
	defer ss.RUnlock()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(ss.data)
}

// Range iterates over all key-value pairs in the slice.
func (ss *SortedSlice[K, V]) Range(f func(k K, v V) bool) {
	ss.RLock()
	defer ss.RUnlock()

	for _, kv := range ss.data {
		if !f(kv.Key, kv.Value) {
			break
		}
	}
}

// Range iterates over all key-value pairs in the slice but iterates backwards
func (ss *SortedSlice[K, V]) RangeBackward(f func(k K, v V) bool) {
	ss.RLock()
	defer ss.RUnlock()

	for i := 0; i < len(ss.data); i++ {
		kv := ss.data[len(ss.data)-i-1]
		if !f(kv.Key, kv.Value) {
			break
		}
	}
}

// First returns the first element in the slice.
func (ss *SortedSlice[K, V]) First() V {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.data) == 0 {
		var zero V
		return zero
	}
	return ss.data[0].Value
}

// Last returns the last element in the slice.
func (ss *SortedSlice[K, V]) Last() V {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.data) == 0 {
		var zero V
		return zero
	}
	return ss.data[len(ss.data)-1].Value
}

// Min is an alias for First.
func (ss *SortedSlice[K, V]) Min() V {
	return ss.First()
}

// Max is an alias for Last.
func (ss *SortedSlice[K, V]) Max() V {
	return ss.Last()
}

// FirstKey returns the first key in the slice.
func (ss *SortedSlice[K, V]) FirstKey() K {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.data) == 0 {
		var zero K
		return zero
	}
	return ss.data[0].Key
}

// LastKey returns the last key in the slice.
func (ss *SortedSlice[K, V]) LastKey() K {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.data) == 0 {
		var zero K
		return zero
	}
	return ss.data[len(ss.data)-1].Key
}

// MinKey is an alias for FirstKey.
func (ss *SortedSlice[K, V]) MinKey() K {
	return ss.FirstKey()
}

// MaxKey is an alias for LastKey.
func (ss *SortedSlice[K, V]) MaxKey() K {
	return ss.LastKey()
}
