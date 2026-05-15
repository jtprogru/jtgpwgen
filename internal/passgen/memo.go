package passgen

import (
	"math"
	"strings"
)

var (
	memoConsonants = []rune("bcdfghjklmnpqrstvwxz")
	memoVowels     = []rune("aeiouy")
)

const memoSuffixLen = 3 // "NN@"

// memoSyllableCount returns how many CVC syllables generateMemo will emit
// for the given requested target length.
func memoSyllableCount(target int) int {
	if target < 4 {
		target = 4
	}
	return (target + 1) / 4
}

// MemoEntropyBits returns the approximate Shannon entropy in bits for a
// memorable password of the requested length, given the current syllable
// alphabet and the fixed "NN@" suffix.
func MemoEntropyBits(length int) float64 {
	n := memoSyllableCount(length)
	syllableSpace := len(memoConsonants) * len(memoVowels) * len(memoConsonants)
	return float64(n)*math.Log2(float64(syllableSpace)) + math.Log2(100)
}

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

	var b strings.Builder

	for b.Len() < target-memoSuffixLen {
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
