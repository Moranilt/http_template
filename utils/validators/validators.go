package validators

import (
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// ValidDate checks if the string is a valid date in the format YYYY-MM-DD
func ValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

// ValidDateTime checks if the string is a valid date time in RFC3339 format
func ValidDateTime(date string) bool {
	_, err := time.Parse(time.RFC3339, date)

	return err == nil
}

// ValidInt checks if the string can be parsed as an integer
func ValidInt(val string) bool {
	_, err := strconv.Atoi(val)
	return err == nil
}

// ValidURL checks if the string is a valid URL that can be parsed
func ValidURL(uri string) bool {
	_, err := url.ParseRequestURI(uri)
	return err == nil
}

// ValidUUID checks if the string is a valid UUID
func ValidUUID(val string) bool {
	_, err := uuid.Parse(val)
	return err == nil
}
