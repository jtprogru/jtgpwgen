package passgen

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
)

func secureIntn(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("secureIntn: n must be positive, got %d", n)
	}
	max := uint64(n)
	limit := ^uint64(0) - (^uint64(0) % max)
	var buf [8]byte
	for {
		if _, err := rand.Read(buf[:]); err != nil {
			return 0, err
		}
		v := binary.BigEndian.Uint64(buf[:])
		if v < limit {
			return int(v % max), nil
		}
	}
}

func pickRune(pool []rune) (rune, error) {
	idx, err := secureIntn(len(pool))
	if err != nil {
		return 0, err
	}
	return pool[idx], nil
}

func shuffleRunes(s []rune) error {
	for i := len(s) - 1; i > 0; i-- {
		j, err := secureIntn(i + 1)
		if err != nil {
			return err
		}
		s[i], s[j] = s[j], s[i]
	}
	return nil
}
