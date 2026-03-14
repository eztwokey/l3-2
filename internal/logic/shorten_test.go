package logic

import (
	"testing"

	"github.com/eztwokey/l3-shortener/internal/models"
)

func TestCreateLink_EmptyURL(t *testing.T) {
	l := &Logic{}

	_, err := l.CreateLink(t.Context(), models.CreateLinkRequest{
		URL: "",
	})
	if err == nil {
		t.Fatal("expected error for empty URL, got nil")
	}
}

func TestCreateLink_InvalidURL(t *testing.T) {
	l := &Logic{}

	cases := []struct {
		name string
		url  string
	}{
		{"no scheme", "example.com"},
		{"javascript", "javascript:alert(1)"},
		{"ftp", "ftp://files.example.com"},
		{"empty scheme", "://example.com"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := l.CreateLink(t.Context(), models.CreateLinkRequest{
				URL: tc.url,
			})
			if err == nil {
				t.Errorf("expected error for URL %q, got nil", tc.url)
			}
		})
	}
}

func TestCreateLink_CustomCodeTooShort(t *testing.T) {
	l := &Logic{}

	_, err := l.CreateLink(t.Context(), models.CreateLinkRequest{
		URL:        "https://example.com",
		CustomCode: "ab",
	})
	if err == nil {
		t.Fatal("expected error for short custom code, got nil")
	}
}

func TestCreateLink_CustomCodeTooLong(t *testing.T) {
	l := &Logic{}

	_, err := l.CreateLink(t.Context(), models.CreateLinkRequest{
		URL:        "https://example.com",
		CustomCode: "abcdefghijklmnopqrstuvwxyz",
	})
	if err == nil {
		t.Fatal("expected error for long custom code, got nil")
	}
}
