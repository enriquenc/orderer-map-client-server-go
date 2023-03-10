package orderedmap

import (
	"sync"
)

type node struct {
	key   string
	value string
	prev  *node
	next  *node
}

type OrderedMap struct {
	items map[string]*node
	head  *node
	tail  *node
	mu    sync.RWMutex
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		items: make(map[string]*node),
	}
}

func (m *OrderedMap) Add(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.items[key]; !exists {
		newNode := &node{
			key:   key,
			value: value,
			prev:  m.tail,
		}
		if m.tail != nil {
			m.tail.next = newNode
		} else {
			m.head = newNode
		}
		m.tail = newNode
		m.items[key] = newNode
	} else {
		m.items[key].value = value
	}
}

func (m *OrderedMap) Remove(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if node, exists := m.items[key]; exists {
		if node.prev != nil {
			node.prev.next = node.next
		} else {
			m.head = node.next
		}
		if node.next != nil {
			node.next.prev = node.prev
		} else {
			m.tail = node.prev
		}
		delete(m.items, key)
		return exists
	}
	return false
}

func (m *OrderedMap) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	node, exists := m.items[key]
	if exists {
		return node.value, true
	}
	return "", false
}

func (m *OrderedMap) GetAll() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]string, 0, len(m.items))
	for node := m.head; node != nil; node = node.next {
		result = append(result, node.key+"="+node.value)
	}
	return result
}
