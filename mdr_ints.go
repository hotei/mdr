package mdr

import (
	"fmt"
	"log"
	"net"
)

type Ints []int

type IntPair struct {
	X, Y int
}

// =====================================

func MinI(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MaxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func InRangeI(a, b, c int) bool {
	if a > c {
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

func (a Ints) ContainsI(b int) bool {
	for _, val := range a {
		if val == b {
			return true
		}
	}
	return false
}

func ForceRangeI64(a, b, c int64) int64 {
	if a > c { // swap bounds if necessary to get a < b < c
		a, c = c, a
	}
	if b < a {
		return a
	}
	if b > c {
		return c
	}
	return b
}

func ConfineLoHi(lo, i, hi int) int {
	if i < lo {
		return lo
	}
	if i > hi {
		return hi
	}
	return i
}

// RotT2H rotates the head of an array to tail position
//	abcd => bcda
func (a Ints) RotH2T() Ints {
	ln := len(a)
	if ln <= 1 {
		return a
	}
	return append(a[1:], a[0])
}

//
// a b c  ln = 3
// 0 1 2
// c a b

// RotT2H rotates the tail of an array to head position
//	abcd => dabc
func (a Ints) RotT2H() Ints {
	ln := len(a)
	if ln <= 1 {
		return a
	}
	var x Ints = Ints{a[ln-1]}
	return append(x, a[:ln-1]...)
	//
	//return
}

// PermutedInts returns all the permutaions of the original array
// length of array must be in range of [0..6]
func PermutedInts(a Ints) []Ints {
	ln := len(a)
	if ln <= 1 {
		return []Ints{a}
	}
	if ln == 2 {
		return []Ints{{a[0], a[1]}, {a[1], a[0]}}
	}
	if ln > 6 {
		log.Panicf("permute array length > 6 (test phase only)\n")
	}
	expLen := Factorial(len(a))
	//fmt.Printf("PermutedInts should return %d elements each with %d elements\n", expLen, len(a))

	// make the empty return value array with zero length but full capacity
	var rv = make([]Ints, 0, Factorial(len(a)))
	for i := 0; i < ln; i++ {
		var x = a
		y := PermutedInts(x[1:]) // permute the tail
		for _, p := range y {
			x = append(Ints{a[0]}, p...)
			rv = append(rv, x)
		}
		a = a.RotH2T()
		//fmt.Printf("rotated a = %v\n", a)
	}
	if len(rv) != expLen {
		log.Panicf("return array has wrong length\n")
	}
	return rv
}

// Factorial computes value recursively.
// BUG(mdr) Factorial using recursion may not be the best choice ...
// BUG(mdr) check for overflow in Factorial or limit choice of n to good range
//  since it gets big quickly
func Factorial(n int) int {
	if n < 0 {
		return -1
	}
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}

// CreateBezierPts start, control, end pts for quadratic bezier.  segments is
// the number of line segments to create, more means smoother curve.
// Intended use is mostly with ring5 package, but MIGHT be genericly useful.
func CreateBezierPts(p1, p2, p3 Pointi, segments int) PointiList {
	//fmt.Printf("Bezier points start(%v) control(%v) end(%v)\n", p1, p2, p3)
	var pts = make(PointiList, segments)
	fx1, fy1 := float64(p1.X), float64(p1.Y)
	fx2, fy2 := float64(p2.X), float64(p2.Y)
	fx3, fy3 := float64(p3.X), float64(p3.Y)
	for i := 0; i < segments; i++ {
		c := float64(i) / float64(segments)
		a := 1 - c
		a, b, c := a*a, 2*c*a, c*c
		pts[i].X = int(a*fx1 + b*fx2 + c*fx3)
		pts[i].Y = int(a*fy1 + b*fy2 + c*fy3)
	}
	pts = append(pts, p3)
	return pts
}

/*
// used in image manipulation
func RangeMinMaxIntPair(v []IntPair) (minPt, maxPt IntPair) {
	vlen := len(v)
	if vlen <= 0 {
		// warn?
		return
	}
	minPt, maxPt = v[0], v[0]
	for i := 1; i < vlen; i++ {
		if v[i].X < minPt.X {
			minPt.X = v[i].X
		}
		if v[i].Y < minPt.Y {
			minPt.Y = v[i].Y
		}
		if v[i].X > maxPt.X {
			maxPt.X = v[i].X
		}
		if v[i].Y > maxPt.Y {
			maxPt.Y = v[i].Y
		}
	}
	return minPt, maxPt
}
*/

// LoHi returns the min and max of an array of ints
func LoHi(ary []int) (lo int, hi int) {
	if len(ary) <= 0 {
		log.Panic("empty array received by LoHi\n")
	}
	lo = ary[0]
	hi = ary[0]
	for i := 1; i < len(ary); i++ {
		if ary[i] < lo {
			lo = ary[i]
		}
		if ary[i] > hi {
			hi = ary[i]
		}
	}
	return lo, hi
}

// ==================   16 bit functions =====================

// SixteenBit converts from little endian two byte slice to int16
// aka Uint16FromLSBytes
func SixteenBit(n []byte) uint16 {
	if len(n) != 2 {
		FatalError(fmt.Errorf("mdr: Slice must be exactly 2 bytes\n"))
	}
	var rc uint16
	rc = uint16(n[1])
	rc <<= 8
	rc |= uint16(n[0])
	return rc
}

// ==================   32 bit functions =====================

// ThirtyTwoNet is a synonym for Uint32FromMSBytes
func ThirtyTwoNet(n []byte) uint32 {
	return Uint32FromMSBytes(n)
}

// convert from BIG endian four byte slice to int32
// reverse function is MSBytesFromUint32
func Uint32FromMSBytes(b []byte) uint32 {
	if len(b) != 4 {
		FatalError(fmt.Errorf("mdr: Slice must be exactly 4 bytes\n"))
	}
	var rc uint32
	rc = uint32(b[0])
	rc <<= 8
	rc |= uint32(b[1])
	rc <<= 8
	rc |= uint32(b[2])
	rc <<= 8
	rc |= uint32(b[3])
	return rc
}

// convert from LITTLE endian four byte slice to int32
// reverse function is LSBytesFromUint32
//  AKA Uint32FromLSBytes
func ThirtyTwoBit(n []byte) uint32 {
	if len(n) != 4 {
		FatalError(fmt.Errorf("mdr: Slice must be exactly 4 bytes\n"))
	}
	var rc uint32
	rc = uint32(n[3])
	rc <<= 8
	rc |= uint32(n[2])
	rc <<= 8
	rc |= uint32(n[1])
	rc <<= 8
	rc |= uint32(n[0])
	return rc
}

// Uint32FromIP returns a uint32 from IPv4 so we can use as index to map
//    -BEWARE- there are magic numbers that
//      presume knowledge of net.IP internals order
func Uint32FromIP(ip net.IP) uint32 {
	if false {
		for ndx, value := range ip {
			fmt.Printf("%d %d\n", ndx, value)
		}
	}
	var x uint32 = uint32(ip[12])
	x <<= 8
	x |= uint32(ip[13])
	x <<= 8
	x |= uint32(ip[14])
	x <<= 8
	x |= uint32(ip[15])
	return x
}

// IPFromUint32 beware - also presumes knowledge of net.IP internals order
//  Inverse of Uint32FromIP
func IPFromUint32(adr uint32) net.IP {
	d := byte(adr & 0xff)
	adr >>= 8
	c := byte(adr & 0xff)
	adr >>= 8
	b := byte(adr & 0xff)
	adr >>= 8
	a := byte(adr)
	return net.IPv4(a, b, c, d)
}

// convert uint32 to [0:4]byte slice in MSB first (aka 'Net') order
// reverse function is ThirtyTwoNet()
func MSBytesFromUint32(u uint32) []byte { //
	b := make([]byte, 4)
	b[3] = byte(u & 0xff)
	u >>= 8
	b[2] = byte(u & 0xff)
	u >>= 8
	b[1] = byte(u & 0xff)
	u >>= 8
	b[0] = byte(u & 0xff)
	return b
}

// convert uint32 to [0:4]byte slice in LSB first order
// reverse function is ThirtyTwoBit()
func LSBytesFromUint32(u uint32) []byte { //
	b := make([]byte, 4)
	b[0] = byte(u & 0xff)
	u >>= 8
	b[1] = byte(u & 0xff)
	u >>= 8
	b[2] = byte(u & 0xff)
	u >>= 8
	b[3] = byte(u & 0xff)
	return b
}

// ==================   64 bit functions below =====================

// true IFF a <= b <= c || a >= b >= c, note a < c not required
func InRangeInt64(a, b, c int64) bool {
	if a > c { // swap ends if necessary to get a < b < c
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

// CommaFmtInt64 returns a comma inserted number string
//  12345 becomes "12,345" NOTE! uses USA format, not internationalized
//  Name of a generic function might be DecimalFmtdInt64()
func CommaFmtInt64(n int64) string {
	//  Test is Test_007
	str := fmt.Sprintf("%d", n)
	nice := ""
	i := 1
	for n := len(str); n > 0; n-- {
		//fmt.Printf("%s\n", str[n-1:n])
		nice = nice + str[n-1:n]
		if (i % 3) == 0 {
			nice = nice + ","
		}
		//fmt.Printf("%s\n",nice)
		i++
	}
	niceNum := ""
	for n := len(nice); n > 0; n-- {
		niceNum = niceNum + nice[n-1:n]
	}
	if niceNum[0:1] == "," {
		niceNum = niceNum[1:]
	}
	return niceNum
}

// AbsInt64 returns the absolute value of an int64
func AbsInt64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

// LSBytesFromInt64 converts an int64 into an 8 byte array, LSB first (LittleEndian)
// test is Test_011
func LSBytesFromInt64(n int64) []byte {
	rv := make([]byte, 0, 8)
	for i := 0; i < 8; i++ {
		rv = append(rv, byte(n%256))
		n >>= 8
	}
	return rv
}

// MSBytesFromInt64 converts int64 to [0:8]byte slice with MSB first (BigEndian)
// reverse function is
func MSBytesFromInt64(n int64) []byte { //
	// test is Test_013
	b := make([]byte, 8)
	b[7] = byte(n & 0xff)
	n >>= 8
	b[6] = byte(n & 0xff)
	n >>= 8
	b[5] = byte(n & 0xff)
	n >>= 8
	b[4] = byte(n & 0xff)
	n >>= 8
	b[3] = byte(n & 0xff)
	n >>= 8
	b[2] = byte(n & 0xff)
	n >>= 8
	b[1] = byte(n & 0xff)
	n >>= 8
	b[0] = byte(n & 0xff)
	return b
}

// Int64FromMSBytes converts an 8 byte slice (BigEndian - MSB First) into an int64
func Int64FromMSBytes(b []byte) int64 {
	// see Test_014
	if len(b) != 8 {
		FatalError(fmt.Errorf("mdr: Slice must be exactly 8 bytes\n"))
	}
	rc := int64(0)
	rc = int64(b[0])
	rc <<= 8
	rc |= int64(b[1])
	rc <<= 8
	rc |= int64(b[2])
	rc <<= 8
	rc |= int64(b[3])
	rc <<= 8
	rc |= int64(b[4])
	rc <<= 8
	rc |= int64(b[5])
	rc <<= 8
	rc |= int64(b[6])
	rc <<= 8
	rc |= int64(b[7])
	return rc
}

// Int64FromLSBytes converts an 8 byte array into an int64
func Int64FromLSBytes(b []byte) int64 {
	// see Test_012
	if len(b) != 8 {
		FatalError(fmt.Errorf("mdr: Slice must be exactly 8 bytes\n"))
	}
	rv := int64(0)
	for i := 7; i >= 0; i-- {
		rv += int64(b[i])
		if i == 0 {
			break
		}
		rv <<= 8
	}
	return rv
}

// Smooth creates a new array with smoothed float64 values where n is number of
//   values over which to average.  Returns an array of same size as original,
//   so the first and last 'n' values will be "short", ie more sensitive to
//   any variation.
func Smooth(d []int64, n int) []float64 {
	//if n <= 1 {
	//return d
	//}
	fmt.Printf("Smoothing %d elements to avg[%d]\n", len(d), n)
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
