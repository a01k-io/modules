package arrayops

import "strings"

// ContainsInt return true/false based on element found
func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ContainsString return true/false based on  element found
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// UniqueIntArray creates an array of int with unique values.
func UniqueIntArray(a []int) []int {
	length := len(a)
	seen := make(map[int]struct{}, length)
	j := 0
	for i := 0; i < length; i++ {
		v := a[i]
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		a[j] = v
		j++
	}
	return a[0:j]
}

// UniqueStringArray creates an array of string with unique values.
func UniqueStringArray(a []string) []string {
	length := len(a)
	seen := make(map[string]struct{}, length)
	j := 0
	for i := 0; i < length; i++ {
		v := a[i]
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		a[j] = v
		j++
	}
	return a[0:j]
}

// SubtractString remove elements from first array passed
func SubtractString(x []string, y []string) []string {
	if len(x) == 0 {
		return []string{}
	}
	if len(y) == 0 {
		return x
	}
	slice := []string{}
	hash := map[string]struct{}{}

	for _, v := range x {
		hash[v] = struct{}{}
	}

	for _, v := range y {
		_, ok := hash[v]
		if ok {
			delete(hash, v)
		}
	}

	for _, v := range x {
		_, ok := hash[v]
		if ok {
			slice = append(slice, v)
		}
	}
	return slice
}

// Intersection it will find the intersection of two slices
func Intersection(a, b []string) []string {
	var c = make([]string, 0)
	m := make(map[string]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return c
}

// StrSliceToLowerCase converts all elements of a string array to lowercase
func StrSliceToLowerCase(s []string) []string {
	for i, v := range s {
		s[i] = strings.ToLower(v)
	}
	return s
}
