// mdr.go   (c) 2013,2014 David Rook - License is BSD style - see LICENSE.md
// Utility Functions
//
//  Features
//  ========
//    Most of the functions are short and easily understood
//    Examples available for the less obvious
//
// Real BUGS  -##0{   None Known - but beware of limitations
//
//  Limitations
//  -----------
//    GetKey() doesn't hide key entry
//
//
//  Also see README-mdr.md for more info
//
package mdr

// BUG(mdr): TODO need more test cases
import (
	// uses standard lib go 1.2 pkgs below
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"syscall"
)

type Ints []int

type IntPair struct {
	X, Y int
}

var (
	g_UnusualMode os.FileMode = os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice
	Verbose                   = false
)

// return true randomly half the time
// Test case is Test_003
func FlipCoin() bool {
	return rand.Int31n(2) == 0
}

// reverse a slice of bytes in place
// Test is Test_006
func Reverse(b []byte) {
	first := 0
	last := len(b) - 1
	for first < last {
		b[first], b[last] = b[last], b[first]
		first++
		last--
	}
}

// used as helper with filters to collect args from input stream
// limited to RAM for size of return []string
func GetAllArgs() []string {
	rv := make([]string, 0, 1000)
	f := os.Stdin // f is * osFile
	rdr := bufio.NewReader(f)
	alldone := false
	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				alldone = true
			} else {
				log.Panicf("mdr: GetAllArgs read error")
			}
		}
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			rv = append(rv, line)
		}
		if alldone {
			break
		}
	}
	if flag.Parsed() {
		args := flag.Args()
		for _, arg := range args {
			rv = append(rv, arg)
		}
	} else {
		fmt.Printf("Warning --> GetAllArgs: flags not parsed yet\n")
	}
	return rv
}

func ValidDecChar(c byte) bool {
	var decchars []byte = []byte("0123456789")
	for _, d := range decchars {
		if c == d {
			return true
		}
	}
	return false
}

func ValidDecString(s string) bool {
	for _, c := range s {
		if !ValidDecChar(byte(c)) {
			return false
		}
	}
	return true
}

// Test_004
func ValidHexChar(c byte) bool {
	var hexchars []byte = []byte("0123456789abcdefABCDEF")
	for _, h := range hexchars {
		if c == h {
			return true
		}
	}
	return false
}

// Test_004
func ValidHexString(s string) bool {
	for _, c := range s {
		if !ValidHexChar(byte(c)) {
			return false
		}
	}
	return true
}

// Test is Test_008
func FileLength(fname string) (int64, error) {
	stats, err := os.Stat(fname)
	if err != nil {
		fmt.Printf("mdr: Can't get fileinfo for %s\n", fname)
		return -1, err
	}
	// fmt.Printf("%s %v\n",fname,stats)
	return stats.Size(), nil
}

///////////////////////////  N E E D   T E S T   C A S E S  ///////////////////

// return files userid number
func FileUID(fname string) (int, error) {
	fi, err := os.Stat(fname)
	if err != nil {
		return -1, err
	}
	//fmt.Printf("fi %v\n\n\n", fi)
	sys := fi.Sys().(*syscall.Stat_t)
	//fmt.Printf("UID = %d\n", int(sys.Uid))
	return int(sys.Uid), nil
}

// return files groupid number
func FileGID(fname string) (int, error) {
	fi, err := os.Stat(fname)
	if err != nil {
		return -1, err
	}
	//fmt.Printf("fi %v\n\n\n", fi)
	sys := fi.Sys().(*syscall.Stat_t)
	//fmt.Printf("GID = %d\n", int(sys.Gid))
	return int(sys.Gid), nil
}

func FileIsRegular(fname string) (bool, error) {
	info, err := os.Stat(fname)
	if err != nil {
		return false, err
	}
	fileMode := info.Mode()
	// directory is Not G_Unusual
	if (fileMode & g_UnusualMode) != 0 { // its not a regular file
		return false, nil
	}
	return true, nil
}

// used in image manipulation
func RangeMinMaxPoint(v []IntPair) (minPt, maxPt IntPair) {
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

func Crash(reason string) {
	fmt.Printf("Stopping now - I crashed because %s\n", reason)
	os.Exit(1)
}

func FatalError(err error) {
	fmt.Printf("%v\n", err)
	os.Exit(1)
}

// <end>
