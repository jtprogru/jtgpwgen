package passgen

import (
	"errors"
	"testing"
)

func TestOptionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr error
	}{
		{"defaults ok", DefaultOptions(), nil},
		{"zero length", Options{Length: 0, UseLetters: true}, ErrLengthOutOfRange},
		{"too long", Options{Length: MaxLength + 1, UseLetters: true}, ErrLengthOutOfRange},
		{"all classes off", Options{Length: 8}, ErrNoCharClasses},
		{"memo bypasses class check", Options{Length: 16, Memo: true}, nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.opts.Validate()
			if tc.wantErr == nil {
				if err != nil {
					t.Fatalf("unexpected err: %v", err)
				}
				return
			}
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("want %v, got %v", tc.wantErr, err)
			}
		})
	}
}
