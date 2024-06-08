package stringops

import (
	"fmt"
	"testing"
)

var givenTestInterfaces = []struct {
	given    interface{}
	expected string
	testName string
}{
	{
		given:    100,
		expected: "100",
		testName: "int_to_string",
	},
	{
		given:    "apple",
		expected: "apple",
		testName: "string_to_string",
	},
	{
		given:    "",
		expected: "",
		testName: "empty_string_to_string",
	},
	{
		given:    nil,
		expected: "",
		testName: "nil_to_string",
	},
	{
		given:    100.2,
		expected: "100.2",
		testName: "float_to_string",
	},
	{
		given:    0.2,
		expected: ".21",
		testName: "float-starting-with-point_to_string",
	},
	{
		given:    true,
		expected: "true",
		testName: "bool_to_string",
	},
}

func TestInterfaceToString(t *testing.T) {
	for _, tc := range givenTestInterfaces {
		t.Run(tc.testName, func(t *testing.T) {
			if actual := InterfaceToString(tc.given); tc.expected != actual {
				fmt.Printf("expected %s, actual %s", tc.expected, actual)
			}
		})
	}
}

func TestIsBlank(t *testing.T) {
	testCases := []struct {
		testName   string
		testString string
		expected   bool
	}{
		{
			testName:   "When string is blank",
			testString: "",
			expected:   true,
		},
		{
			testName:   "When string has spaces",
			testString: "        ",
			expected:   true,
		},
		{
			testName:   "When string has tabs",
			testString: "\t\n",
			expected:   true,
		},
		{
			testName:   "When string has special characters",
			testString: "!%#   ",
			expected:   false,
		},
		{
			testName:   "when string has spaces and is non-empty",
			testString: "   testme   ",
			expected:   false,
		},
		{
			testName:   "when string is not blank",
			testString: "testme",
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			if result := IsBlank(tc.testString); tc.expected != result {
				t.Errorf("expected %t, actual %t", tc.expected, result)
			}
		})
	}
}

func TestIsNilOrBlank(t *testing.T) {
	type S struct{ s string } // Helper type to assign pointer string value

	testCases := []struct {
		testName   string
		testString *string
		expected   bool
	}{
		{
			testName:   "When string is nil",
			testString: nil,
			expected:   true,
		},
		{
			testName:   "When string is blank",
			testString: &(&S{""}).s,
			expected:   true,
		},
		{
			testName:   "When string has spaces",
			testString: &(&S{"        "}).s,
			expected:   true,
		},
		{
			testName:   "When string has tabs",
			testString: &(&S{"\t\n"}).s,
			expected:   true,
		},
		{
			testName:   "When string has special characters",
			testString: &(&S{"!%#   "}).s,
			expected:   false,
		},
		{
			testName:   "when string has spaces and is non-empty",
			testString: &(&S{"   testme   "}).s,
			expected:   false,
		},
		{
			testName:   "when string is not blank",
			testString: &(&S{"testme"}).s,
			expected:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			if result := IsNilOrBlank(tc.testString); tc.expected != result {
				t.Errorf("expected %t, actual %t", tc.expected, result)
			}
		})
	}
}

func TestEqualFoldTrimSpace(t *testing.T) {
	testCases := []struct {
		testName    string
		testStringA string
		testStringB string
		expected    bool
	}{
		{
			testName:    "When both strings are blank",
			testStringA: "",
			testStringB: "",
			expected:    true,
		},
		{
			testName:    "When strings has spaces",
			testStringA: "        ",
			testStringB: "   ",
			expected:    true,
		},
		{
			testName:    "when strings has spaces and are non-empty",
			testStringA: "   testme   ",
			testStringB: "testme         ",
			expected:    true,
		},
		{
			testName:    "when strings are equal and not blank",
			testStringA: "testme",
			testStringB: "testme",
			expected:    true,
		},
		{
			testName:    "when string are not equal",
			testStringA: "testmehere",
			testStringB: "testmethere",
			expected:    false,
		},
		{
			testName:    "when strings have tabs/newlines but are equal",
			testStringA: "testme\t\n",
			testStringB: "\t\ttestme",
			expected:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			if result := EqualFoldTrimSpace(tc.testStringA, tc.testStringB); tc.expected != result {
				t.Errorf("expected %t, actual %t", tc.expected, result)
			}
		})
	}
}
