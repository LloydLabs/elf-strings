// Radix sort for []uint64.
package zuint64

import (
	"sort"
)

const (
	// Calling Sort() on slices smaller than this will result is sorting with sort.Sort() instead.
	MinSize      = 256
	radix   uint = 8
	bitSize uint = 64
)

// Sorts x using a Radix sort (Small slices are sorted with sort.Sort() instead).
func Sort(x []uint64) {
	if len(x) < MinSize {
		sort.Sort(uint64Sortable(x))
	} else {
		buffer := make([]uint64, len(x))
		SortBYOB(x, buffer)
	}
}

// Similar to Sort(), but returns a sorted copy of x, leaving x unmodified.
func SortCopy(x []uint64) []uint64 {
	y := make([]uint64, len(x))
	copy(y, x)
	Sort(y)
	return y
}

// Sorts x using a Radix sort, using supplied buffer space. Panics if
// len(x) is greater than len(buffer). Uses radix sort even on small slices.
func SortBYOB(x, buffer []uint64) {
	if len(x) > len(buffer) {
		panic("Buffer too small")
	}
	if len(x) < 2 {
		return
	}

	// Each pass processes a byte offset, copying back and forth between slices
	from := x
	to := buffer[:len(x)]
	var key uint8
	var offset [256]int // Keep track of where groups start

	for keyOffset := uint(0); keyOffset < bitSize; keyOffset += radix {
		keyMask := uint64(0xFF << keyOffset) // Current 'digit' to look at
		var counts [256]int                  // Keep track of the number of elements for each kind of byte
		sorted := true                       // Check for already sorted
		prev := uint64(0)                    // if elem is always >= prev it is already sorted
		for _, elem := range from {
			key = uint8((elem & keyMask) >> keyOffset) // fetch the byte at current 'digit'
			counts[key]++                              // count of elems to put in this digit's bucket
			if sorted {                                // Detect sorted
				sorted = elem >= prev
				prev = elem
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
		for _, elem := range from {
			key = uint8((elem & keyMask) >> keyOffset) // Get the digit
			to[offset[key]] = elem                     // Copy the element to the digit's bucket
			offset[key]++                              // One less space, move the offset
		}
		// On next pass copy data the other way
		to, from = from, to
	}
}

// Implements sort.Interface for small slices
type uint64Sortable []uint64

func (r uint64Sortable) Len() int           { return len(r) }
func (r uint64Sortable) Less(i, j int) bool { return r[i] < r[j] }
func (r uint64Sortable) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
