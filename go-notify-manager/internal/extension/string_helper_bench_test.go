package extension

import (
	"testing"
)

func BenchmarkFormatToHMS(b *testing.B) {
	ds := DateString("2026-08-19T15:30:45Z")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ds.FormatToHMS()
	}
}
