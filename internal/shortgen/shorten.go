package shortgen

import (
	"crypto/rand"
	"math/big"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const defaultLength = 6

func Generate() (string, error) {
	result := make([]byte, defaultLength)
	alphabetLen := big.NewInt(int64(len(alphabet)))

	for i := 0; i < defaultLength; i++ {
		num, err := rand.Int(rand.Reader, alphabetLen)
		if err != nil {
			return "", err
		}
		result[i] = alphabet[num.Int64()]
	}

	return string(result), nil
}
