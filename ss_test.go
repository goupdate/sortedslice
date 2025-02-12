package sortedslice

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ss := New[int, string]()
	assert.NotNil(t, ss)
	assert.Equal(t, 0, ss.Len())
}

func TestAddAndGet(t *testing.T) {
	ss := New[int, string]()
	ss.Add(1, "one")
	ss.Add(2, "two")
	ss.Add(3, "three")

	val, found := ss.Get(2)
	assert.True(t, found)
	assert.Equal(t, "two", val)

	val, found = ss.Get(4)
	assert.False(t, found)
	assert.Empty(t, val)
}

func TestSetAndStore(t *testing.T) {
	ss := New[int, string]()
	ss.Set(1, "one")
	ss.Store(2, "two")

	val, found := ss.Get(1)
	assert.True(t, found)
	assert.Equal(t, "one", val)

	val, found = ss.Get(2)
	assert.True(t, found)
	assert.Equal(t, "two", val)
}

func TestExist(t *testing.T) {
	ss := New[int, string]()
	ss.Add(1, "one")
	ss.Add(2, "two")

	assert.True(t, ss.Exist(1))
	assert.False(t, ss.Exist(3))
}

func TestDeleteAndRemove(t *testing.T) {
	ss := New[int, string]()
	ss.Add(1, "one")
	ss.Add(2, "two")

	val, found := ss.Delete(1)
	assert.True(t, found)
	assert.Equal(t, "one", val)

	val, found = ss.Remove(2)
	assert.True(t, found)
	assert.Equal(t, "two", val)

	assert.Equal(t, 0, ss.Len())
}

func TestLen(t *testing.T) {
	ss := New[int, string]()
	ss.Add(1, "one")
	ss.Add(2, "two")

	assert.Equal(t, 2, ss.Len())
}

func TestClear(t *testing.T) {
	ss := New[int, string]()
	ss.Add(1, "one")
	ss.Add(2, "two")

	ss.Clear()
	assert.Equal(t, 0, ss.Len())
}

func TestLoadAndSave(t *testing.T) {
	ss := New[int, string]()
	ss.Add(1, "one")
	ss.Add(2, "two")

	filename := "test_save.gob"
	err := ss.Save(filename)
	assert.NoError(t, err)

	newSS := New[int, string]()
	err = newSS.Load(filename)
	assert.NoError(t, err)

	assert.Equal(t, ss.Len(), newSS.Len())
	val, found := newSS.Get(1)
	assert.True(t, found)
	assert.Equal(t, "one", val)

	os.Remove(filename)
}

func TestRange(t *testing.T) {
	ss := New[int, string]()
	ss.Add(1, "one")
	ss.Add(2, "two")
	ss.Add(3, "three")

	keys := []int{}
	values := []string{}
	ss.Range(func(k int, v string) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})

	assert.Equal(t, []int{1, 2, 3}, keys)
	assert.Equal(t, []string{"one", "two", "three"}, values)
}

func TestRangeBackward(t *testing.T) {
	ss := New[int, string]()
	ss.Add(1, "one")
	ss.Add(2, "two")
	ss.Add(3, "three")

	keys := []int{}
	values := []string{}
	ss.RangeBackward(func(k int, v string) bool {
		keys = append(keys, k)
		values = append(values, v)
		return true
	})

	assert.Equal(t, []int{3, 2, 1}, keys)
	assert.Equal(t, []string{"three", "two", "one"}, values)
}

func TestFirstAndLast(t *testing.T) {
	ss := New[int, string]()
	ss.Add(2, "two")
	ss.Add(1, "one")
	ss.Add(3, "three")

	assert.Equal(t, "one", ss.First())
	assert.Equal(t, "three", ss.Last())
}

func TestMinAndMax(t *testing.T) {
	ss := New[int, string]()
	ss.Add(2, "two")
	ss.Add(1, "one")
	ss.Add(3, "three")

	assert.Equal(t, "one", ss.Min())
	assert.Equal(t, "three", ss.Max())
}

func TestFirstKeyAndLastKey(t *testing.T) {
	ss := New[int, string]()
	ss.Add(2, "two")
	ss.Add(1, "one")
	ss.Add(3, "three")

	assert.Equal(t, 1, ss.FirstKey())
	assert.Equal(t, 3, ss.LastKey())
}

func TestMinKeyAndMaxKey(t *testing.T) {
	ss := New[int, string]()
	ss.Add(2, "two")
	ss.Add(1, "one")
	ss.Add(3, "three")

	assert.Equal(t, 1, ss.MinKey())
	assert.Equal(t, 3, ss.MaxKey())
}
