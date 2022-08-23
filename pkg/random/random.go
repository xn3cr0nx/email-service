package random

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const digits = "0123456789"

// Time generates a random time.Time, which has its UNIX timestamp
// set in [startSec,endSec). It panics if endSec <= startSec.
func Time(startSec, endSec int64) time.Time {
	if endSec <= startSec {
		log.Panic("endSec must be higher than startSec")
	}
	sec := rand.Int63n(endSec-startSec) + startSec
	return time.Unix(sec, 0).UTC()
}

// Int returns a random int64 in range between min and max
func Int(min, max int) int {
	if min == 0 && max == 0 {
		return 0
	}

	if min == max {
		return max
	}

	return rand.Intn(max-min) + min
}

// IntRange returns a random int64 in range between min and max
func IntRange(min, max int64) int64 {
	return rand.Int63n(max-min) + min
}

// FloatRange returns a random float64 in range between min and max
func FloatRange(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

// String returns a random generated string with fixed length
func String() string {
	return StringWithLen(20)
}

// StringWithLen returns a random string with length
func StringWithLen(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// Password returns a random valid password
func Password() string {
	return String() + DigitsWithLen(2)
}

// MultilineString creates a string with nLines lines, each of length lineLength
func MultilineString(lineLength int, nLines int) string {
	lines := make([]string, 0, nLines)
	for i := 0; i < nLines; i++ {
		lines = append(lines, StringWithLen(lineLength))
	}
	return strings.Join(lines, "\\n")
}

// DigitsWithLen returns a random string of digits with length
func DigitsWithLen(length int) string {
	result := make([]byte, length)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}

// Email returns a randomly generated email address
func Email() string {
	localPart := fmt.Sprintf("%s.%s", StringWithLen(5), StringWithLen(7))
	domain := Domain()
	return strings.ToLower(fmt.Sprintf("%s@%s", localPart, domain))
}

// Domain returns a randomly generated domain
func Domain() string {
	domain := fmt.Sprintf("%s.%s", StringWithLen(10), StringWithLen(3))
	return strings.ToLower(domain)
}

func URL() string {
	return strings.ToLower(fmt.Sprintf("https://%s/%s", Domain(), StringWithLen(7)))
}
