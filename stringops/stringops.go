package stringops

import (
	"fmt"
	"strings"
)

// NormalizeSpace multiple white spaces between word to single space
// Also remove leading and trailing spaces
func NormalizeSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// InterfaceToString converts the value of interface to string representation according to https://golang.org/pkg/fmt/ for types:
//
// `bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, complex64, complex128, string`.
//
// It returns empty("") string for types other than defined above.
//
// the value in a default format.
// follow this link https://golang.org/pkg/fmt/
// for more information.
func InterfaceToString(i interface{}) string {
	switch i.(type) {
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, complex64, complex128, string:
		return fmt.Sprintf("%v", i)
	default:
		return ""
	}
}

// IsBlank returns true if the string 's' is blank.
func IsBlank(s string) bool {
	if s = strings.TrimSpace(s); s == "" {
		return true
	}
	return false
}

// IsNilOrBlank return true of s is nil or *s in blank.
func IsNilOrBlank(s *string) bool {
	if s == nil {
		return true
	}

	return IsBlank(*s)
}

// EqualFoldTrimSpace reports whether space trimed s and t, interpreted as UTF-8 strings, are equal under Unicode case-folding.
func EqualFoldTrimSpace(s, t string) bool {
	return strings.EqualFold(strings.TrimSpace(s), strings.TrimSpace(t))
}
