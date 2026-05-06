package passgen

const (
	lettersLower = "abcdefghijklmnopqrstuvwxyz"
	lettersUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits       = "0123456789"
)

func uniqueRunes(s string) []rune {
	seen := make(map[rune]struct{}, len(s))
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if _, ok := seen[r]; ok {
			continue
		}
		seen[r] = struct{}{}
		out = append(out, r)
	}
	return out
}
