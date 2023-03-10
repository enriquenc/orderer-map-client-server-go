package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAndGet(t *testing.T) {
	m := NewOrderedMap()
	m.Add("key1", "value1")
	m.Add("key2", "value2")

	if value, ok := m.Get("key1"); !ok || value != "value1" {
		t.Errorf("m.Get(\"key1\") = (%v, %v), expected (%v, %v)", value, ok, "value1", true)
	}

	if value, ok := m.Get("key2"); !ok || value != "value2" {
		t.Errorf("m.Get(\"key2\") = (%v, %v), expected (%v, %v)", value, ok, "value2", true)
	}
}

func TestAddAndRemove(t *testing.T) {
	m := NewOrderedMap()
	m.Add("key1", "value1")
	m.Add("key2", "value2")

	m.Remove("key1")

	if _, ok := m.Get("key1"); ok {
		t.Errorf("m.Get(\"key1\") = (_, %v), expected (_, %v)", ok, false)
	}

	if value, ok := m.Get("key2"); !ok || value != "value2" {
		t.Errorf("m.Get(\"key2\") = (%v, %v), expected (%v, %v)", value, ok, "value2", true)
	}
}

func TestGetAll(t *testing.T) {
	m := NewOrderedMap()
	m.Add("key1", "value1")
	m.Add("key2", "value2")

	expected := []string{"key1=value1", "key2=value2"}
	if result := m.GetAll(); !equal(result, expected) {
		t.Errorf("m.GetAll() = %v, expected %v", result, expected)
	}
}

func TestOrderedMap_GetAll(t *testing.T) {
	// Initialize the map.
	m := NewOrderedMap()
	m.Add("a", "1")
	m.Add("b", "2")
	m.Add("c", "3")

	// Get all items and check the order.
	all := m.GetAll()
	assert.Equal(t, []string{"a=1", "b=2", "c=3"}, all)

	// Remove an item and check the order again.
	m.Remove("b")
	all = m.GetAll()
	assert.Equal(t, []string{"a=1", "c=3"}, all)

	// Add a new item and check the order again.
	m.Add("d", "4")
	all = m.GetAll()
	assert.Equal(t, []string{"a=1", "c=3", "d=4"}, all)
}

func TestOrderedMap_Add(t *testing.T) {
	// Initialize the map.
	m := NewOrderedMap()

	// Add some items and check the order.
	m.Add("a", "1")
	m.Add("b", "2")
	m.Add("c", "3")
	assert.Equal(t, []string{"a=1", "b=2", "c=3"}, m.GetAll())

	// Add an item that already exists and check the order.
	m.Add("a", "4")
	assert.Equal(t, []string{"a=4", "b=2", "c=3"}, m.GetAll())

	// Add another item and check the order.
	m.Add("d", "5")
	assert.Equal(t, []string{"a=4", "b=2", "c=3", "d=5"}, m.GetAll())
}

func TestOrderedMap_Remove(t *testing.T) {
	// Initialize the map.
	m := NewOrderedMap()
	m.Add("a", "1")
	m.Add("b", "2")
	m.Add("c", "3")
	m.Add("d", "4")

	// Remove an item from the middle and check the order.
	m.Remove("b")
	assert.Equal(t, []string{"a=1", "c=3", "d=4"}, m.GetAll())

	// Remove the first item and check the order.
	m.Remove("a")
	assert.Equal(t, []string{"c=3", "d=4"}, m.GetAll())

	// Remove the last item and check the order.
	m.Remove("d")
	assert.Equal(t, []string{"c=3"}, m.GetAll())

	// Remove an item that doesn't exist and check the order.
	m.Remove("e")
	assert.Equal(t, []string{"c=3"}, m.GetAll())
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
