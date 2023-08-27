package query

import (
	"fmt"
	"strings"

	"github.com/Moranilt/http_template/utils/validators"
)

const (
	// DESC is descending order
	DESC = "DESC"
	// ASC is ascending order
	ASC = "ASC"
)

// ValidOrderType checks if the given string is a valid ordering type
func ValidOrderType(val string) bool {
	val = strings.ToUpper(val)
	return val == DESC || val == ASC
}

type Query struct {
	// main is the base SQL query string
	main string

	// order is the ORDER BY clause
	order string

	// limit is the LIMIT clause
	limit string

	// offset is the OFFSET clause
	offset string

	// where is the WHERE clause
	where *Where
}

// New creates a new Query with the given base SQL query string
func New(q string) *Query {
	return &Query{
		main:  q,
		where: &Where{},
	}
}

// Where returns the WHERE clause for adding conditions
func (q *Query) Where() *Where {
	return q.where
}

// String returns the full SQL query string
func (q *Query) String() string {
	var where string
	if len(q.where.chunks) > 0 {
		where = "WHERE " + strings.Join(q.where.chunks, " AND ")
	}

	items := []string{where, q.order, q.limit, q.offset}
	var result strings.Builder
	result.WriteString(q.main)
	for _, item := range items {
		if len(item) != 0 {
			result.WriteString(" " + item)
		}
	}

	return result.String()
}

// Order adds an ORDER BY clause
func (q *Query) Order(by string, orderType string) *Query {
	if !ValidOrderType(orderType) {
		return q
	}
	q.order = fmt.Sprintf("ORDER BY %s %s", by, orderType)
	return q
}

// Limit adds a LIMIT clause
func (q *Query) Limit(val string) *Query {
	if !validators.ValidInt(val) {
		return q
	}
	q.limit = fmt.Sprintf("LIMIT %s", val)
	return q
}

// Offset adds an OFFSET clause
func (q *Query) Offset(val string) *Query {
	if !validators.ValidInt(val) {
		return q
	}
	q.offset = fmt.Sprintf("OFFSET %s", val)
	return q
}

// The Where struct represents the WHERE clause in a SQL query.
// It contains a chunks slice to hold the individual WHERE conditions.

type Where struct {
	chunks []string
}

// EQ adds an equality condition to the WHERE clause.
func (w *Where) EQ(fieldName string, value string) *Where {
	w.chunks = append(w.chunks, EQ(fieldName, value))
	return w
}

// LIKE adds a LIKE condition to the WHERE clause.
func (w *Where) LIKE(fieldName string, value string) *Where {
	w.chunks = append(w.chunks, LIKE(fieldName, value))
	return w
}

// OR adds an OR condition to the WHERE clause.
func (w *Where) OR(args ...string) *Where {
	if len(args) < 2 {
		return w
	}
	w.chunks = append(w.chunks, OR(args...))
	return w
}

// AND adds an AND condition to the WHERE clause.
func (w *Where) AND(args ...string) *Where {
	if len(args) < 2 {
		return w
	}
	w.chunks = append(w.chunks, AND(args...))
	return w
}

// AND creates an AND condition string from the provided arguments.
func AND(args ...string) string {
	if len(args) < 2 {
		return ""
	}
	return fmt.Sprintf("(%s)", strings.Join(args, " AND "))
}

// EQ creates an equality condition string.
func EQ(fieldName string, value string) string {
	return fmt.Sprintf("%s='%s'", fieldName, value)
}

// LIKE creates a LIKE condition string.
func LIKE(fieldName string, value string) string {
	return fmt.Sprintf("%s LIKE '%%%s%%'", fieldName, value)
}

// OR creates an OR condition string from the provided arguments.
func OR(args ...string) string {
	if len(args) < 2 {
		return ""
	}
	return fmt.Sprintf("(%s)", strings.Join(args, " OR "))
}
