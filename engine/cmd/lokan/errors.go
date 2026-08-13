package main

import "fmt"

// NotFoundError is returned when no task matches the requested id.
type NotFoundError struct {
	ID string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("task not found: %s", e.ID)
}

// ValidationError is returned for invalid command input.
type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string {
	return e.msg
}

// cliErrorf builds a ValidationError with the given format string.
func cliErrorf(format string, args ...interface{}) error {
	return &ValidationError{msg: fmt.Sprintf(format, args...)}
}
