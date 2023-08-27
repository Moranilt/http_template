package mock

import (
	"fmt"
	"sync"
)

const (
	ERR_Not_Expected_Call = "not expected call of %q. Expected %q"
	ERR_Events_Is_Empty   = "not expected call of %q. Events is empty"
)

type MockHistoryItem[T any] struct {
	Data T
	Err  error
}

// MockHistory stores the history of mock calls.
type MockHistory[T any] struct {
	events  []string
	storage map[string][]MockHistoryItem[T]
	mu      sync.RWMutex
}

// NewMockHistory creates a new MockHistory.
func NewMockHistory[T any]() *MockHistory[T] {
	return &MockHistory[T]{
		storage: make(map[string][]MockHistoryItem[T]),
	}
}

// Push adds a new mock call to the history.
func (m *MockHistory[T]) Push(name string, data T, err error) {
	m.mu.Lock()
	m.events = append(m.events, name)
	m.storage[name] = append(m.storage[name], MockHistoryItem[T]{
		Data: data,
		Err:  err,
	})
	m.mu.Unlock()
}

// Get retrieves a mock call from the history by name.
func (m *MockHistory[T]) Get(name string) (*MockHistoryItem[T], error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.events) != 0 {
		if m.events[0] == name {
			_, firstItem := m.getNextItem()
			return firstItem, nil
		} else {
			return nil, fmt.Errorf(ERR_Not_Expected_Call, name, m.events[0])
		}
	} else {
		return nil, fmt.Errorf(ERR_Events_Is_Empty, name)
	}
}

// Clear resets the history.
func (m *MockHistory[T]) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = []string{}
	m.storage = make(map[string][]MockHistoryItem[T])
}

// AllExpectationsDone checks if all expected calls were made.
func (m *MockHistory[T]) AllExpectationsDone() error {
	if len(m.events) != 0 {
		name, data := m.getNextItem()
		return fmt.Errorf("expected %q to be called with %#v", *name, *data)
	}

	return nil
}

// IsEmpty checks if the history is empty.
func (m *MockHistory[T]) IsEmpty() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.events) == 0 && len(m.storage) == 0
}

// getNextItem returns the next mock call details from history.
func (m *MockHistory[T]) getNextItem() (*string, *MockHistoryItem[T]) {
	if len(m.events) != 0 {
		eventName := m.events[0]
		eventData := m.storage[eventName]
		item := eventData[0]
		m.events = m.events[1:]
		m.storage[eventName] = m.storage[eventName][1:]
		return &eventName, &item
	}

	return nil, nil
}
