package internal

import (
	"math"
	"unicode"
)

// Percentage calculates the percentage amount
func Percentage(amount float64, rate float64) float64 {
	return amount * rate / 100
}

// Round rounds a float to 2 decimal places
func Round(val float64) float64 {
	return math.Round(val*100) / 100
}

// ArrayMap maps a slice using a function
func ArrayMap[T any, U any](data []T, f func(T) U) []U {
	res := make([]U, len(data))
	for i, e := range data {
		res[i] = f(e)
	}
	return res
}

// ToLower converts a rune to lowercase
func ToLower(r rune) rune {
	return unicode.ToLower(r)
}
