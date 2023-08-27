package validators

import (
	"testing"
)

type validatorTest struct {
	name     string
	arg      string
	expected bool
}

var urlTests = []validatorTest{
	{
		name:     "valid url",
		arg:      "http://test.com",
		expected: true,
	},
	{
		name:     "not valid url",
		arg:      "?test=name",
		expected: false,
	},
}

func TestValidURL(t *testing.T) {
	runValidatorTest(t, ValidURL, urlTests)
}

var intTests = []validatorTest{
	{
		name:     "valid int",
		arg:      "100",
		expected: true,
	},
	{
		name:     "int with letters",
		arg:      "10O",
		expected: false,
	},
	{
		name:     "int with special symbols",
		arg:      "10☐",
		expected: false,
	},
	{
		name:     "int with spaces",
		arg:      "1 0",
		expected: false,
	},
}

func TestValidInt(t *testing.T) {
	runValidatorTest(t, ValidInt, intTests)
}

var dateTimeTests = []validatorTest{
	{
		name:     "valid date time",
		arg:      "2023-08-15T09:29:49.047Z",
		expected: true,
	},
	{
		name:     "not valid date time",
		arg:      "2023-01-02",
		expected: false,
	},
	{
		name:     "empty value",
		arg:      "",
		expected: false,
	},
}

func TestValidDateTime(t *testing.T) {
	runValidatorTest(t, ValidDateTime, dateTimeTests)
}

var dateTests = []validatorTest{
	{
		name:     "valid date",
		arg:      "2023-08-15",
		expected: true,
	},
	{
		name:     "not valid date",
		arg:      "2023-08-15T09:29:49.047Z",
		expected: false,
	},
	{
		name:     "empty value",
		arg:      "",
		expected: false,
	},
}

func TestValidDate(t *testing.T) {
	runValidatorTest(t, ValidDate, dateTests)
}

var uuidTests = []validatorTest{
	{
		name:     "valid uuid",
		arg:      "9aef6831-8538-4b03-8f61-7a687861584a",
		expected: true,
	},
	{
		name:     "not valid uuid",
		arg:      "9aef6831",
		expected: false,
	},
	{
		name:     "empty value",
		arg:      "",
		expected: false,
	},
}

func TestValidUUID(t *testing.T) {
	runValidatorTest(t, ValidUUID, uuidTests)
}

func runValidatorTest(t *testing.T, caller func(val string) bool, tests []validatorTest) {
	t.Helper()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			valid := caller(test.arg)
			if valid != test.expected {
				t.Errorf("expected %t, got %t", test.expected, valid)
			}
		})
	}
}
