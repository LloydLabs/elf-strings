// Zermelo is a library for sorting slices in Go.
package zermelo

import (
	"errors"
	"github.com/shawnsmithdev/zermelo/zfloat32"
	"github.com/shawnsmithdev/zermelo/zfloat64"
	"github.com/shawnsmithdev/zermelo/zint"
	"github.com/shawnsmithdev/zermelo/zint32"
	"github.com/shawnsmithdev/zermelo/zint64"
	"github.com/shawnsmithdev/zermelo/zuint"
	"github.com/shawnsmithdev/zermelo/zuint32"
	"github.com/shawnsmithdev/zermelo/zuint64"
	"sort"
)

// Attempts to sort x.
//
// If x is a supported slice type, this library will be used to sort it. Otherwise,
// if x implements sort.Interface it will passthrough to the sort.Sort() algorithm.
// Returns an error on unsupported types.
func Sort(x interface{}) error {
	switch xAsCase := x.(type) {
	case []float32:
		zfloat32.Sort(xAsCase)
	case []float64:
		zfloat64.Sort(xAsCase)
	case []int:
		zint.Sort(xAsCase)
	case []int32:
		zint32.Sort(xAsCase)
	case []int64:
		zint64.Sort(xAsCase)
	case []string:
		sort.Strings(xAsCase)
	case []uint:
		zuint.Sort(xAsCase)
	case []uint32:
		zuint32.Sort(xAsCase)
	case []uint64:
		zuint64.Sort(xAsCase)
	case sort.Interface:
		sort.Sort(xAsCase)
	default:
		return errors.New("type not supported")
	}
	return nil
}

// A Sorter can sort things like slices. Returns an error on unsupported types.
type Sorter interface {
	Sort(x interface{}) error
	CopySort(x interface{}) (interface{}, error)
}

// Reuseable buffers
type zSorter struct {
	bufFloat32 []float32
	bufFloat64 []float64
	bufInt     []int
	bufInt32   []int32
	bufInt64   []int64
	bufUint    []uint
	bufUint32  []uint32
	bufUint64  []uint64
}

// Given an existing buffer capacity and a requested one, finds a new buffer size.
// For the first alloc this will equal requested size, then after at it leaves
// a 25% buffer for future growth.
func allocSize(bufCap, reqLen int) int {
	if bufCap == 0 {
		return reqLen
	}
	return 5 * reqLen / 4
}

// Checks that buffers are large enough.
func (z *zSorter) prepBuffers(x interface{}) {
	switch xAsCase := x.(type) {
	case []float32:
		if cap(z.bufFloat32) < len(xAsCase) {
			z.bufFloat32 = make([]float32, allocSize(cap(z.bufFloat32), len(xAsCase)))
		}
	case []float64:
		if cap(z.bufFloat64) < len(xAsCase) {
			z.bufFloat64 = make([]float64, allocSize(cap(z.bufFloat64), len(xAsCase)))
		}
	case []int:
		if cap(z.bufInt) < len(xAsCase) {
			z.bufInt = make([]int, allocSize(cap(z.bufInt), len(xAsCase)))
		}
	case []int32:
		if cap(z.bufInt32) < len(xAsCase) {
			z.bufInt32 = make([]int32, allocSize(cap(z.bufInt32), len(xAsCase)))
		}
	case []int64:
		if cap(z.bufInt64) < len(xAsCase) {
			z.bufInt64 = make([]int64, allocSize(cap(z.bufInt64), len(xAsCase)))
		}
	case []uint:
		if cap(z.bufUint) < len(xAsCase) {
			z.bufUint = make([]uint, allocSize(cap(z.bufUint), len(xAsCase)))
		}
	case []uint32:
		if cap(z.bufUint32) < len(xAsCase) {
			z.bufUint32 = make([]uint32, allocSize(cap(z.bufUint32), len(xAsCase)))
		}
	case []uint64:
		if cap(z.bufUint64) < len(xAsCase) {
			z.bufUint64 = make([]uint64, allocSize(cap(z.bufUint64), len(xAsCase)))
		}
	}
}

func (z *zSorter) Sort(x interface{}) error {
	z.prepBuffers(x)
	switch xAsCase := x.(type) {
	case []float32:
		zfloat32.SortBYOB(xAsCase, z.bufFloat32[:len(xAsCase)])
	case []float64:
		zfloat64.SortBYOB(xAsCase, z.bufFloat64[:len(xAsCase)])
	case []int:
		zint.SortBYOB(xAsCase, z.bufInt[:len(xAsCase)])
	case []int32:
		zint32.SortBYOB(xAsCase, z.bufInt32[:len(xAsCase)])
	case []int64:
		zint64.SortBYOB(xAsCase, z.bufInt64[:len(xAsCase)])
	case []string:
		sort.Strings(xAsCase)
	case []uint:
		zuint.SortBYOB(xAsCase, z.bufUint[:len(xAsCase)])
	case []uint32:
		zuint32.SortBYOB(xAsCase, z.bufUint32[:len(xAsCase)])
	case []uint64:
		zuint64.SortBYOB(xAsCase, z.bufUint64[:len(xAsCase)])
	case sort.Interface:
		sort.Sort(xAsCase)
	default:
		return errors.New("type not supported")
	}
	return nil
}

func (z *zSorter) CopySort(x interface{}) (interface{}, error) {
	y := makeCopy(x)
	if y == nil {
		return x, errors.New("type not supported")
	}
	err := z.Sort(y)
	return y, err
}

func makeCopy(x interface{}) interface{} {
	switch xAsCase := x.(type) {
	case []float32:
		y := make([]float32, len(xAsCase))
		copy(y, xAsCase)
		return y
	case []float64:
		y := make([]float64, len(xAsCase))
		copy(y, xAsCase)
		return y
	case []int:
		y := make([]int, len(xAsCase))
		copy(y, xAsCase)
		return y
	case []int32:
		y := make([]int32, len(xAsCase))
		copy(y, xAsCase)
		return y
	case []int64:
		y := make([]int64, len(xAsCase))
		copy(y, xAsCase)
		return y
	case []string:
		y := make([]string, len(xAsCase))
		copy(y, xAsCase)
		return y
	case []uint:
		y := make([]uint, len(xAsCase))
		copy(y, xAsCase)
		return y
	case []uint32:
		y := make([]uint32, len(xAsCase))
		copy(y, xAsCase)
		return y
	case []uint64:
		y := make([]uint64, len(xAsCase))
		copy(y, xAsCase)
		return y
	default:
		return nil
	}
}

// Creates a Sorter that reuses buffers on repeated Sort() or CopySort() calls on the same type.
// This is not thread safe. CopySort() does not support passthrough of sort.Interface values.
func New() Sorter {
	return new(zSorter)
}
