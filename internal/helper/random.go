package helper

import (
	"crypto/rand"
	"math/big"
)

const (
	letterCharset       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	passwordCharset     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}<>?"
)

func GenerateLetter(length int) string {
	return generateString(length, letterCharset)
}

func GenerateString(length int) string {
	return generateString(length, alphanumericCharset)
}

func GeneratePassword(length int) string {
	return generateString(length, passwordCharset)
}

func generateString(length int, charset string) string {
	if length <= 0 {
		return ""
	}

	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := range result {
		n, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			panic("helper: failed to generate random string: " + err.Error())
		}

		result[i] = charset[n.Int64()]
	}

	return string(result)
}
