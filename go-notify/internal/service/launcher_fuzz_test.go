package service

import (
	"testing"
)

func FuzzExtractURLs(f *testing.F) {
	seeds := []string{
		"Visit https://example.com/test",
		"Check www.google.com for information",
		"Plain text with no URL",
		"http://domain.org/page?query=1&var=2#anchor",
		"<a href='https://foo.bar'>click</a>",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	svc := NewDesktopLauncherService()

	f.Fuzz(func(t *testing.T, input string) {
		urls := svc.ExtractURLs(input)
		for _, u := range urls {
			if len(u) == 0 {
				t.Errorf("Extracted empty URL from input %q", input)
			}
		}
	})
}

func FuzzCleanExecLine(f *testing.F) {
	seeds := []string{
		"firefox %u",
		"/usr/bin/code --new-window %F",
		"spotify %U %i %c",
		"simple_command --flag=value",
		"",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		_ = cleanExecLine(input)
	})
}
