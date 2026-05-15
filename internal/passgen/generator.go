package passgen

import "fmt"

func Generate(opts Options) (string, error) {
	if err := opts.Validate(); err != nil {
		return "", err
	}
	if opts.Memo {
		return generateMemo(opts)
	}

	var classes [][]rune
	if opts.UseLetters {
		classes = append(classes, []rune(lettersLower+lettersUpper))
	}
	if opts.UseDigits {
		classes = append(classes, []rune(digits))
	}
	if opts.UseSpecial {
		classes = append(classes, uniqueRunes(DefaultSpecialChar+opts.ExtraSpecial))
	}

	if len(classes) == 0 {
		return "", ErrNoCharClasses
	}
	if opts.Length < len(classes) {
		return "", fmt.Errorf("length %d is smaller than number of required classes %d", opts.Length, len(classes))
	}

	pool := make([]rune, 0)
	for _, c := range classes {
		pool = append(pool, c...)
	}

	out := make([]rune, 0, opts.Length)
	for _, c := range classes {
		r, err := pickRune(c)
		if err != nil {
			return "", err
		}
		out = append(out, r)
	}
	for len(out) < opts.Length {
		r, err := pickRune(pool)
		if err != nil {
			return "", err
		}
		out = append(out, r)
	}
	if err := shuffleRunes(out); err != nil {
		return "", err
	}
	s := string(out)
	zeroRunes(out)
	zeroRunes(pool)
	return s, nil
}

func zeroRunes(s []rune) {
	for i := range s {
		s[i] = 0
	}
}
