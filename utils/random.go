package utils

import (
	"crypto/rand"
	"math/big"
)

type randomType struct{}

var Random = randomType{}

func (randomType) String(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	return randomString(length, letters)
}

func (randomType) Number(length int) string {
	const digits = "0123456789"
	return randomString(length, digits)
}

func randomString(length int, digits string) string {
	result := make([]byte, length)
	for i := range result {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return ""
		}
		result[i] = digits[index.Int64()]
	}
	return string(result)
}
