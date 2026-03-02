package repository

import (
	"database/sql"
	"time"
)

// toNullString converts a string pointer to sql.NullString.
func toNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *s, Valid: true}
}

// toNullInt64 converts an int pointer to sql.NullInt64.
func toNullInt64(i *int) sql.NullInt64 {
	if i == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: int64(*i), Valid: true}
}

// toNullFloat64 converts a float64 pointer to sql.NullFloat64.
func toNullFloat64(f *float64) sql.NullFloat64 {
	if f == nil {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{Float64: *f, Valid: true}
}

// toNullTime converts a time.Time pointer to sql.NullTime.
func toNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// toNullBool converts a bool pointer to sql.NullBool.
func toNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

// stringPtr creates a pointer to a string.
func stringPtr(s string) *string {
	return &s
}

// intPtr creates a pointer to an int.
func intPtr(i int) *int {
	return &i
}

// float64Ptr creates a pointer to a float64.
func float64Ptr(f float64) *float64 {
	return &f
}

// boolPtr creates a pointer to a bool.
func boolPtr(b bool) *bool {
	return &b
}

// timePtr creates a pointer to a time.Time.
func timePtr(t time.Time) *time.Time {
	return &t
}
