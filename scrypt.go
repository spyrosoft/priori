package main

import (
	"crypto/rand"

	scrypt "golang.org/x/crypto/scrypt"
)

func scryptHashAndSalt(plaintext string) (hash []byte, salt []byte, err error) {
	salt = generateByteSliceToken(10)
	hash, err = scryptHash(plaintext, salt)
	return
}

func scryptHash(plaintext string, salt []byte) (hash []byte, err error) {
	hashBytes, err := scrypt.Key([]byte(plaintext), salt, 16384, 8, 1, 32)
	hash = hashBytes
	return
}

func generateStringToken(length int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	const letterIdxBits = 6
	const letterIdxMask = 1<<letterIdxBits - 1
	result := make([]byte, length)
	bufferSize := int(float64(length) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < length; j++ {
		if j%bufferSize == 0 {
			randomBytes = generateByteSliceToken(bufferSize)
		}
		if idx := int(randomBytes[j%length] & letterIdxMask); idx < len(letterBytes) {
			result[i] = letterBytes[idx]
			i++
		}
	}
	return string(result)
}

func generateByteSliceToken(length int) []byte {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		notifyAdmin(err.Error())
	}
	return token
}
