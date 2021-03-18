// 999_test.go

// go test '-bench=.'   # expression matches everything, runs all benchmarks
// go test -run="Test_000" to run just one function

package mdr

import (
	//"fmt"
	"bytes"
	"hash/crc64"
	"io"
	"io/ioutil"
	"testing"
	//"time"
)

/////////////////////////  B E N C H M A R K S  ////////////////////////////
// go test -bench="*."   # re2 expression matches everything, runs all benchmarks
/*  4 GHz AMD-64  8120 8 core
Benchmark_PseudoRandomBlock-8	      50	  36,471,550 ns/op
Benchmark_FileLength-8	         1000000	       1,531 ns/op
Benchmark_BufSHA256-8	             200	  13,853,744 ns/op
Benchmark_BufMD5-8	                1000	   1,370,967 ns/op
Benchmark_BufCRC64-8	             500	   3,192,160 ns/op
*/

// 46.9e6 ns/op on 4Ghz AMD64 with 1.0.3
// 36.6e6 ns/op on 4Ghz AMD64 with 1.1 << 22% better >>
// 35.1e6 ns/op on 4Ghz AMD64 with 1.2
//  4.2e6 ns/op on 4Ghz i7    with 1.12.4
//  5.6e6 ns/op on merc       with 1.12.5
func Benchmark_PseudoRandomBlock(b *testing.B) {
	PRBsize := 1000000
	for i := 0; i < b.N; i++ {
		x := PseudoRandomBlock(PRBsize)
		r := bytes.NewReader(x)
		if _, err := io.Copy(ioutil.Discard, r); err != nil {
			panic(err)
		}
	}
}

// 1471 ns/op on 4Ghz AMD64 with 1.0.3
// 1341 ns/op on 4Ghz AMD64 with 1.1 << 8% better >>
// 1384 ns/op on 4Ghz AMD64 with 1.2
//  546 ns/op on merc       with 1.12.5
func Benchmark_FileLength(b *testing.B) {
	targetFile := "test-data/do_NOT_modify.txt"
	for i := 0; i < b.N; i++ {
		if _, err := FileLength(targetFile); err != nil {
			panic(err)
		}
	}
}

// 21.4e6 ns/op on 4Ghz AMD64 with 1.0.3
// 14.0e6 ns/op on 4Ghz AMD64 with 1.1  << 30% better >>
// 13.8e6 ns/op on 4Ghz AMD64 with 1.2
//  2.2e6 ns/op on 4Ghz i7    with 1.12.4
//  2.8e6 ns/op on merc i7    with 1.12.5
func Benchmark_BufSHA256(b *testing.B) {
	testBuf := PseudoRandomBlock(1024 * 1024)
	for i := 0; i < b.N; i++ {
		_ = BufSHA256(testBuf)
	}
}

//  1.5e6 ns/op on merc i7    with 1.12.5
func Benchmark_BufMD5(b *testing.B) {
	testBuf := PseudoRandomBlock(1024 * 1024)
	for i := 0; i < b.N; i++ {
		_ = BufMD5(testBuf)
	}
}

// 2.71e6 ns/op on 4Ghz AMD64 with 1.0.3
// 3.25e6 ns/op on 4Ghz AMD64 with 1.1  << 18% worse >>
// 3.19e6 ns/op on 4Ghz AMD64 with 1.2
// 0.49e6 ns/op on 4Ghz i7    with 1.12.4
// 0.66e6 ns/op on merc i7    with 1.12.5
func Benchmark_BufCRC64(b *testing.B) {
	testBuf := PseudoRandomBlock(1024 * 1024)
	for i := 0; i < b.N; i++ {
		_ = crc64.Checksum(testBuf[:], G_crcTable)
	}
}

// 117 ns/op on 4Ghz i7    with 1.12.4 (2.0e7 reps)
// 153 ns/op on merc i7    with 1.12.5
func Benchmark_KmBetweenGC(b *testing.B) {
	var (
		apt GPS2dT = GPS2dT{-76.0, 36.0}
		bpt GPS2dT = GPS2dT{-77.0, 37.0}
	)
	Verbose.Printf("%v %v %6.2f\n", apt, bpt, KmBetweenGC(apt, bpt))
	for i := 0; i < b.N; i++ {
		_ = KmBetweenGC(apt, bpt)
	}
}

func Looper(a, b GPS2dT) int {
	return 0
}

// 0.24 ns/op on 4Ghz i7    with 1.12.4 (2.0e9 reps) seems TOO fast?
// 0.29 ns/op on merc i7    with 1.12.5
func Benchmark_Looper(b *testing.B) {
	var (
		apt GPS2dT = GPS2dT{-76.0, 36.0}
		bpt GPS2dT = GPS2dT{-77.0, 37.0}
	)
	Verbose.Printf("%v %v %6.2f\n", apt, bpt, KmBetweenGC(apt, bpt))
	for i := 0; i < b.N; i++ {
		_ = Looper(apt, bpt)
	}
}
