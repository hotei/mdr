// mdr_misc.go (c) 2015 David Rook

package mdr

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	//"math/rand"
	"math"
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
)

// FlipCoin returns true randomly half the time
// Test with mdr_test.go:Test_003
func FlipCoin() bool {
	return GenFlipCoin()
}

// reverse a slice of bytes in place
// Test with mdr_test.go:Test_006
func Reverse(b []byte) {
	first := 0
	last := len(b) - 1
	for first < last {
		b[first], b[last] = b[last], b[first]
		first++
		last--
	}
}

// singleCharRead might return newline if no other input
//  otherwise first character if more than one on a line
// ^D will cause EOF to be printed and ? as the char returned - better choice is ...
func SingleCharRead() byte {
	var buf = []byte{0}
	_, err := os.Stdin.Read(buf)
	if err != nil {
		fmt.Printf("%v\n", err)
		return '?'
	}
	return buf[0]
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

// ZRotate rotates a point about 0,0 (z axis).
func (loc *Point) ZRotate(radians float64) {
	if radians == 0.0 {
		return
	}
	if (*loc == Point{0, 0}) {
		return
	}
	Verbose.Printf("ZRotate starts with loc(%v), turn by  %v radians --> ", loc, radians)
	polarAngle, polarRadius := Polar(*loc)
	Verbose.Printf("A %v R %v\n", polarAngle, polarRadius)
	polarAngle += radians
	*loc = Cartesian(polarAngle, polarRadius)
	Verbose.Printf("new loc(%v)\n", *loc)
	return
}

// RotateOnPivotPt rotate a point using another point as a pivot (mimic mechanical compass drawing).
func (loc *Point) RotateOnPivotPt(p Point, radians float64) {
	loc.X -= p.X
	loc.Y -= p.Y
	loc.ZRotate(radians)
	loc.X += p.X
	loc.Y += p.Y
}

// Cartesian returns Cartesian point from polar point.
// NoteBene: NOT the usual [in trig x=cos(a) y=sin(a)]
// because y axis is inverted in 'conventional' computer graphics.
func Cartesian(angle, radius float64) Point {
	var rv Point
	rv.X = int(math.Sin(angle) * radius)
	rv.Y = int(math.Cos(angle) * radius)
	return rv
}

// video   3 | 2       function graph Quad 2 | 1
//         --|--                           --|--
//         4 | 1                           3 | 4
//

// Polar returns the polar coords (angle and radius) for a Cartesian point.
func Polar(loc Point) (theta, r float64) {
	if (loc == Point{0, 0}) {
		return 0, 0 // The angle is actually undefined but this will do.
	}
	x, y := float64(loc.X), float64(loc.Y)
	return math.Mod(math.Atan2(x, y)+2*math.Pi, 2*math.Pi), math.Hypot(x, y)
}

// radians returns the equivalent to degrees.
func Radians(degrees float64) float64 {
	return degrees * RadiansPerDegree
}

// PolarAngle returns the angle produced by a line from Point{0,0} to loc.
func PolarAngle(loc Point) float64 {
	if (loc == Point{0, 0}) {
		// possibly not an error in some cases so
		return 0.0 // The angle is actually undefined but this will do.
	}
	x, y := float64(loc.X), float64(loc.Y)
	rv := math.Atan2(x, y)
	if rv < 0.0 {
		rv += PiX2
	}
	return rv
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

func FileLinkCt(fname string) (int, error) {
	fi, err := os.Stat(fname)
	if err != nil {
		return -1, err
	}
	//fmt.Printf("fi %v\n\n\n", fi)
	sys := fi.Sys().(*syscall.Stat_t)
	//fmt.Printf("Nlink = %d\n", int(sys.Nlink))
	return int(sys.Nlink), nil
}

// BUG(mdr) FileIsRegular now OBE (since go 1.1)
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

func Crash(reason string) {
	fmt.Printf("Stopping now - I crashed because %s\n", reason)
	os.Exit(1)
}

func FatalError(err error) {
	fmt.Printf("%v\n", err)
	os.Exit(1)
}

// checkInterfaces - see if listener is bound to correct interface
// first is localhost, second should be IP4 of active card,
// third is IP6 localhost, fourth is IP6 for active card (on this system)
// on BSD it's [ IP4 IP6 LocalHostIP4 LocalHostIP6 LocalHostIP4]
//
// Order isn't important as long as requested inteface is there somewhere
// actual check needs to do a string match of interfaces we have with a
// target of the requested interface.  If target isn't present then stop.
//
func HasInterface(hostIPStr string) bool {
	Verbose.Printf("running interface check\n")
	ifa, err := net.InterfaceAddrs()
	if err != nil || (len(ifa) < 1) {
		fmt.Printf("!Err---> HasInterface: Can't list interfaces\n")
		log.Panic("")
		os.Exit(1)
	}
	for i := 0; i < len(ifa); i++ {
		Verbose.Printf("Interface[%d] = %v\n", i, ifa[i])
	}

	for i := 0; i < len(ifa); i++ {
		myIfs := strings.Split(ifa[i].String(), "/")
		for _, v := range myIfs {
			if strings.Contains(v, hostIPStr) {
				Verbose.Printf("Found the requested interface %s\n", v)
				return true
			}
		}
	}
	return false
}

// <end>
