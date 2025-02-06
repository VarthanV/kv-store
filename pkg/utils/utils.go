package utils

func Insert[T any](slice []T, index int, value T) []T {
	if index < 0 || index > len(slice) {
		return slice
	}

	// Grow the slice by one element.
	slice = append(slice, value)

	// Use copy to shift the elements after the index to the right.
	copy(slice[index+1:], slice[index:])

	// Insert the new value at the specified index.
	slice[index] = value

	return slice
}
