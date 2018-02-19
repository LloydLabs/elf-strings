// Radix sort for []int64.
package zint64

import (
	"sort"
)

const (
	// Calling Sort() on slices smaller than this will result is sorting with sort.Sort() instead.
	MinSize        = 256
	radix    uint  = 8
	bitSize  uint  = 64
	minInt64 int64 = -1 << 63
)

// Sorts x using a Radix sort (Small slices are sorted with sort.Sort() instead).
func Sort(x []int64) {
	if len(x) < MinSize {
		sort.Sort(int64Sortable(x))
	} else {
		buffer := make([]int64, len(x))
		SortBYOB(x, buffer)
	}
}

// Similar to Sort(), but returns a sorted copy of x, leaving x unmodified.
func SortCopy(x []int64) []int64 {
	y := make([]int64, len(x))
	copy(y, x)
	Sort(y)
	return y
}

// Sorts a []int64 using a Radix sort, using supplied buffer space. Panics if
// len(x) does not equal len(buffer). Uses radix sort even on small slices.
func SortBYOB(x, buffer []int64) {
	if len(x) > len(buffer) {
		panic("Buffer too small")
	}
	if len(x) < 2 {
		return
	}

	from := x
	to := buffer[:len(x)]
	var key uint8
	var offset [256]int // Keep track of where groups start

	for keyOffset := uint(0); keyOffset < bitSize; keyOffset += radix {
		keyMask := int64(0xFF << keyOffset)
		var counts [256]int // Keep track of the number of elements for each kind of byte
		sorted := true
		prev := minInt64

		for _, elem := range from {
			// For each elem to sort, fetch the byte at current radix
			key = uint8((elem & keyMask) >> keyOffset)
			// inc count of bytes of this type
			counts[key]++
			if sorted { // Detect sorted
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

		if keyOffset == bitSize-radix {
			// Handle signed values
			// Count negative elements (last 128 counts)
			negCnt := 0
			for i := 128; i < 256; i++ {
				negCnt += counts[i]
			}

			offset[0] = negCnt // Start of positives
			offset[128] = 0    // Start of negatives
			for i := 1; i < 128; i++ {
				// Positive values
				offset[i] = offset[i-1] + counts[i-1]
				// Negative values
				offset[i+128] = offset[i+127] + counts[i+127]
			}
		} else {
			offset[0] = 0
			for i := 1; i < len(offset); i++ {
				offset[i] = offset[i-1] + counts[i-1]
			}
		}

		// Swap values between the buffers by radix
		for _, elem := range from {
			key = uint8((elem & keyMask) >> keyOffset) // Get the byte of each element at the radix
			to[offset[key]] = elem                     // Copy the element depending on byte offsets
			offset[key]++
		}
		from, to = to, from
	}
}

type int64Sortable []int64

func (r int64Sortable) Len() int           { return len(r) }
func (r int64Sortable) Less(i, j int) bool { return r[i] < r[j] }
func (r int64Sortable) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
