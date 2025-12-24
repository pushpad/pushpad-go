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

// Int returns a pointer to an int value.
func Int(value int) *int {
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
