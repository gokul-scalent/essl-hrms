package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"math/big"
	"strings"
)

func GenerateRandomPassword(length int) (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	password := make([]byte, length)

	for i := 0; i < length; i++ {
		// Generate a secure random index
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

const SALT = "G8sK8xNm0qAr2zY6" //"G8nL2aQpL8rU9vGn6TpQe7GhLk3Ve7StV2"

func DecodeBase64Password(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	decoded := string(decodedBytes)

	if !strings.HasPrefix(decoded, SALT) {
		return "", errors.New("invalid salt")
	}

	return strings.TrimPrefix(decoded, SALT), nil
}

// DecodeBase64Password decodes a base64 encoded password string and returns the original password.
func DecodePassword(encoded string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
