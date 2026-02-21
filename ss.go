package sortedslice

import (
	"encoding/gob"
	"os"
	"sort"
	"sync"

	"golang.org/x/exp/constraints"
)

// SortedSlice is a thread-safe, sorted slice that stores key-value pairs.
// All operations are protected by RWMutex for concurrent access.
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
	return &SortedSlice[K, V]{
		data: make([]kv[K, V], 0),
	}
}

// Add adds a value to the slice associated with the given key.
// If key exists, value is replaced. O(N) due to slice insertion.
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

// Get retrieves the value associated with the given key. O(log N)
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

// Exist checks if a key exists in the slice. O(log N)
func (ss *SortedSlice[K, V]) Exist(key K) bool {
	ss.RLock()
	defer ss.RUnlock()

	index := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key >= key
	})

	return index < len(ss.data) && ss.data[index].Key == key
}

// Delete removes a key-value pair from the slice. O(N)
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

// Len returns the number of elements in the slice. O(1)
func (ss *SortedSlice[K, V]) Len() int {
	ss.RLock()
	defer ss.RUnlock()

	return len(ss.data)
}

// LenNoLock returns the number of elements without locking.
// Use only when you're sure no concurrent modifications happen.
func (ss *SortedSlice[K, V]) LenNoLock() int {
	return len(ss.data)
}

// Clear removes all elements from the slice.
func (ss *SortedSlice[K, V]) Clear() {
	ss.Lock()
	defer ss.Unlock()

	ss.data = make([]kv[K, V], 0)
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

// Range iterates over all key-value pairs in the slice in ascending order.
// Callback can return false to stop iteration.
func (ss *SortedSlice[K, V]) Range(f func(k K, v V) bool) {
	ss.RLock()
	defer ss.RUnlock()

	for _, kv := range ss.data {
		if !f(kv.Key, kv.Value) {
			break
		}
	}
}

// RangeBackward iterates over all key-value pairs in descending order.
// FIXED: Now handles empty slice without panic.
func (ss *SortedSlice[K, V]) RangeBackward(f func(k K, v V) bool) {
	ss.RLock()
	defer ss.RUnlock()

	l := len(ss.data)
	if l == 0 {
		return
	}
	for i := l - 1; i >= 0; i-- {
		kv := ss.data[i]
		if !f(kv.Key, kv.Value) {
			break
		}
	}
}

// First returns the first element in the slice (lowest key).
func (ss *SortedSlice[K, V]) First() V {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.data) == 0 {
		var zero V
		return zero
	}
	return ss.data[0].Value
}

// Last returns the last element in the slice (highest key).
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

// FirstKey returns the first key in the slice (lowest key).
func (ss *SortedSlice[K, V]) FirstKey() K {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.data) == 0 {
		var zero K
		return zero
	}
	return ss.data[0].Key
}

// LastKey returns the last key in the slice (highest key).
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

// ============================================================================
// NEW OPTIMIZED METHODS - O(log N) binary search operations
// ============================================================================

// LowerBound returns the key and value of the element with the largest key
// that is strictly LESS than the given key.
// Returns zero values and false if no such element exists.
// Time complexity: O(log N)
func (ss *SortedSlice[K, V]) LowerBound(key K) (K, V, bool) {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.data) == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	// Find first element >= key
	index := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key >= key
	})

	// Element before index is < key
	if index > 0 {
		kv := ss.data[index-1]
		return kv.Key, kv.Value, true
	}

	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// UpperBound returns the key and value of the element with the smallest key
// that is strictly GREATER than the given key.
// Returns zero values and false if no such element exists.
// Time complexity: O(log N)
func (ss *SortedSlice[K, V]) UpperBound(key K) (K, V, bool) {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.data) == 0 {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	// Find first element > key
	index := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key > key
	})

	if index < len(ss.data) {
		kv := ss.data[index]
		return kv.Key, kv.Value, true
	}

	var zeroK K
	var zeroV V
	return zeroK, zeroV, false
}

// FindLower returns the value of the element with the largest key
// that is strictly LESS than the given key.
// Time complexity: O(log N)
func (ss *SortedSlice[K, V]) FindLower(key K) (V, bool) {
	_, v, ok := ss.LowerBound(key)
	return v, ok
}

// FindHigher returns the value of the element with the smallest key
// that is strictly GREATER than the given key.
// Time complexity: O(log N)
func (ss *SortedSlice[K, V]) FindHigher(key K) (V, bool) {
	_, v, ok := ss.UpperBound(key)
	return v, ok
}

// FindLowerWithKey returns the key and value of the element with the largest key
// that is strictly LESS than the given key.
// Time complexity: O(log N)
func (ss *SortedSlice[K, V]) FindLowerWithKey(key K) (K, V, bool) {
	return ss.LowerBound(key)
}

// FindHigherWithKey returns the key and value of the element with the smallest key
// that is strictly GREATER than the given key.
// Time complexity: O(log N)
func (ss *SortedSlice[K, V]) FindHigherWithKey(key K) (K, V, bool) {
	return ss.UpperBound(key)
}

