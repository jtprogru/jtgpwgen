package passgen

import "strings"

var (
	memoConsonants = []rune("bcdfghjklmnpqrstvwxz")
	memoVowels     = []rune("aeiouy")
)

// generateMemo builds a syllable-based pronounceable password.
// Pattern: CVC syllables joined by '-', with a numeric+'@' suffix
// to satisfy common password policies. Output length is at least
// opts.Length; caller-visible length may exceed it slightly to
// preserve syllable boundaries (documented behavior).
func generateMemo(opts Options) (string, error) {
	target := opts.Length
	if target < 4 {
		target = 4
	}

	const suffixLen = 3 // 2 digits + '@'
	var b strings.Builder

	for b.Len() < target-suffixLen {
		if b.Len() > 0 {
			b.WriteRune('-')
		}
		syl, err := memoSyllable()
		if err != nil {
			return "", err
		}
		b.WriteString(syl)
	}

	for i := 0; i < 2; i++ {
		idx, err := secureIntn(len(digits))
		if err != nil {
			return "", err
		}
		b.WriteByte(digits[idx])
	}
	b.WriteRune('@')

	return b.String(), nil
}

func memoSyllable() (string, error) {
	c1, err := pickRune(memoConsonants)
	if err != nil {
		return "", err
	}
	v, err := pickRune(memoVowels)
	if err != nil {
		return "", err
	}
	c2, err := pickRune(memoConsonants)
	if err != nil {
		return "", err
	}
	return string([]rune{c1, v, c2}), nil
}
