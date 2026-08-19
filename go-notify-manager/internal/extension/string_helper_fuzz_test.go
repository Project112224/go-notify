package extension

import (
	"testing"
)

func FuzzDateFormatToHMS(f *testing.F) {
	seeds := []string{
		"2026-08-19T15:30:45Z",
		"2026-08-19T15:30:45+08:00",
		"invalid-date",
		"",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		ds := DateString(input)
		_ = ds.FormatToHMS()
		_ = ds.FormatToHMS("2006-01-02")
	})
}
