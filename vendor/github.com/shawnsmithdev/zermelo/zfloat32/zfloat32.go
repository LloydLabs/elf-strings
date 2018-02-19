// Radix sort for []float32.
package zfloat32

import (
	"math"
	"sort"
)

const (
	// Calling Sort() on slices smaller than this will result is sorting with sort.Sort() instead.
	MinSize         = 256
	radix      uint = 8
	radixShift uint = 3
	bitSize    uint = 32
)

// Sorts x using a Radix sort (Small slices are sorted with sort.Sort() instead).
func Sort(x []float32) {
	if len(x) < MinSize {
		sort.Sort(float32Sortable(x))
	} else {
		SortBYOB(x, make([]float32, len(x)))
	}
}

// Similar to Sort(), but returns a sorted copy of x, leaving x unmodified.
func SortCopy(x []float32) []float32 {
	y := make([]float32, len(x))
	copy(y, x)
	Sort(y)
	return y
}

// Sorts x using a Radix sort, using supplied buffer space. Panics if
// len(x) is greater than len(buffer). Uses radix sort even on small slices.
func SortBYOB(x, buffer []float32) {
	if len(x) > len(buffer) {
		panic("Buffer too small")
	}
	if len(x) < 2 {
		return
	}

	// Don't sort NaNs, just put them up front and skip them
	nans := 0
	for idx, val := range x {
		if math.IsNaN(float64(val)) {
			x[idx] = x[nans]
			x[nans] = val
			nans++
		}
	}

	// Each pass processes a byte offset, copying back and forth between slices
	from := x[nans:]
	to := buffer[:len(from)]
	var key uint8
	var uintVal uint32
	var offset [256]int // Keep track of where room is made for byte groups in the buffer

	for keyOffset := uint(0); keyOffset < bitSize; keyOffset += radix {
		keyMask := uint32(0xFF << keyOffset) // Current 'digit' to look at
		var counts [256]int                  // Keep track of the number of elements for each kind of byte
		sorted := true                       // Check for already sorted
		prev := float32(0)                   // if elem is always >= prev it is already sorted

		for _, val := range from {
			uintVal = floatFlip(math.Float32bits(val))
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
			uintVal = floatFlip(math.Float32bits(val))
			key = uint8((uintVal & keyMask) >> keyOffset) // Get the digit
			to[offset[key]] = val                         // Copy the element to the digit's bucket
			offset[key]++                                 // One less space, move the offset
		}
		// On next pass copy data the other way
		to, from = from, to
	}
}

// Converts a uint32 that represents a true float to one sorts properly
func floatFlip(x uint32) uint32 {
	if (x & 0x80000000) == 0x80000000 {
		return x ^ 0xFFFFFFFF
	}
	return x ^ 0x80000000
}

type float32Sortable []float32

func (r float32Sortable) Len() int           { return len(r) }
func (r float32Sortable) Less(i, j int) bool { return r[i] < r[j] }
func (r float32Sortable) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
