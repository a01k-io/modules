package rand

import (
	"math/rand"
	"time"
)

const (
	// DefaultLength when passwed as an argument, it generates otp having 6-digit
	DefaultLength    = 6
	numbers          = "0123456789"
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Type describes different types of otp
type Type string

const (
	//Numberic defines numeric only type
	Numberic Type = "NumbericOnly"

	//LowercaseLetters defines lower case only
	LowercaseLetters Type = "LowercaseLetters"
	//UppercaseLetters defines upper case only
	UppercaseLetters Type = "UppercaseLetters"
	//MixedcaseLetters defines mixed case only
	MixedcaseLetters Type = "MixedcaseLetters"

	//AlphaNumericWithLowercaseLetters defines lower case only
	AlphaNumericWithLowercaseLetters Type = "AlphaNumericWithLowercaseLetters"
	//AlphaNumericWithUppercaseLetters defines upper case only
	AlphaNumericWithUppercaseLetters Type = "AlphaNumericWithUppercaseLetters"
	//AlphaNumericWithMixedCaseLetters defines mixed case only
	AlphaNumericWithMixedCaseLetters Type = "AlphaNumericWithMixedCaseLetters"
)

func getTokenStore(t Type) string {
	store := ""
	switch t {
	case Numberic:
		store = numbers
	case LowercaseLetters:
		store = lowercaseLetters
	case UppercaseLetters:
		store = uppercaseLetters
	case MixedcaseLetters:
		store = lowercaseLetters + uppercaseLetters
	case AlphaNumericWithLowercaseLetters:
		store = numbers + lowercaseLetters
	case AlphaNumericWithUppercaseLetters:
		store = numbers + uppercaseLetters
	case AlphaNumericWithMixedCaseLetters:
		store = numbers + lowercaseLetters + uppercaseLetters
	}
	return store
}

func getRandString(input string, length int) string {
	b := make([]byte, length)
	for i := range b {
		n := rand.Intn(len(input))
		b[i] = input[n]
	}
	return string(b)
}

// GenerateString generates random string based on given type and length
func GenerateString(t Type, length int) string {
	tokenStore := getTokenStore(t)
	return getRandString(tokenStore, length)
}
