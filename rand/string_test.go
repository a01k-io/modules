package rand

import (
	"math/rand"
	"regexp"
	"testing"
	"time"
)

func isValidStringType(otp string, t Type) bool {
	switch t {
	case Numberic:
		return regexp.MustCompile(`^[0-9]+$`).MatchString(otp)
	case LowercaseLetters:
		return regexp.MustCompile(`^[a-z]+$`).MatchString(otp)
	case UppercaseLetters:
		return regexp.MustCompile(`^[A-Z]+$`).MatchString(otp)
	case MixedcaseLetters:
		return regexp.MustCompile(`^[A-Za-z]+$`).MatchString(otp)
	case AlphaNumericWithLowercaseLetters:
		return regexp.MustCompile(`^[0-9a-z]+$`).MatchString(otp)
	case AlphaNumericWithUppercaseLetters:
		return regexp.MustCompile(`^[0-9A-Z]+$`).MatchString(otp)
	case AlphaNumericWithMixedCaseLetters:
		return regexp.MustCompile(`^[0-9A-Za-z]+$`).MatchString(otp)
	}
	return false
}

var generateRandomStringCases = []struct {
	name         string
	stringType   Type
	stringLength int
}{
	{
		stringType:   Numberic,
		stringLength: 4,
		name:         "numeric otp only",
	},
	{
		stringType:   LowercaseLetters,
		stringLength: 5,
		name:         "lowercase otp only",
	},
	{
		stringType:   Numberic,
		stringLength: 6,
		name:         "uppercase letters otp only",
	},
	{
		stringType:   Numberic,
		stringLength: 7,
		name:         "mixed case letters otp",
	},
	{
		stringType:   Numberic,
		stringLength: 8,
		name:         "alpha numeric with lowercase letters",
	},
	{
		stringType:   Numberic,
		stringLength: 9,
		name:         "alpha numeric with uppercase letters",
	},
	{
		stringType:   Numberic,
		stringLength: 10,
		name:         "alpha numeric with mixed case letters",
	},
}

func TestGenerateString(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	for _, c := range generateRandomStringCases {
		t.Run(c.name, func(t *testing.T) {
			got := GenerateString(c.stringType, c.stringLength)
			if c.stringLength != len(got) {
				t.Errorf("expected length[%d], actual[%d]", c.stringLength, len(got))
			} else if !isValidStringType(got, c.stringType) {
				t.Errorf("expected of type[%s], actual otp[%s]", c.stringType, got)
			}
		})
	}
}
