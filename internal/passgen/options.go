package passgen

import (
	"errors"
	"fmt"
)

const (
	DefaultLength      = 24
	DefaultSpecialChar = "@"
	MaxLength          = 32768
)

type Options struct {
	Length       int
	UseLetters   bool
	UseDigits    bool
	UseSpecial   bool
	ExtraSpecial string
	Memo         bool
}

func DefaultOptions() Options {
	return Options{
		Length:     DefaultLength,
		UseLetters: true,
		UseDigits:  true,
		UseSpecial: true,
	}
}

var (
	ErrLengthOutOfRange = errors.New("length must be between 1 and 4096")
	ErrNoCharClasses    = errors.New("at least one character class must be enabled")
	ErrConflictDigits   = errors.New("--digits and --no-digits are mutually exclusive")
	ErrConflictSpecial  = errors.New("--special and --no-special are mutually exclusive")
)

func (o Options) Validate() error {
	if o.Length < 1 || o.Length > MaxLength {
		return fmt.Errorf("%w: got %d", ErrLengthOutOfRange, o.Length)
	}
	if o.Memo {
		return nil
	}
	if !o.UseLetters && !o.UseDigits && !o.UseSpecial {
		return ErrNoCharClasses
	}
	return nil
}
