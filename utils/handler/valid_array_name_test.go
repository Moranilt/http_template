package handler

import "testing"

var nameTests = []struct {
	name     string
	field    string
	expected bool
	index    int
}{
	{
		name:     "valid array name",
		field:    "files[]",
		index:    5,
		expected: true,
	},
	{
		name:     "ends not with close bracket",
		field:    "files[]sad",
		index:    -1,
		expected: false,
	},
	{
		name:     "not array element",
		field:    "files",
		index:    -1,
		expected: false,
	},
	{
		name:     "empty element",
		field:    "",
		index:    -1,
		expected: false,
	},
	{
		name:     "not valid index element",
		field:    "files[ashds]",
		index:    5,
		expected: false,
	},
	{
		name:     "empty index",
		field:    "files[]",
		expected: true,
		index:    5,
	},
	{
		name:     "with index",
		field:    "files[1]",
		index:    5,
		expected: false,
	},
}

func TestValidArrayName(t *testing.T) {
	for _, test := range nameTests {
		t.Run(test.name, func(t *testing.T) {
			i, valid := isValidNameArray(test.field)

			if valid != test.expected {
				t.Errorf("got %t, expected %t", valid, test.expected)
			}

			if i != test.index {
				t.Errorf("got %d, expected %d", i, test.index)
			}
		})
	}
}

func BenchmarkValidArrayName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isValidNameArray("sahgdjhgajhds[]")
	}
}
