package shortgen

import (
	"testing"
)

func TestGenerate_ReturnsCorrectLength(t *testing.T) {
	code, err := Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(code) != defaultLength {
		t.Errorf("expected length %d, got %d", defaultLength, len(code))
	}
}

func TestGenerate_OnlyAlphanumeric(t *testing.T) {
	code, err := Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, ch := range code {
		found := false
		for _, a := range alphabet {
			if ch == a {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("character '%c' is not in alphabet", ch)
		}
	}
}

func TestGenerate_Uniqueness(t *testing.T) {
	seen := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		code, err := Generate()
		if err != nil {
			t.Fatalf("unexpected error on iteration %d: %v", i, err)
		}
		if seen[code] {
			t.Errorf("duplicate code generated: %s", code)
		}
		seen[code] = true
	}
}
