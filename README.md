# SortedSlice

**SortedSlice** is a thread-safe, generic Go library that provides a sorted slice for storing key-value pairs. It is designed to be fast, efficient, and easy to use, with support for common operations like adding, deleting, and retrieving elements while maintaining sorted order.

## Features

- **Thread-safe**: Uses `sync.RWMutex` to ensure safe concurrent access.
- **Sorted storage**: Keys are stored in sorted order, enabling fast lookups and range queries.
- **Generic**: Works with any ordered key type (`constraints.Ordered`) and any value type.
- **Common operations**:
  - Add, Set, Store: Insert or update key-value pairs.
  - Get, Find: Retrieve values by key.
  - Delete, Remove: Remove key-value pairs.
  - Exist: Check if a key exists.
  - Len: Get the number of elements.
  - Clear: Remove all elements.
  - Range: Iterate over all key-value pairs in sorted order.
  - First, Last, Min, Max: Get the first or last element (by key or value).
  - Iter: Get an iterator for looping through elements.
- **Serialization**: Save and load data to/from a file using `gob` encoding.

## Installation

To use `SortedSlice` in your Go project, run:

```bash
go get github.com/goupdate/sortedslice
```

Replace `yourusername` with your GitHub username or the actual import path.

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"github.com/goupdate/sortedslice"
)

func main() {
	// Create a new SortedSlice
	ss := sortedslice.New[int, string]()

	// Add key-value pairs
	ss.Add(3, "three")
	ss.Add(1, "one")
	ss.Add(2, "two")

	// Retrieve a value
	val, found := ss.Get(2)
	if found {
		fmt.Println("Found:", val) // Output: Found: two
	}

	// Iterate over all elements
	ss.Range(func(k int, v string) bool {
		fmt.Printf("Key: %d, Value: %s\n", k, v)
		return true
	})

	// Delete a key
	ss.Delete(2)

	// Check if a key exists
	if ss.Exist(2) {
		fmt.Println("Key 2 exists")
	} else {
		fmt.Println("Key 2 does not exist") // Output: Key 2 does not exist
	}
}
```

### Advanced Example

```go
package main

import (
	"fmt"
	"github.com/goupdate/sortedslice"
)

func main() {
	ss := sortedslice.New[int, string]()

	// Add multiple values
	ss.Add(5, "five")
	ss.Add(4, "four")
	ss.Add(6, "six")

	// Get the first and last elements
	fmt.Println("First:", ss.First()) // Output: First: four
	fmt.Println("Last:", ss.Last())   // Output: Last: six

	// Get the first and last keys
	fmt.Println("First Key:", ss.FirstKey()) // Output: First Key: 4
	fmt.Println("Last Key:", ss.LastKey())   // Output: Last Key: 6

	// Save to a file
	err := ss.Save("data.gob")
	if err != nil {
		fmt.Println("Error saving:", err)
	}

	// Load from a file
	newSS := sortedslice.New[int, string]()
	err = newSS.Load("data.gob")
	if err != nil {
		fmt.Println("Error loading:", err)
	}

	// Verify loaded data
	newSS.Range(func(k int, v string) bool {
		fmt.Printf("Loaded Key: %d, Value: %s\n", k, v)
		return true
	})
}
```

## API Reference

### Types

- **SortedSlice[K, V]**: The main type representing a sorted slice of key-value pairs.

### Methods

- **New[K, V]()**: Creates a new `SortedSlice`.
- **Add(k K, v V)**: Adds a key-value pair (or updates if the key exists).
- **Set(k K, v V)**: Alias for `Add`.
- **Store(k K, v V)**: Alias for `Add`.
- **Get(k K) (V, bool)**: Retrieves a value by key.
- **Find(k K) (V, bool)**: Alias for `Get`.
- **Exist(k K) bool**: Checks if a key exists.
- **Delete(k K) (V, bool)**: Removes a key-value pair.
- **Remove(k K) (V, bool)**: Alias for `Delete`.
- **Len() int**: Returns the number of elements.
- **Clear()**: Removes all elements.
- **Save(filename string) error**: Saves the data to a file.
- **Load(filename string) error**: Loads the data from a file.
- **Range(f func(k K, v V) bool)**: Iterates over all key-value pairs.
- **First() V**: Returns the first value.
- **Last() V**: Returns the last value.
- **Min() V**: Alias for `First`.
- **Max() V**: Alias for `Last`.
- **FirstKey() K**: Returns the first key.
- **LastKey() K**: Returns the last key.
- **MinKey() K**: Alias for `FirstKey`.
- **MaxKey() K**: Alias for `LastKey`.
- **Iter() <-chan kv[K, V]**: Returns an iterator for looping through elements.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.