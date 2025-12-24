package pushpad

import "time"

// String returns a pointer to a string value.
func String(value string) *string {
	return &value
}

// Bool returns a pointer to a bool value.
func Bool(value bool) *bool {
	return &value
}

// Int64 returns a pointer to an int64 value.
func Int64(value int64) *int64 {
	return &value
}

// Time returns a pointer to a time value.
func Time(value time.Time) *time.Time {
	return &value
}

// Strings returns a pointer to a string slice.
func Strings(value []string) *[]string {
	return &value
}
