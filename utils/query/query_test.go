package query

import "testing"

type testItem struct {
	name     string
	callback func(t *testing.T) string
	expected string
}

var tests = []testItem{
	{
		name: "default query",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "limit",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Limit("10")
			return query.String()
		},
		expected: "SELECT * FROM test_table LIMIT 10",
	},
	{
		name: "offset",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Offset("1")
			return query.String()
		},
		expected: "SELECT * FROM test_table OFFSET 1",
	},
	{
		name: "limit offset",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Limit("5").Offset("1")
			return query.String()
		},
		expected: "SELECT * FROM test_table LIMIT 5 OFFSET 1",
	},
	{
		name: "order by",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Order("name", DESC)
			return query.String()
		},
		expected: "SELECT * FROM test_table ORDER BY name DESC",
	},
	{
		name: "limit offset order by",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Order("name", DESC).Limit("5").Offset("1")
			return query.String()
		},
		expected: "SELECT * FROM test_table ORDER BY name DESC LIMIT 5 OFFSET 1",
	},
	{
		name: "where EQ",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().EQ("name", "testname")
			return query.String()
		},
		expected: "SELECT * FROM test_table WHERE name='testname'",
	},
	{
		name: "where LIKE",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().LIKE("name", "testname")
			return query.String()
		},
		expected: "SELECT * FROM test_table WHERE name LIKE '%testname%'",
	},
	{
		name: "where OR",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().OR(EQ("name", "testname"), EQ("age", "12"))
			return query.String()
		},
		expected: "SELECT * FROM test_table WHERE (name='testname' OR age='12')",
	},
	{
		name: "where OR AND",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().OR(EQ("name", "testname"), EQ("age", "12")).AND(EQ("id", "123"), EQ("email", "test@mail.com"))
			return query.String()
		},
		expected: "SELECT * FROM test_table WHERE (name='testname' OR age='12') AND (id='123' AND email='test@mail.com')",
	},
	{
		name: "not valid ORDER",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Order("test", "random")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "not valid LIMIT",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Limit("s")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "not valid OFFSET",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table").Offset("s")
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "not valid OR",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().OR(EQ("name", "testname"))
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "not valid AND",
		callback: func(t *testing.T) string {
			query := New("SELECT * FROM test_table")
			query.Where().AND(EQ("name", "testname"))
			return query.String()
		},
		expected: "SELECT * FROM test_table",
	},
	{
		name: "not valid DEFAULT AND",
		callback: func(t *testing.T) string {
			query := AND(EQ("name", "testname"))
			return query
		},
		expected: "",
	},
	{
		name: "not valid DEFAULT OR",
		callback: func(t *testing.T) string {
			query := OR(EQ("name", "testname"))
			return query
		},
		expected: "",
	},
}

func TestQuery(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.callback(t)
			if result != test.expected {
				t.Errorf("expected %q, got %q", test.expected, result)
			}
		})
	}
}
