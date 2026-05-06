package passgen

import (
	"strings"
	"testing"
	"unicode"
)

func TestGenerateDefaults(t *testing.T) {
	opts := DefaultOptions()
	for i := 0; i < 50; i++ {
		got, err := Generate(opts)
		if err != nil {
			t.Fatalf("Generate err: %v", err)
		}
		if len([]rune(got)) != DefaultLength {
			t.Fatalf("len = %d, want %d (got %q)", len([]rune(got)), DefaultLength, got)
		}
		if !containsLetter(got) {
			t.Fatalf("missing letter: %q", got)
		}
		if !containsDigit(got) {
			t.Fatalf("missing digit: %q", got)
		}
		if !strings.ContainsRune(got, '@') {
			t.Fatalf("missing default special @: %q", got)
		}
	}
}

func TestGenerateNoSpecialNoDigits(t *testing.T) {
	opts := DefaultOptions()
	opts.UseDigits = false
	opts.UseSpecial = false
	opts.Length = 32
	for i := 0; i < 50; i++ {
		got, err := Generate(opts)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range got {
			if !unicode.IsLetter(r) {
				t.Fatalf("unexpected non-letter %q in %q", r, got)
			}
		}
	}
}

func TestGenerateExtraSpecial(t *testing.T) {
	opts := DefaultOptions()
	opts.ExtraSpecial = "#$%"
	opts.Length = 64
	got, err := Generate(opts)
	if err != nil {
		t.Fatal(err)
	}
	allowedSpecial := "@#$%"
	for _, r := range got {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			continue
		}
		if !strings.ContainsRune(allowedSpecial, r) {
			t.Fatalf("rune %q not in allowed set %q (full: %q)", r, allowedSpecial, got)
		}
	}
}

func TestGenerateLengthTooSmall(t *testing.T) {
	opts := Options{Length: 1, UseLetters: true, UseDigits: true, UseSpecial: true}
	if _, err := Generate(opts); err == nil {
		t.Fatal("expected error for length smaller than class count")
	}
}

func TestGenerateValidationErr(t *testing.T) {
	if _, err := Generate(Options{Length: 0}); err == nil {
		t.Fatal("expected validation error")
	}
}

func containsLetter(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func containsDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}
