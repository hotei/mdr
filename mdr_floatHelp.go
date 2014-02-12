// mdr_floatHelp.go (c) 2012-2014 David Rook - see LICENSE.md

package mdr

// true IFF a <= b <= c || a >= b >= c, note a < c not required
func InRangeF(a, b, c float64) bool {
	if a > c { // swap bounds if necessary to get a < b < c
		a, c = c, a
	}
	if b < a {
		return false
	}
	if b > c {
		return false
	}
	return true
}

// Return the lo and hi of given array
func RangeLoHi64(v []float64) (lo, hi float64) {
	vlen := len(v)
	if vlen <= 0 {
		// warn?
		return
	}
	lo, hi = v[0], v[0]
	for i := 1; i < vlen; i++ {
		if v[i] < lo {
			lo = v[i]
		}
		if v[i] > hi {
			hi = v[i]
		}
	}
	return lo, hi
}
