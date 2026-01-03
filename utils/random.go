package utils

import "math/rand"

// RandomSelect returns a random element from a slice
// Returns the zero value of T if the slice is empty
func RandomSelect[T any](items []T) T {
	if len(items) == 0 {
		var zero T
		return zero
	}
	return items[rand.Intn(len(items))]
}
