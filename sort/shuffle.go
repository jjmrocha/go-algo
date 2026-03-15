package sort

import (
	"crypto/rand"
	"math/big"
)

// Shuffle randomly permutes arr in place using a cryptographically secure random source.
// It returns an error only if the underlying random source fails.
func Shuffle[T any](arr []T) error {
	for i := 1; i < len(arr); i++ {
		j, err := nextInt(i + 1)
		if err != nil {
			return err
		}

		Swap(arr, i, j)
	}

	return nil
}

func nextInt(n int) (int, error) {
	i, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0, err
	}

	return int(i.Int64()), nil
}
