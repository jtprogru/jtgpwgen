package passgen

import (
	"regexp"
	"strings"
	"testing"
)

var memoRe = regexp.MustCompile(`^[a-z]{3}(-[a-z]{3})*[0-9]{2}@$`)

func TestGenerateMemoShape(t *testing.T) {
	for _, n := range []int{8, 16, 24, 32, 64} {
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
