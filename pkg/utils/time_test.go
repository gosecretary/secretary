package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration string
		expected time.Duration
	}{
		{
			name:     "seconds",
			duration: "30s",
			expected: 30 * time.Second,
		},
		{
			name:     "minutes",
			duration: "5m",
			expected: 5 * time.Minute,
		},
		{
			name:     "hours",
			duration: "2h",
			expected: 2 * time.Hour,
		},
		{
			name:     "mixed duration",
			duration: "1h30m",
			expected: 1*time.Hour + 30*time.Minute,
		},
		{
			name:     "invalid duration",
			duration: "invalid",
			expected: 0,
		},
		{
			name:     "empty duration",
			duration: "",
			expected: 0,
		},
		{
			name:     "microseconds",
			duration: "500us",
			expected: 500 * time.Microsecond,
		},
		{
			name:     "milliseconds",
			duration: "100ms",
			expected: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseDuration(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}
