package mdr

import (
	"fmt"
	"log"
	"net"
)

// ==================   start with "int" functions =====================

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

// BUG(mdr) Factorial using recursion may not be the best choice ...
// BUG(mdr) check for overflow in Factorial or limit choice of n to good range
func Factorial(n int) int {
	if n < 0 {
		return -1
	}
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}

// ==================   next are the 32 bit functions =====================

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

// returns a uint32 from IPv4 so we can use as index to map
//    - beware - there are magic numbers here that
//    - presume knowledge of net.IP internals order
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

//  beware - presumes knowledge of net.IP internals order
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

// returns arg as comma inserted number string
//  12345 becomes "12,345" (not locale sensitive)
//  This is USA format, not internationalized
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

// returns the absolute value of an int64
func AbsInt64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

// convert an int64 into an 8 byte array, LSB first (LittleEndian)
// test is Test_011
func LSBytesFromInt64(n int64) []byte {
	rv := make([]byte, 0, 8)
	for i := 0; i < 8; i++ {
		rv = append(rv, byte(n%256))
		n >>= 8
	}
	return rv
}

// convert int64 to [0:8]byte slice with MSB first (BigEndian)
// reverse function is
// test is Test_013
func MSBytesFromInt64(n int64) []byte { //
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

// convert an 8 byte slice (BigEndian - MSB First) into an int64
// see Test_014
func Int64FromMSBytes(b []byte) int64 {
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

// convert an 8 byte array into an int64
// see Test_012
func Int64FromLSBytes(b []byte) int64 {
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
