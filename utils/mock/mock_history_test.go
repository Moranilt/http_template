package mock

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMockHistory(t *testing.T) {
	t.Run("valid calls", func(t *testing.T) {
		var validTests = []struct {
			name string
			data []byte
			err  error
		}{
			{
				name: "operation_1",
				data: []byte{'a', 'a'},
				err:  nil,
			},
			{
				name: "operation_2",
				data: []byte{'a', 'b'},
				err:  nil,
			},
			{
				name: "operation_3",
				data: []byte{'a', 'c'},
				err:  nil,
			},
			{
				name: "operation_4",
				data: []byte{'a', 'd'},
				err:  nil,
			},
		}
		history := NewMockHistory[[]byte]()
		for _, test := range validTests {
			history.Push(test.name, test.data, test.err)
		}

		for _, test := range validTests {
			item, err := history.Get(test.name)
			if err != nil {
				t.Errorf("name: %q. got error: %v", test.name, err)
				return
			}

			if !reflect.DeepEqual(item.Data, test.data) {
				t.Errorf("name: %q. not expected data. Expected %v, got %v", test.name, test.data, item.Data)
				return
			}

			if item.Err != test.err {
				t.Errorf("name %q. not expected error. Expected %v, got %v", test.name, test.err, item.Err)
				return
			}
		}

		if len(history.events) != 0 {
			t.Errorf("events is not empty. Value %v", history.events)
		}

		for name, value := range history.storage {
			if len(value) != 0 {
				t.Errorf("name %q. Value is not empty. Got %v", name, value)
				return
			}
		}
	})

	t.Run("item is not expected to be called", func(t *testing.T) {
		var validTests = []struct {
			name string
			data []byte
			err  error
		}{
			{
				name: "operation_1",
				data: []byte{'a', 'a'},
				err:  nil,
			},
			{
				name: "operation_2",
				data: []byte{'a', 'b'},
				err:  nil,
			},
			{
				name: "operation_3",
				data: []byte{'a', 'c'},
				err:  nil,
			},
			{
				name: "operation_4",
				data: []byte{'a', 'd'},
				err:  nil,
			},
		}
		history := NewMockHistory[[]byte]()
		for _, test := range validTests {
			history.Push(test.name, test.data, test.err)
		}

		item, err := history.Get("invalid_operation")
		if err != nil {
			expectedErr := fmt.Errorf(ERR_Not_Expected_Call, "invalid_operation", validTests[0].name)
			if err.Error() != expectedErr.Error() {
				t.Errorf("expected error %q, got %q", expectedErr, err)
			}
			return
		}

		if item != nil {
			t.Errorf("item should be nil for invalid call, got %v", item)
		}
	})

	t.Run("empty events expected", func(t *testing.T) {
		history := NewMockHistory[[]byte]()
		item, err := history.Get("unexpected_operation")
		if err != nil {
			expectedErr := fmt.Errorf(ERR_Events_Is_Empty, "unexpected_operation")
			if err.Error() != expectedErr.Error() {
				t.Errorf("expected error %q, got %q", expectedErr, err)
			}
			return
		}

		if item != nil {
			t.Errorf("item should be nil for invalid call, got %v", item)
		}
	})

	t.Run("empty getNextItem", func(t *testing.T) {
		history := NewMockHistory[byte]()
		name, item := history.getNextItem()
		if name != nil {
			t.Errorf("name should be nil, got %q", *name)
		}

		if item != nil {
			t.Errorf("item data should be nil, got %v", *item)
		}
	})

	t.Run("all expectations were not done", func(t *testing.T) {
		history := NewMockHistory[[]byte]()
		history.Push("ExpectedEvent", nil, nil)

		if err := history.AllExpectationsDone(); err == nil {
			t.Errorf("expected error, but got nil")
		}
	})

	t.Run("clear history", func(t *testing.T) {
		history := NewMockHistory[[]byte]()
		history.Push("event_1", []byte("data 1"), nil)
		history.Push("event_2", []byte("data 2"), nil)

		history.Clear()

		if len(history.events) != 0 {
			t.Errorf("events array should be empty, got %#v", history.events)
		}

		if len(history.storage) != 0 {
			t.Errorf("storage should be empty, got %#v", history.storage)
		}
	})

	t.Run("is empty history", func(t *testing.T) {
		history := NewMockHistory[[]byte]()
		history.Push("event_1", []byte("data 1"), nil)
		history.Push("event_2", []byte("data 2"), nil)

		if history.IsEmpty() {
			t.Errorf("expected history to be not empty")
		}

		history.Clear()

		if !history.IsEmpty() {
			t.Error("expected history to be empty")
		}

	})
}
