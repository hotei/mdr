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
//  Also see README.md for more info
//
package mdr

// BUG(mdr): TODO need more test cases
import (
	"github.com/hotei/datatable"

	// also uses standard lib go 1.2 pkgs below
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"flag"
	"fmt"
	"hash"
	"hash/crc64"
	"io"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

/*
type Point struct {
	X, Y int
}
*/

type Ints []int

type IntPair struct {
	X, Y int
}

var NormalZtable *datatable.Table

var (
	G_crcTable    *crc64.Table
	g_UnusualMode os.FileMode = os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice
	Verbose                   = false
	g_strictZ                 = false
)

func init() {
	G_crcTable = crc64.MakeTable(crc64.ECMA)
}

// turns duration into decimal minutes,hours,days as appropriate,
// maximum duration is about 290 years
// Test_002
func HumanTime(t time.Duration) (rs string) {
	sec := int64(t.Seconds()) // converts duration in nanosec to seconds
	if sec < 60 {
		rs = fmt.Sprintf("%d seconds", sec)
		return
	}
	if sec < 3600 {
		rs = fmt.Sprintf("%5.2f minutes", float64(sec)/60.0)
		return
	}
	if sec < 86400 {
		rs = fmt.Sprintf("%5.2f hours", float64(sec)/3600.0)
		return
	}
	if sec < int64(86400.0*30.4375) {
		rs = fmt.Sprintf("%5.2f days", float64(sec)/86400.0)
		return
	}
	if sec < int64(86400.0*365.25) {
		rs = fmt.Sprintf("%5.2f months", float64(sec)/(86400.0*30.4375))
		return
	}
	rs = fmt.Sprintf("%5.2f years", float64(sec)/(86400.0*365.25))
	return
}

func LeapYear(when time.Time) bool {
	year := when.Year()
	if year%400 == 0 {
		return true
	}
	if year%100 == 0 {
		return false
	}
	if year%4 == 0 {
		return true
	}
	return false
}

// returns decimal year to nearest 52 minutes - usually printed with %9.4f
func StarDate(when time.Time) float64 {
	yr := float64(when.Year())
	dayofyear := float64(when.YearDay())
	//fmt.Printf("%v %v\n", yr, dayofyear)
	var daysinyear float64
	if LeapYear(when) {
		daysinyear = 366
	} else {
		daysinyear = 365
	}
	hrs := (dayofyear-1)*24 + float64(when.Hour())
	//fmt.Printf("hrs = %v \n",hrs)
	return yr + hrs/(daysinyear*24)
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

func ThirtyTwoNet(n []byte) uint32 {
	return Uint32FromMSBytes(n)
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
//    - beware magic numbers
//    - presumes knowledge of net.IP internals order
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

// split n into NumCPUs ranges,
//   JobSplit(10,1) -> returns [ {0,9} ]
//   JobSplit(10,2) -> returns [ {0,4},{5,9} ]
//   JobSplit(10,3) -> returns [ {0,3}, {4,6}, {7,9} ]
//       if not all slices are same length, longer ones will occur first
// Test_001
// See also ExampleJobSplit()
func JobSplit(n int, NumCPUs int) []IntPair {
	if NumCPUs < 1 {
		fmt.Printf("mdr: Jobsplit() thinks %d is too few CPUS \n", NumCPUs)
		os.Exit(-1)
	}
	if NumCPUs > 1000 {
		fmt.Printf("mdr: Jobsplit() thinks %d is too many CPUS \n", NumCPUs)
		os.Exit(-1)
	}
	if Verbose {
		fmt.Printf("mdr: Jobsplit() NumCPUs(%d)\n", NumCPUs)
	}
	rc := make([]IntPair, 0, NumCPUs)
	if Verbose {
		fmt.Printf("mdr: Jobsplit() splitting %d into %d pieces\n", n, NumCPUs)
	}
	if NumCPUs == 1 {
		if Verbose {
			rc = append(rc, IntPair{0, n - 1})
			fmt.Printf("mdr: Jobsplit() no split required, range is 0 to %d\n", n-1)
		}
		return rc
	}
	splitInc := n / NumCPUs
	excess := n - (splitInc * NumCPUs)

	if Verbose {
		fmt.Printf("mdr: Jobsplit() increment =  %d, excess = %d\n", splitInc, excess)
	}
	leftSide := 0
	rightSide := splitInc - 1
	rc = append(rc, IntPair{leftSide, rightSide})
	maxRight := n - 1

	for i := 1; i < NumCPUs; i++ {
		if excess != 0 {
			rightSide++
			excess--
		}
		if i == NumCPUs {
			rightSide = maxRight
		}
		pcs := rightSide - leftSide + 1

		if Verbose {
			fmt.Printf("mdr: Jobsplit()  [ %d , %d ]  %d items in this piece \n", leftSide, rightSide, pcs)
		}

		rc[i-1].X = leftSide
		rc[i-1].Y = rightSide
		leftSide = rightSide + 1
		rightSide += splitInc
		rc = append(rc, IntPair{leftSide, rightSide})
	}
	pcs := maxRight - leftSide + 1
	if Verbose {
		fmt.Printf("mdr: Jobsplit()  [ %d , %d ]  %d items in this piece \n", leftSide, maxRight, pcs)
	}

	return rc
}

// return true randomly half the time
// Test case is Test_003
func FlipCoin() bool {
	return rand.Int31n(2) == 0
}

// true IFF a <= b <= c || a >= b >= c, note a < c not a given
func InRangeInt64(a, b, c int64) bool {
	if a > c { // swap ends if necessary
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

func DigestToString(digest []byte) string {
	rv := ""
	for _, hex := range digest {
		rv = rv + fmt.Sprintf("%02x", hex)
	}
	return string(rv)
}

func BufMD5(buf []byte) string {
	var h hash.Hash = md5.New()
	h.Write(buf[:])
	digest := h.Sum(nil)
	//fmt.Printf("%s digest ->%s\n",fname,rv)
	return DigestToString(digest)
}

func BufSHA256(buf []byte) string {
	var h hash.Hash = sha256.New()
	h.Write(buf[:])
	digest := h.Sum(nil)
	//fmt.Printf("%s digest ->%s\n",fname,rv)
	return DigestToString(digest)
}

// Compute md5 for given file
// Test is Test_005
func FileMD5(fname string) (string, error) {
	//fmt.Printf("Computing digest for %s\n",fname)
	const NBUF = 1 << 20 // 20 -> 1 MB
	file, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if file == nil {
		fmt.Printf("!Err--> mdr.FileMD5() can't open file %s as readonly; err=%v\n", fname, err)
		return "", err
	}
	defer file.Close()
	buf := make([]byte, NBUF)
	var h hash.Hash = md5.New()
	for {
		numRead, err := file.Read(buf)
		if (err == io.EOF) && (numRead == 0) {
			break
		} // end of file reached
		if (err != nil) || (numRead < 0) {
			fmt.Fprintf(os.Stderr, "!Err--> mdr.FileMD5: error reading from %s: %v\n", fname, err)
			return "", err
		}
		//		fmt.Printf("read(%d) bytes\n",numRead)
		h.Write(buf[0:numRead])
	}
	digest := h.Sum(nil)
	rv := ""
	for _, hex := range digest {
		rv = rv + fmt.Sprintf("%02x", hex)
	}
	//fmt.Printf("%s digest ->%s\n",fname,rv)
	return string(rv), nil
}

// Compute sha256 for given file
// Test is Test_005
func FileSHA256(fname string) (string, error) {
	//fmt.Printf("Computing digest for %s\n",fname)
	const NBUF = 1 << 20 // 20 -> 1 MB
	file, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if file == nil {
		fmt.Printf("mdr: FileSHA256() can't open file %s as readonly; err=%v\n", fname, err)
		return "", err
	}
	defer file.Close()
	buf := make([]byte, NBUF)
	var h hash.Hash = sha256.New()
	for {
		numRead, err := file.Read(buf)
		if (err == io.EOF) && (numRead == 0) {
			break
		} // end of file reached
		if (err != nil) || (numRead < 0) {
			fmt.Fprintf(os.Stderr, "error reading from %s: %v\n", fname, err)
			return "", err
		}
		//		fmt.Printf("read(%d) bytes\n",numRead)
		h.Write(buf[0:numRead])
	}
	digest := h.Sum(nil)
	rv := ""
	for _, hex := range digest {
		rv = rv + fmt.Sprintf("%02x", hex)
	}
	//fmt.Printf("%s digest ->%s\n",fname,rv)
	return string(rv), nil
}

///////////////////////////  N E E D   T E S T   C A S E S  ///////////////////

// returns a buffer of specified size filled with random bits from /dev/urandom
// there is a limited supply (about 4096 bits) so read will fail if urandom has not had time
// to accumulate enough bits to satisfy the read.
// Max readable under optimum conditios ?  actual test => 1020,1932,1020
// 1 ms sleep not enough(1940fail), 10 not enough (1020fail), 50 not enough 1020fail

func UdevRandomBlock(blksize int) []byte {
	var buf = make([]byte, blksize)
	input, err := os.Open("/dev/urandom")
	if err != nil {
		FatalError(err)
	}
	n, err := input.Read(buf)
	if err != nil {
		FatalError(err)
	}
	if n != blksize { // must read /dev/urandom for full block of data
		fmt.Printf("mdr: RandomBlock failed, only %d bytes read - expected %d\n", n, blksize)
		os.Exit(1)
	}
	return buf
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

// see also go.crypto/ssh/terminal
func SetTtyRawMode() {
	cmd := exec.Command("/bin/stty", "-F", "/dev/tty", "-icanon", "min", "1") // worked ok
	err := cmd.Run()
	if err != nil {
		FatalError(err)
	}
}

// returns buffer of specified size filled with random bytes
// test funcs are  Test_010 and Benchmark_010
func PseudoRandomBlock(blksize int) []byte {
	// tmp buffer initial size must be a multiple of SizeOf(Uint32) ie. *4
	rvsize := blksize + 4
	rvsize >>= 2
	rvsize <<= 2
	rvb := make([]byte, 0, rvsize)
	seed := UdevRandomBlock(8)
	seedInt64 := Int64FromLSBytes(seed)
	rand.Seed(seedInt64)
	for {
		r := rand.Uint32()
		tmp := LSBytesFromUint32(r)
		rvb = append(rvb, tmp...)
		if len(rvb) >= blksize {
			break
		}
	}
	// slice it so we only return number of bytes requested
	return rvb[:blksize]
}

// BUG(mdr): GetKey() key is visible during entry
func GetKey() string {
	var buf = make([]byte, 200)
	fmt.Printf("Enter key: ")
	num_in, err := os.Stdin.Read(buf)
	// fmt.Printf("buf(%v), num_in(%d), err(%v)\n",buf, num_in, err)
	var localkey string
	if err == nil {
		localkey = string(buf[:num_in-1])
	}
	if len(localkey) <= 0 {
		localkey = "DefaultKey"
	}
	if false {
		fmt.Printf("Key(%q)\n", localkey)
	}
	return localkey
}

// returns the absolute value of an int64
func AbsInt64(a int64) int64 {
	if a < 0 {
		return -a
	}
	return a
}

// convert an int64 into an 8 byte array, LSB first (littleEndian)
// test is Test_011
func LSBytesFromInt64(n int64) []byte {
	rv := make([]byte, 0, 8)
	for i := 0; i < 8; i++ {
		rv = append(rv, byte(n%256))
		n >>= 8
	}
	return rv
}

// convert int64 to [0:8]byte slice with MSB first
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
// Test is Test_014
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
// Test is Test_012
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

// syntax sugar
func BufCRC64(bufr []byte) uint64 {
	return crc64.Checksum(bufr[:], G_crcTable)
}

// reads file in 512KB chunks
func FileCRC64(fname string) (uint64, error) {
	var checksum uint64
	var bufSize = 1024 * 512
	fp, err := os.Open(fname)
	if err != nil {
		return 0, err
		//fmt.Printf("!Err-> can't read %s\n", fname)
		//log.Panicf("!Err-> can't read %s %s\n", fname, err)
	}
	defer fp.Close()
	bufr := make([]byte, bufSize)
	checksum = 0
	for {
		n, err := io.ReadFull(fp, bufr)
		if err != nil {
			if n > 0 {
				// fmt.Printf("not full block but read %d bytes\n",n)
			}
		}
		if err == io.EOF {
			if n != 0 {
				return 0, err
			} else {
				break // normal end of file reached
			}
		}
		// fmt.Printf(" read %d bytes\n",n)
		checksum = crc64.Update(checksum, G_crcTable, bufr[0:n])
		// checksum = crc64.Checksum(bufr[0:n],crcTable)  if entire string in bufr
	}
	return checksum, nil
}

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

/////////////////////////////// A R C H I V E ///////////////////////////
/// BEGIN

type ArchiveType int

const (
	ArchiveNoMatchType ArchiveType = iota
	ArchiveTarType
	ArchiveZipType
	ArchiveArkType
	ArchiveCpioType
)

type dispatch struct {
	pattern string
	arch    ArchiveType
}

var dispatcher = []dispatch{}

// examine name to determine archive method used
// should match magic numbers but that's a different function
func WhichArchiveType(s string) ArchiveType {
	ext := strings.ToLower(filepath.Ext(s))
	// ArchiveZipType -----------------------------------------
	if ext == ".zip" {
		return ArchiveZipType
	}
	// ArchiveTarType -----------------------------------------
	if ext == ".tar" {
		return ArchiveTarType
	}
	if ext == ".tgz" {
		return ArchiveTarType
	}
	if ext == ".taz" {
		return ArchiveTarType
	}
	if ext == ".tbz" {
		return ArchiveTarType
	}
	// ArchiveCpioType -----------------------------------------
	if ext == ".cpio" {
		return ArchiveCpioType
	}
	// ArchiveArkType -----------------------------------------
	if ext == ".ark" {
		return ArchiveArkType
	}

	if len(dispatcher) == 0 {
		if Verbose {
			fmt.Printf("Building dispatcher\n")
		}
		// compressed tar collection
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.gz\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.Z\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.z\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.bz\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.bz2\\E$", ArchiveTarType})
		dispatcher = append(dispatcher, dispatch{".*\\Q.tar.bzip2\\E$", ArchiveTarType})
		// compressed ark collection
		dispatcher = append(dispatcher, dispatch{".*\\Q.ark.Z\\E$", ArchiveArkType})
	}

	nameBytes := []byte(s)
	for _, d := range dispatcher {
		isMatch, err := regexp.Match(d.pattern, nameBytes)
		if err != nil {
			fmt.Printf("!Err-> ?pattern error in %s\n", d.pattern)
			log.Panicf("better fix it now")
		}
		if isMatch {
			return d.arch
		} // only need the first match
	}

	// NO MATCH -----------------------------------------
	return ArchiveNoMatchType
}

type CompressType int

const (
	CompressNoMatchType CompressType = iota
	CompressZipType
	CompressGzipType
	CompressBz2Type
	CompressZcompressType // .Z not common except possibly in Japan
	CompressPackType      // .z deprecated - very rare in last 20 years
	CompressBz1Type       // .bz deprecated - very rare in last 20+ years
)

// examine name to determine compression method used
// should match magic numbers but that's a different function
func WhichCompressType(s string) CompressType {
	origExt := filepath.Ext(s)
	// note - we don't do the next step because a lot of the .z files we have
	// found are simply .Z that have been mislabeled, thus strictZ defaults false
	// if you have the opposite, then set the global true instead
	// the caller  should test magic header bytes in file body to see which is correct
	if g_strictZ {
		if origExt == ".z" {
			return CompressPackType
		}
	}
	// unless we have a good reason not to we fold case for the following tests
	ext := strings.ToLower(origExt)
	// CompressZipType ---------------------------------
	if ext == ".zip" {
		return CompressZipType
	}
	// CompressGzipType ---------------------------------
	if ext == ".gz" {
		return CompressGzipType
	}
	if ext == ".gzip" {
		return CompressGzipType
	}
	if ext == ".tgz" {
		return CompressGzipType
	}
	// CompressZcompressType ---------------------------------
	if ext == ".z" { // .Z converted to .z by tolower - see above - usually not a bug
		return CompressZcompressType
	}
	if ext == ".taz" {
		return CompressZcompressType
	}
	// CompressBz2Type ---------------------------------
	if ext == ".bz" {
		return CompressBz2Type
	}
	if ext == ".tbz" {
		return CompressBz2Type
	}
	if ext == ".bz2" {
		return CompressBz2Type
	}
	if ext == ".tbz2" {
		return CompressBz2Type
	}
	if ext == ".bzip2" {
		return CompressBz2Type
	}
	// NO MATCH -----------------------------------------
	return CompressNoMatchType
}

///  END
/////////////////////////////// A R C H I V E ///////////////////////////

// keep user entertained while something happens behind the curtain
//  see example from mdr_test.go for usage
//  see also Spinner()
func ProgressBar(barWidth int, p chan int64, alldone int64) {
	fmt.Fprintf(os.Stderr, "Progress bar:\n")
	bar := "++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
	nobar := "...................................................................."
	barlen := 0
	s := bar[:barlen] + nobar[:barWidth-barlen]
	lastBarlen := 0
	fmt.Fprintf(os.Stderr, "%s\r", s)
	for {
		progress := <-p
		if progress < 0 {
			break
		}
		if progress > alldone {
			progress = alldone
		}
		barlen = int(int64(progress*int64(barWidth)) / alldone)
		if barlen != lastBarlen {
			s = bar[:barlen] + nobar[:barWidth-barlen]
			fmt.Fprintf(os.Stderr, "%s\r", s)
			lastBarlen = barlen
		} else {
			// nothing
		}
	}
	fmt.Fprintf(os.Stderr, "\nDone\n")
}

// these would normally be C static vars
var spinCt int8

const spinchars string = "|/-\\|/-\\ "

// keep user entertained while something happens behind the curtain
//  see example from mdr_test.go for usage
//  Choose this if progress bar can't be used because endpoint not known
func Spinner() {
	fmt.Fprintf(os.Stderr, "%s\r", spinchars[spinCt:spinCt+1])
	spinCt++
	spinCt &= 0x7 // mod 8 which is length of spinchars by design
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

func InRangeF(a, b, c float64) bool {
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

func init() {
	// table contains : z value, area (ie. probability) to left of z value
	// upper half of table only, negative values = (1.0 - upper value)
	NormalZtable = new(datatable.Table)
	NormalZtable.Name = "normal_Z_table"
	NormalZtable.Data =
		[]datatable.DblPair{{0.0, 0.5000},
			{0.25, 0.5987},
			{0.5, 0.6915},
			{0.75, 0.7734},
			{1.0, 0.8413},
			{1.5, 0.9332},
			{2.0, 0.9722},
			{2.5, 0.9938},
			{3.0, 0.9987},
			{3.49, 0.9998},
			{4.00, 0.99999},
			{5.00, 0.999995},
			{6.00, 0.9999966},
			{60.00, 1.00},
		}
}

// return float64 in range {0 .. r}
func GenRandomUniform(r float64) float64 {
	return rand.Float64() * r
}

func GenRandomUniformLoHi(low, high float64) float64 {
	r := high - low
	return low + GenRandomUniform(r)
}

func GenRandomZNormal() float64 {
	rnd := 0.5 + (rand.Float64() / 2.0)
	rv, err := NormalZtable.ReverseEval(rnd)
	if err != nil {
		Crash(fmt.Sprintf("reverse eval of normalZtable failed with err %v", err))
	}
	if len(rv) > 1 {
		fmt.Printf("got more than one return value\n")
	}
	if FlipCoin() {
		rv[0] = -rv[0]
	}
	return rv[0]
}

func GenRandomNormal(mu, stdev float64) float64 {
	rnd := GenRandomZNormal()
	dev := stdev * rnd
	if FlipCoin() {
		dev = -dev
	}
	rv := mu + dev
	return rv
}

func GenRandomPoisson(lambda float64) int {
	L := 1.0 / math.Exp(lambda)
	k := 0
	p := 1.0
	for {
		k += 1
		u := rand.Float64()
		p *= u
		if p <= L {
			break
		}
	}
	return k
}

func (a Ints) ContainsI(b int) bool {
	for _, val := range a {
		if val == b {
			return true
		}
	}
	return false
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
