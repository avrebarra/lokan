package id

import "strconv"

// GenerateID returns the plain counter as the task id, e.g. "2".
func GenerateID(counter int) string {
	return strconv.Itoa(counter)
}
