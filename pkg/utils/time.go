package utils

import "time"

// ParseDuration parses a duration string into a time.Duration value.
func ParseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return 0
	}
	return d
}
