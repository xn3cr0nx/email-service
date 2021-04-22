package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"

// RandomString returns a random generated string with fixed length
func RandomString() string {
	return RandomStringWithLen(20)
}

// RandomStringWithLen returns a random string with length
func RandomStringWithLen(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// RandomEmail returns a randomly generated email address
func RandomEmail() string {
	localPart := fmt.Sprintf("%s.%s", RandomStringWithLen(5), RandomStringWithLen(7))
	domain := RandomDomain()
	return strings.ToLower(fmt.Sprintf("%s@%s", localPart, domain))
}

// RandomDomain returns a randomly generated domain
func RandomDomain() string {
	domain := fmt.Sprintf("%s.%s", RandomStringWithLen(10), RandomStringWithLen(3))
	return strings.ToLower(domain)
}
