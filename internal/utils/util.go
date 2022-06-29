package utils

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	var output strings.Builder
	for i := 0; i < n; i++ {
		random := rand.Intn(len(letterBytes))
		randomChar := letterBytes[random]
		output.WriteString(string(randomChar))
	}
	return output.String()
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 15)
	return string(hash), err
}
