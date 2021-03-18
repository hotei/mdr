// mdr_strings.go (c) 2019 David Rook all rights reserved

package mdr

import (
	"strings"
)

// StringIsMember() returns true if s is a member of list ary.
// working
// test_?
func StringIsMember(s string, ary []string) bool {
	for i := 0; i < len(ary); i++ {
		if s == ary[i] {
			return true
		}
	}
	return false
}

// StringsUnion() returns a list which combines strings and returns the mininimum
// length result.  Same as StringsUnique(StringsCombine(a+b)) but a little more efficient.
// test_?

func StringsUnion(a, b []string) []string {
	var rv []string
	rv = a
	for i := 0; i < len(b); i++ {
		if StringIsMember(b[i], a) {
			continue
		}
		rv = append(rv, b[i])
	}
	return rv
}

// StringsCombine() returns a list which combines two lists of strings.
// test_?
func StringsCombine(a, b []string) []string {
	return append(a, b...)
}

// StringsSimilar() returns true if both string lists have the same members
// duplicate members are ignored.  To be SAME the member count must be equal.
// test_?
func StringsSimilar(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if StringIsMember(a[i], b) == false {
			return false
		}
	}
	return true
}

// StringsIntersection() returns a list where all members appear in both a,b.
// test_?
func StringsIntersection(a, b []string) []string {
	var rv []string

	for i := 0; i < len(b); i++ {
		if StringIsMember(b[i], a) {
			rv = append(rv, b[i])
		}
	}
	return rv
}

// StringsUnique returns a list with duplicate members removed.
// test_?
func StringsUnique(tags []string) (rv []string) {
	var myMap map[string]bool
	myMap = make(map[string]bool, 50)
	for _, val := range tags {
		//		fmt.Printf("val[%s],ndx[%d]\n", val, ndx)
		_, ok := myMap[val]
		//		fmt.Printf("ok %v, dat %v\n", ok, dat)
		if ok == false {
			myMap[val] = true
			rv = append(rv, val)
		}
	}
	//	for ndx, val := range rv {
	//		fmt.Printf("val[%s],ndx[%d]\n", val, ndx)
	//	}
	return rv
}

// RmPrefix() removes prefix from a string, if it exists
// test_?
func RmPrefix(s, pre string) string {
	if strings.HasPrefix(s, pre) {
		return s[len(pre):]
	}
	return s
}

// RmSuffix() removes suffix from a string, if it exists
// test_?
func RmSuffix(s, suf string) string {
	if strings.HasSuffix(s, suf) {
		return s[:len(s)-len(suf)]
	}
	return s
}

func ValidDecChar(c byte) bool {
	// test_?
	var decchars []byte = []byte("0123456789")
	for _, d := range decchars {
		if c == d {
			return true
		}
	}
	return false
}

func ValidDecString(s string) bool {
	// test_?
	for _, c := range s {
		if !ValidDecChar(byte(c)) {
			return false
		}
	}
	return true
}

func ValidHexChar(c byte) bool {
	// test_?
	var hexchars []byte = []byte("0123456789abcdefABCDEF")
	for _, h := range hexchars {
		if c == h {
			return true
		}
	}
	return false
}

func ValidHexString(s string) bool {
	// test_?
	for _, c := range s {
		if !ValidHexChar(byte(c)) {
			return false
		}
	}
	return true
}

// Reverse a slice of bytes in place

func Reverse(b []byte) {
	// Test with mdr_test.go:Test_006
	first := 0
	last := len(b) - 1
	for first < last {
		b[first], b[last] = b[last], b[first]
		first++
		last--
	}
}
