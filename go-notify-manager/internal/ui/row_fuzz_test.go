package ui

import (
	"testing"
)

func FuzzLinkify(f *testing.F) {
	seeds := []string{
		"Check https://example.com/test",
		"http://foo.bar/path",
		"No links here",
		"Special <xml> chars http://test.org",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		res := linkify(input)
		if len(input) > 0 && len(res) == 0 {
			t.Errorf("linkify(%q) returned empty string", input)
		}
	})
}

func FuzzResolveAppIconName(f *testing.F) {
	seeds := []string{
		"Discord",
		"Firefox Developer Edition",
		"Google-Chrome",
		"VSCode",
		"UnknownApp",
		"",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		icon := resolveAppIconName(input)
		if len(icon) == 0 {
			t.Errorf("resolveAppIconName(%q) returned empty string", input)
		}
	})
}
