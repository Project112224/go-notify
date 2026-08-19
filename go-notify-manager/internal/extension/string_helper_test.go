package extension

import (
	"testing"
)

func TestDateString_FormatToHMS(t *testing.T) {
	tests := []struct {
		input    DateString
		customFmt []string
		expected string
	}{
		{
			input:    DateString("2026-08-19T15:30:45Z"),
			expected: "15:30:45",
		},
		{
			input:     DateString("2026-08-19T15:30:45Z"),
			customFmt: []string{"2006-01-02"},
			expected:  "2026-08-19",
		},
		{
			input:    DateString("invalid-date-string"),
			expected: "invalid-date-string",
		},
	}

	for _, tt := range tests {
		res := tt.input.FormatToHMS(tt.customFmt...)
		if res != tt.expected {
			t.Errorf("DateString(%q).FormatToHMS() = %q; want %q", tt.input, res, tt.expected)
		}
	}
}
