package generate

import (
	"crypto/rand"
	"math/big"
)

func SecureToken(n int) (string, error) {
	const from = "abcdefghijklmnopqrstuvwxyz0123456789"

	max := big.NewInt(int64(len(from)))
	res := make([]byte, n)

	for i := range res {
		posBig, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}

		pos := posBig.Int64()

		res[i] = from[pos]
	}

	return string(res), nil
}
