// mdr_test.go

// go test -bench="*."   # re2 expression matches everything, runs all benchmarks

package mdr

import (
	"bytes"
	"fmt"
	"hash/crc64"
	"io"
	"io/ioutil"
	"testing"
	"time"
)

// Template, call t.Errorf with reason for failure
// BUG(mdr): what's diff between Errorf and Fatalf ?
func Test_000(t *testing.T) {
	fmt.Printf("Test_000 \n")
	if false {
		t.Errorf("print fail, but keep testing")
	}
	if false {
		t.Fatalf("print fail and keep testing")
	}
	fmt.Printf("go test -bench=\"*.\" to run all benchmarks\n")
	fmt.Printf("Pass - test 000\n")
}

func Test_001(t *testing.T) {
	Verbose = false
	fmt.Printf("\nTest_001 Jobsplit\n")
	x := JobSplit(10, 3)
	fmt.Printf("%v\n", x)
	goodsplit := true
	if len(x) != 3 {
		goodsplit = false
	}
	var y = IntPair{0, 3}
	if x[0] != y {
		goodsplit = false
	}
	y = IntPair{4, 6}
	if x[1] != y {
		goodsplit = false
	}
	y = IntPair{7, 9}
	if x[2] != y {
		goodsplit = false
	}

	// fails if x[2] != IntPair{7,9} {goodsplit = false }

	if !goodsplit {
		t.Errorf("split failed\n")
	} else {
		fmt.Printf("Pass - test 001\n")
	}
}

func Test_002(t *testing.T) {
	Verbose = false
	fmt.Printf("\nTest_002 HumanTime\n")
	var tsec int64 = 10
	fmt.Printf("%d seconds is %s\n", tsec, HumanTime(time.Duration(tsec)*time.Second))

	tsec *= 10
	fmt.Printf("%d seconds is %s\n", tsec, HumanTime(time.Duration(tsec)*time.Second))

	tsec *= 10
	fmt.Printf("%d seconds is %s\n", tsec, HumanTime(time.Duration(tsec)*time.Second))

	tsec *= 10
	fmt.Printf("%d seconds is %s\n", tsec, HumanTime(time.Duration(tsec)*time.Second))

	tsec *= 10
	fmt.Printf("%d seconds is %s\n", tsec, HumanTime(time.Duration(tsec)*time.Second))

	tsec *= 10
	fmt.Printf("%d seconds is %s\n", tsec, HumanTime(time.Duration(tsec)*time.Second))

	tsec *= 10
	fmt.Printf("%d seconds is %s\n", tsec, HumanTime(time.Duration(tsec)*time.Second))

	tsec *= 10
	fmt.Printf("%d seconds is %s\n", tsec, HumanTime(time.Duration(tsec)*time.Second))
	tsec *= 10
	fmt.Printf("%d seconds is %s\n", tsec, HumanTime(time.Duration(tsec)*time.Second))

	fmt.Printf("Pass - test 002\n")
	return
	// exceeds maximum duration of 290 years
	tsec *= 10
	fmt.Printf("%d seconds is %s\n", tsec, HumanTime(time.Duration(tsec)*time.Second))
}

func Test_003(t *testing.T) {
	Verbose = false
	fmt.Printf("\nTest_003 FlipCoin\n")
	var s int64 = 0
	for i := 0; i < 1000; i++ {
		if FlipCoin() {
			s++
		}
	}
	fmt.Printf("flip of 1000 coins gives %d heads\n", s)
	if !InRangeInt64(450, s, 550) {
		t.Errorf("FlipCoin() not in expected range of [450..550]")
	}
	fmt.Printf("Pass - test 003\n")
}

func Test_004(t *testing.T) {
	fmt.Printf("\nTest_004 ValidHexString and ValidHexChar\n")
	if !ValidHexString("abcdef0123456789") {
		t.Errorf("failed on abcdef0123456789")
	}
	if ValidHexString("abcdefg0123456789") {
		t.Errorf("failed on abcdefg0123456789")
	}
	fmt.Printf("Pass - test 004\n")
}

func Test_005(t *testing.T) {
	fmt.Printf("\nTest_005 FileSHA256\n")
	targetFile := "test-data/do_NOT_modify.txt"
	sha256, err := FileSHA256(targetFile)
	if err != nil {
		t.Errorf("FileSHA256(%s) failed\n", targetFile)
		return
	}
	target := "756eee95d094a7c0e1db8f9f952eab2e499fcea82ebf4fca59802870d1c2e7a6"
	if sha256 != target {
		t.Errorf("FileSHA256 did not match for %s\n", targetFile)
		return
	}
	myBuf := []byte("David Rook\n")
	myDigest := BufSHA256(myBuf)
	if myDigest != target {
		t.Errorf("BufSHA256 did not match\n")
		return
	}
	fmt.Printf("Pass - test 005\n")
}

