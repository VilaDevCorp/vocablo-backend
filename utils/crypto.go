package utils

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", nil
	}
	return string(hash), nil
}

func GenerateRandomToken(length int) (string, error) {
	b := make([]byte, length/2)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func CompareHash(hash string, pass string) bool {
	byteHash := []byte(hash)
	bytePass := []byte(pass)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePass)
	return err == nil
}
