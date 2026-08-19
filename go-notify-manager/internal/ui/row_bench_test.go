package ui

import (
	"testing"
)

func BenchmarkLinkify(b *testing.B) {
	text := "Check out https://github.com/june/go-notify for details and http://example.com"

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = linkify(text)
	}
}

func BenchmarkResolveAppIconName(b *testing.B) {
	appNames := []string{"Discord", "Firefox Developer Edition", "Google-Chrome", "VSCode", "Alacritty"}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = resolveAppIconName(appNames[i%len(appNames)])
	}
}