func Test_005a(t *testing.T) {
	fmt.Printf("\nTest_005a FileMD5\n")
	targetFile := "test-data/do_NOT_modify.txt"
	md5, err := FileMD5(targetFile)
	if err != nil {
		t.Errorf("FileMD5(%s) failed\n", targetFile)
		return
	}
	target := "31b08d26a2ff669a35c20cf561083918"
	if md5 != target {
		t.Errorf("FileMD5 did not match for %s\n", targetFile)
		return
	}
	myBuf := []byte("David Rook\n")
	myDigest := BufMD5(myBuf)
	if myDigest != target {
		t.Errorf("BufMD5 did not match\n")
		return
	}
	fmt.Printf("Pass - test 005a\n")
}

func Test_005b(t *testing.T) {
	if !ValidHexString("abcdef0123456789") {
		t.Errorf("failed on abcdef0123456789\n")
		return
	}
	if ValidHexString("abcdefg0123456789") {
		t.Errorf("failed on abcdefg0123456789\n")
		return
	}
	fmt.Printf("Pass - test 005b\n")
}

func Test_006(t *testing.T) {
	fmt.Printf("\nTest_006 Reverse([]byte)\n")
	var before, after []byte
	before = []byte{1, 2, 3, 4, 5}
	after = []byte{5, 4, 3, 2, 1}
	fmt.Printf("before %v expected after reverse %v\n", before, after)
	Reverse(before)
	if len(before) != len(after) {
		t.Errorf("Reverse failed")
	}
	fmt.Printf("reversed %v should equal after %v\n", before, after)
	for i := 0; i < len(before); i++ {
		if before[i] != after[i] {
			t.Errorf("Reverse failed")
		}
	}
	fmt.Printf("Pass - test 006\n")
}

func Test_007(t *testing.T) {
	fmt.Printf("\nTest_007 CommaFmtInt64\n")
	s := CommaFmtInt64(1234567)
	fmt.Printf("formatting 1234567: expecting 1,234,567, got %s\n", s)
	if s != "1,234,567" {
		t.Errorf("CommaFmtInt64(123456) failed")
	} else {
		fmt.Printf("Pass - test 007\n")
	}
}

func Test_008(t *testing.T) {
	fmt.Printf("\nTest_008 FileLength\n")
	targetFile := "test-data/do_NOT_modify.txt"
	flen, err := FileLength(targetFile)
	if (flen != 11) || (err != nil) {
		t.Errorf("FileLength(%s) failed", targetFile)
	} else {
		fmt.Printf("Pass - test 008\n")
	}
}

func Test_009(t *testing.T) {
	fmt.Printf("\nTest_009 Archive Type Test\n")

	type target struct {
		s string
		a ArchiveType
		c CompressType
	}

	var targetList = []target{
		// just compression
		{"abc.Z", ArchiveNoMatchType, CompressZcompressType},
		{"abc.z", ArchiveNoMatchType, CompressZcompressType}, // possible conflict with pack files
		{"abc.gz", ArchiveNoMatchType, CompressGzipType},
		{"abc.gzip", ArchiveNoMatchType, CompressGzipType},
		{"abc.bz", ArchiveNoMatchType, CompressBz2Type}, // possible conflict with pack files
		{"abc.bz2", ArchiveNoMatchType, CompressBz2Type},
		{"abc.bzip2", ArchiveNoMatchType, CompressBz2Type},
		// combined compression and archive in one ext
		{"abc.taz", ArchiveTarType, CompressZcompressType},
		{"abc.tbz", ArchiveTarType, CompressBz2Type},
		{"abc.tgz", ArchiveTarType, CompressGzipType},
		{"abc.zip", ArchiveZipType, CompressZipType},
		// combined compression and archive in last two
		{"abc.tar.Z", ArchiveTarType, CompressZcompressType},
		{"abc.ark.Z", ArchiveArkType, CompressZcompressType},
		{"abc.tar.bz", ArchiveTarType, CompressBz2Type},
		{"abc.tar.gz", ArchiveTarType, CompressGzipType},
		{"abc.tar.bz2", ArchiveTarType, CompressBz2Type},
		{"abc.tar.bzip2", ArchiveTarType, CompressBz2Type},
		// just archive not compressed
		{"abc.ark", ArchiveArkType, CompressNoMatchType},
		{"abc.cpio", ArchiveCpioType, CompressNoMatchType},
		{"abc.tar", ArchiveTarType, CompressNoMatchType},

		// not implemented yet - and many others...
		{"abc.ar", ArchiveNoMatchType, CompressNoMatchType},
	}

	for ndx, x := range targetList {
		fmt.Printf("%d %s %v %v ", ndx, x.s, x.a, x.c)
		if WhichArchiveType(x.s) != x.a {
			t.Errorf("WhichArchiveType(%s) failed", x.s)
		}
		if WhichCompressType(x.s) != x.c {
			t.Errorf("WhichCompressType(%s) failed", x.s)
		}
		fmt.Printf(" PASS\n")
	}
}