// FindRange returns all key-value pairs where key is in range [start, end].
// Time complexity: O(log N + M) where M is number of elements in range
func (ss *SortedSlice[K, V]) FindRange(start, end K) []kv[K, V] {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.data) == 0 || start > end {
		return []kv[K, V]{}
	}

	// Find first element >= start
	startIndex := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key >= start
	})

	if startIndex >= len(ss.data) {
		return []kv[K, V]{}
	}

	// Collect elements in range
	result := make([]kv[K, V], 0)
	for i := startIndex; i < len(ss.data); i++ {
		if ss.data[i].Key > end {
			break
		}
		result = append(result, ss.data[i])
	}

	return result
}

// FindRangeCallback iterates over all key-value pairs where key is in range [start, end].
// Callback can return false to stop iteration.
// Time complexity: O(log N + M) where M is number of elements in range
func (ss *SortedSlice[K, V]) FindRangeCallback(start, end K, f func(k K, v V) bool) {
	ss.RLock()
	defer ss.RUnlock()

	if len(ss.data) == 0 || start > end {
		return
	}

	// Find first element >= start
	startIndex := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key >= start
	})

	if startIndex >= len(ss.data) {
		return
	}

	// Iterate elements in range
	for i := startIndex; i < len(ss.data); i++ {
		if ss.data[i].Key > end {
			break
		}
		if !f(ss.data[i].Key, ss.data[i].Value) {
			break
		}
	}
}

// IndexOf returns the index of the element with the given key.
// Returns -1 if not found.
// Time complexity: O(log N)
func (ss *SortedSlice[K, V]) IndexOf(key K) int {
	ss.RLock()
	defer ss.RUnlock()

	index := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key >= key
	})

	if index < len(ss.data) && ss.data[index].Key == key {
		return index
	}
	return -1
}

// At returns the key-value pair at the given index.
// Returns zero values and false if index is out of bounds.
// Time complexity: O(1)
func (ss *SortedSlice[K, V]) At(index int) (K, V, bool) {
	ss.RLock()
	defer ss.RUnlock()

	if index < 0 || index >= len(ss.data) {
		var zeroK K
		var zeroV V
		return zeroK, zeroV, false
	}

	kv := ss.data[index]
	return kv.Key, kv.Value, true
}

// Keys returns all keys in ascending order.
// Time complexity: O(N)
func (ss *SortedSlice[K, V]) Keys() []K {
	ss.RLock()
	defer ss.RUnlock()

	keys := make([]K, len(ss.data))
	for i, kv := range ss.data {
		keys[i] = kv.Key
	}
	return keys
}

// Values returns all values in ascending key order.
// Time complexity: O(N)
func (ss *SortedSlice[K, V]) Values() []V {
	ss.RLock()
	defer ss.RUnlock()

	values := make([]V, len(ss.data))
	for i, kv := range ss.data {
		values[i] = kv.Value
	}
	return values
}

// Clone creates a shallow copy of the SortedSlice.
// Time complexity: O(N)
func (ss *SortedSlice[K, V]) Clone() *SortedSlice[K, V] {
	ss.RLock()
	defer ss.RUnlock()

	newSlice := &SortedSlice[K, V]{
		data: make([]kv[K, V], len(ss.data)),
	}
	copy(newSlice.data, ss.data)
	return newSlice
}

// IsEmpty returns true if the slice contains no elements.
// Time complexity: O(1)
func (ss *SortedSlice[K, V]) IsEmpty() bool {
	ss.RLock()
	defer ss.RUnlock()
	return len(ss.data) == 0
}

// IsEmptyNoLock returns true if the slice contains no elements without locking.
// Use only when you're sure no concurrent modifications happen.
func (ss *SortedSlice[K, V]) IsEmptyNoLock() bool {
	return len(ss.data) == 0
}

// ============================================================================
// BULK OPERATIONS - More efficient for multiple operations
// ============================================================================

// BulkAdd adds multiple key-value pairs in a single operation.
// More efficient than calling Add multiple times (single lock acquisition).
// Time complexity: O(N*M) where N is existing size, M is new items
func (ss *SortedSlice[K, V]) BulkAdd(items []kv[K, V]) {
	if len(items) == 0 {
		return
	}

	ss.Lock()
	defer ss.Unlock()

	// Sort incoming items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key < items[j].Key
	})

	// Merge with existing data
	result := make([]kv[K, V], 0, len(ss.data)+len(items))
	i, j := 0, 0

	for i < len(ss.data) && j < len(items) {
		if ss.data[i].Key < items[j].Key {
			result = append(result, ss.data[i])
			i++
		} else if ss.data[i].Key > items[j].Key {
			result = append(result, items[j])
			j++
		} else {
			// Same key, use new value
			result = append(result, items[j])
			i++
			j++
		}
	}

	for i < len(ss.data) {
		result = append(result, ss.data[i])
		i++
	}

	for j < len(items) {
		result = append(result, items[j])
		j++
	}

	ss.data = result
}

// DeleteRange removes all elements with keys in range [start, end].
// Returns number of deleted elements.
// Time complexity: O(N)
func (ss *SortedSlice[K, V]) DeleteRange(start, end K) int {
	ss.Lock()
	defer ss.Unlock()

	if len(ss.data) == 0 || start > end {
		return 0
	}

	// Find first element >= start
	startIndex := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key >= start
	})

	if startIndex >= len(ss.data) {
		return 0
	}

	// Find first element > end
	endIndex := sort.Search(len(ss.data), func(i int) bool {
		return ss.data[i].Key > end
	})

	deleted := endIndex - startIndex
	ss.data = append(ss.data[:startIndex], ss.data[endIndex:]...)

	return deleted
}
