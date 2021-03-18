// mdr_sliceI64.go (c) 2014-2019 David Rook - all rights reserved

package mdr

import (
//"fmt"
)

// ============================================================================  SliceI64
type SliceI64 []int64

func (d SliceI64) SliceI64Avg() float64 {
	l := len(d)
	if l <= 0 {
		return 0.0
	}
	if l == 1 {
		return float64(d[0])
	}
	sum := float64(0)
	for i := 0; i < l; i++ {
		sum += float64(d[i])
	}
	return sum / float64(l)
}

// create new array with smoothed values where n is number of values to average over
func (d SliceI64) Smooth(n int) []float64 {
	Verbose.Printf("Smoothing %d elements to avg[%d]\n", len(d), n)
	var (
		rv []float64
		x  SliceI64
	)
	rv = make([]float64, len(d))
	for i := 0; i < len(d); i++ {
		x = append(x, d[i])
		if len(x) > n {
			x = x[1:]
		}
		//fmt.Printf("i = %d len(x) = %d,len(rv) = %d\n", i,len(x),len(rv))
		rv[i] = x.SliceI64Avg()
	}
	return rv
}
