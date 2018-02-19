// Radix sort for []float64.
package zfloat64

import (
	"math"
	"sort"
)

const (
	// Calling Sort() on slices smaller than this will result is sorting with sort.Sort() instead.
	MinSize         = 256
	radix      uint = 8
	radixShift uint = 3
	bitSize    uint = 64
)

// Sorts x using a Radix sort (Small slices are sorted with sort.Sort() instead).
func Sort(x []float64) {
	if len(x) < MinSize {
		sort.Float64s(x)
	} else {
		SortBYOB(x, make([]float64, len(x)))
	}
}

// Similar to Sort(), but returns a sorted copy of x, leaving x unmodified.
func SortCopy(x []float64) []float64 {
	y := make([]float64, len(x))
	copy(y, x)
	Sort(y)
	return y
}

// Sorts x using a Radix sort, using supplied buffer space. Panics if
// len(x) is greater than len(buffer). Uses radix sort even on small slices.
func SortBYOB(x, buffer []float64) {
	if len(x) > len(buffer) {
		panic("Buffer too small")
	}
	if len(x) < 2 {
		return
	}

	// Don't sort NaNs, just put them up front and skip them
	nans := 0
	for idx, val := range x {
		if math.IsNaN(val) {
			x[idx] = x[nans]
			x[nans] = val
			nans++
		}
	}

	// Each pass processes a byte offset, copying back and forth between slices
	from := x[nans:]
	to := buffer[:len(from)]
	var key uint8
	var uintVal uint64
	var offset [256]int // Keep track of where room is made for byte groups in the buffer

	for keyOffset := uint(0); keyOffset < bitSize; keyOffset += radix {
		keyMask := uint64(0xFF << keyOffset) // Current 'digit' to look at
		var counts [256]int                  // Keep track of the number of elements for each kind of byte
		sorted := true                       // Check for already sorted
		prev := float64(0)                   // if elem is always >= prev it is already sorted

		for _, val := range from {
			uintVal = floatFlip(math.Float64bits(val))
			key = uint8((uintVal & keyMask) >> keyOffset) // fetch the byte at current 'digit'
			counts[key]++                                 // count of values to put in this digit's bucket
			if sorted {                                   // Detect sorted
				sorted = val >= prev
				prev = val
			}
		}

		if sorted { // Short-circuit sorted
			if (keyOffset/radix)%2 == 1 {
				copy(to, from)
			}
			return
		}

		// Find target bucket offsets
		offset[0] = 0
		for i := 1; i < len(offset); i++ {
			offset[i] = offset[i-1] + counts[i-1]
		}

		// Rebucket while copying to other buffer
		for _, val := range from {
			uintVal = floatFlip(math.Float64bits(val))
			key = uint8((uintVal & keyMask) >> keyOffset) // Get the digit
			to[offset[key]] = val                         // Copy the element to the digit's bucket
			offset[key]++                                 // One less space, move the offset
		}
		// On next pass copy data the other way
		to, from = from, to
	}
}

// Converts a uint64 that represents a true float to one sorts properly
func floatFlip(x uint64) uint64 {
	if (x & 0x8000000000000000) == 0x8000000000000000 {
		return x ^ 0xFFFFFFFFFFFFFFFF
	}
	return x ^ 0x8000000000000000
}
