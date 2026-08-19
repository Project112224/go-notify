package service

import (
	"testing"
)

func BenchmarkExtractURLs(b *testing.B) {
	svc := NewDesktopLauncherService()
	text := "Visit https://github.com/june/go-notify and check www.google.com for updates!"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = svc.ExtractURLs(text)
	}
}

func BenchmarkCleanExecLine(b *testing.B) {
	line := "/usr/bin/code --new-window %F %u --user-data-dir %i"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = cleanExecLine(line)
	}
}