func Test_010(t *testing.T) {
	fmt.Printf("\nTest_010 PseudoRandomBlock() Test\n")
	for i := 0; i < 10; i++ {
		x := PseudoRandomBlock(10)
		fmt.Printf("x(%x)\n", x)
	}
}

func Test_011(t *testing.T) {
	fmt.Printf("\nTest_011 LSBytesFromInt64 Test\n")
	i := int64(1)
	b := LSBytesFromInt64(i)
	fmt.Printf("%d\n", i)
	fmt.Printf("%x\n", b)
	if !bytes.Equal([]byte{1, 0, 0, 0, 0, 0, 0, 0}, b) {
		t.Errorf("LSBytesFromInt64 failed\n")
	}
}

func Test_012(t *testing.T) {
	fmt.Printf("\nTest_012 Int64FromLSBytes Test\n")
	i := int64(15)
	b := LSBytesFromInt64(i)
	fmt.Printf("%x\n", b)
	j := Int64FromLSBytes(b)
	fmt.Printf("%d\n", j)
	if i != j {
		t.Errorf("Int64FromLSBytes failed at A\n")
	}
}

func Test_013(t *testing.T) {
	fmt.Printf("\nTest_013 MSBytesFromInt64 Test\n")
	i := int64(15 * 256)
	b := MSBytesFromInt64(i)
	if !bytes.Equal([]byte{0, 0, 0, 0, 0, 0, 0xf, 0}, b) {
		t.Errorf("MSBytesFromInt64 failed at A\n")
	}
	fmt.Printf("%x\n", b)

	i = int64(-1)
	b = MSBytesFromInt64(i)
	fmt.Printf("%x\n", b)
	if !bytes.Equal([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, b) {
		t.Errorf("MSBytesFromInt64 failed at B\n")
	}
}

func Test_014(t *testing.T) {
	fmt.Printf("\nTest_014 Int64FromMSBytes Test\n")
	b := []byte{0, 0, 0, 1, 0, 0, 2, 0}
	i := Int64FromMSBytes(b)
	fmt.Printf("i(%d)\n", i)
	var rv int64
	rv = (1 << 32) + 2*(1<<8)
	if i != rv {
		t.Errorf("MSBytesFromInt64 failed at A\n")
	}
	fmt.Printf("%x\n", b)
}

/////////////////////////  B E N C H M A R K S  ////////////////////////////

// 46.9e6 ns/op on 4Ghz AMD64 with 1.0.3
// 36.6e6 ns/op on 4Ghz AMD64 with 1.1 << 22% better >>
// 35.1e6 ns/op on 4Ghz AMD64 with 1.2
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
func Benchmark_BufSHA256(b *testing.B) {
	testBuf := PseudoRandomBlock(1024 * 1024)
	for i := 0; i < b.N; i++ {
		_ = BufSHA256(testBuf)
	}
}

// 2.71e6 ns/op on 4Ghz AMD64 with 1.0.3
// 3.25e6 ns/op on 4Ghz AMD64 with 1.1  << 18% worse >>
// 3.19e6 ns/op on 4Ghz AMD64 with 1.2
func Benchmark_BufCRC64(b *testing.B) {
	testBuf := PseudoRandomBlock(1024 * 1024)
	for i := 0; i < b.N; i++ {
		_ = crc64.Checksum(testBuf[:], G_crcTable)
	}
}

/////////////////////////  E X A M P L E S  ////////////////////////////
func workerFunction(w IntPair) {
	fmt.Printf("work on items from %d through %d\n", w.X, w.Y)
}

func ExampleJobSplit() {
	nCPU := 4
	totalWork := 100
	jobrange := JobSplit(totalWork, nCPU)
	for i := 0; i < nCPU; i++ {
		go workerFunction(jobrange[i])
	}
}

func ExampleProgressBar() {
	var (
		status    int64
		endNumber int64 = 100
	)

	progChan := make(chan int64, 2)
	go ProgressBar(50, progChan, endNumber) // start the display handler
	progChan <- 0                           // make first progress display visible
	for {
		//      ... do something to advance status towards endNumber ...
		time.Sleep(time.Second)
		status += 10
		progChan <- status
		if status >= endNumber {
			break
		}
	}
	progChan <- -1 // close up shop
}
