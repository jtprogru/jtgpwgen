package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/jtprogru/jtgpwgen/internal/passgen"
)

func runRoot(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	cmd := newRootCmd()
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return stdout.String(), err
}

func TestRunE_MemoConflictsWithClassFlags(t *testing.T) {
	cases := [][]string{
		{"--memo", "--no-special"},
		{"--memo", "--no-digits"},
		{"--memo", "--special", "#"},
		{"--memo", "--digits"},
	}
	for _, args := range cases {
		_, err := runRoot(t, args...)
		if !errors.Is(err, passgen.ErrMemoIncompatibleFlag) {
			t.Fatalf("args=%v: want ErrMemoIncompatibleFlag, got %v", args, err)
		}
	}
}

func TestRunE_MemoTooShortReturnsError(t *testing.T) {
	_, err := runRoot(t, "--memo", "--length", "16")
	if !errors.Is(err, passgen.ErrMemoEntropyTooLow) {
		t.Fatalf("want ErrMemoEntropyTooLow, got %v", err)
	}
}

func TestRunE_DefaultsHappyPath(t *testing.T) {
	out, err := runRoot(t, "--length", "24")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	out = strings.TrimRight(out, "\n")
	if len(out) != 24 {
		t.Fatalf("expected length 24, got %d (%q)", len(out), out)
	}
}
