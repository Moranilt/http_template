package handler

import "testing"

var tests = []struct {
	name     string
	field    string
	expected string
	valid    bool
}{
	{
		name:     "array with index",
		field:    "files[0]",
		expected: "files[0]",
		valid:    false,
	},
	{
		name:     "array",
		field:    "files[]",
		expected: "files",
		valid:    true,
	},
	{
		name:     "ends not with close bracket",
		field:    "files[0]sad",
		expected: "files[0]sad",
		valid:    false,
	},
	{
		name:     "not array element",
		field:    "files",
		expected: "files",
		valid:    false,
	},
	{
		name:     "empty element",
		field:    "",
		expected: "",
		valid:    false,
	},
}

func TestExtractName(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			name, valid := extractArrayName(test.field)

			if name != test.expected {
				t.Errorf("expected %q, got %q", test.expected, name)
			}
			if valid != test.valid {
				t.Errorf("expected %t, got %t", test.valid, valid)
			}
		})
	}
}

func BenchmarkExtractName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		extractArrayName("sahgdjhgajhds[]")
	}
}
