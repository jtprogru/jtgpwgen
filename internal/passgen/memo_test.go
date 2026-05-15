package passgen

import (
	"errors"
	"math"
	"regexp"
	"strings"
	"testing"
)

var memoRe = regexp.MustCompile(`^[a-z]{3}(-[a-z]{3})*[0-9]{2}@$`)

func TestGenerateMemoShape(t *testing.T) {
	for _, n := range []int{24, 32, 64} {
		opts := Options{Length: n, Memo: true}
		got, err := Generate(opts)
		if err != nil {
			t.Fatalf("len=%d err: %v", n, err)
		}
		if !memoRe.MatchString(got) {
			t.Fatalf("memo password %q does not match expected shape", got)
		}
		if len(got) < n {
			t.Fatalf("memo length %d < requested %d (%q)", len(got), n, got)
		}
		if !strings.HasSuffix(got, "@") {
			t.Fatalf("missing trailing @: %q", got)
		}
	}
}

func TestGenerateMemoEntropyTooLow(t *testing.T) {
	for _, n := range []int{4, 8, 16, 22} {
		_, err := Generate(Options{Length: n, Memo: true})
		if !errors.Is(err, ErrMemoEntropyTooLow) {
			t.Fatalf("len=%d: want ErrMemoEntropyTooLow, got %v", n, err)
		}
	}
}

func TestMemoEntropyBits(t *testing.T) {
	cases := []struct {
		length    int
		syllables int
	}{
		{4, 1},
		{8, 2},
		{16, 4},
		{20, 5},
		{23, 6},
		{24, 6},
		{32, 8},
	}
	syllableBits := math.Log2(float64(len(memoConsonants) * len(memoVowels) * len(memoConsonants)))
	for _, c := range cases {
		got := MemoEntropyBits(c.length)
		want := float64(c.syllables)*syllableBits + math.Log2(100)
		if math.Abs(got-want) > 1e-9 {
			t.Fatalf("len=%d: want %.4f, got %.4f", c.length, want, got)
		}
	}
}
