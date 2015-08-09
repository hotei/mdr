<center>
mdr
===
</center>

License details are at the end of this document. 
This document is (c) 2012-2015 David Rook.

Comments can be sent to <hotei1352@gmail.com> .

This is a package of "utility" code I wrote.  I use it frequently.  If you have
received one of my other packages from github (ansiterm, bits, etc ...) you may have
gotten __mdr__ as a dependency.  It pulls in a fairly large range of standard lib
packages so if you only need a smallish set of things it might make sense to just
copy them individually or possibly make a package subset from __mdr__.  If you find it
useful - or find a bug - please send an email.  

## Installation

If you have a working go installation on a Unix-like OS:

> ```go get github.com/hotei/mdr```

## Features

## Style

```
I prefer not to use the single line form :

	if x,err := foo(); err != nil {
	// stuff
	}

instead you'll see:
	x,err := foo()
	if err != nil {
	// stuff
	}
```

## Configuration

* Note that the CRC64 table is Public if you need to replace it without changing
the library.
* I use fatal errors rather than panics in most places.  This is a habit of mine and may
not always lead to the fastest debugging.  I just prefer the end user see a smaller
understandable (I hope) message and not the gut-spilling verbosity of panic.

### Change Log

* 2015-05-01 updated progress bar functions
 * validate with 1.4.2
* 2013-04-10 updated docs, posted at github.com/hotei/MDR.git
* additions
* 2010-04-20 started, working

### Resources

* [go language reference] [1] 
* [go standard library package docs] [2]
* [Source for mdr package on github] [3]

[1]: http://golang.org/ref/spec/ "go reference spec"
[2]: http://golang.org/pkg/ "go package docs"
[3]: http://github.com/hotei/mdr "github.com/hotei/mdr"


License
-------
The 'mdr' go package is distributed under the Simplified BSD License:

> Copyright (c) 2010-2015 David Rook. All rights reserved.
> 
> Redistribution and use in source and binary forms, with or without modification, are
> permitted provided that the following conditions are met:
> 
>    1. Redistributions of source code must retain the above copyright notice, this list of
>       conditions and the following disclaimer.
> 
>    2. Redistributions in binary form must reproduce the above copyright notice, this list
>       of conditions and the following disclaimer in the documentation and/or other materials
>       provided with the distribution.
> 
> THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDER ``AS IS'' AND ANY EXPRESS OR IMPLIED
> WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND
> FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> OR
> CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
> CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
> SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON
> ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
> NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF
> ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// EOF README-mdr-pkg.md  (this is a markdown document and tested OK with blackfriday)
DO NOT EDIT BELOW - Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)

# mdr
    import "."

mdr.go   (c) 2013-20155 David Rook - all rights reserved

Eclectic Utility Functions


	Features
	========
	  Most of the functions are short and easily understood
	  Examples available for the less obvious

Real BUGS  -##0{   None Known - but beware of limitations


	Limitations
	-----------
	  GetKey() doesn't hide key entry
	
	Also see README-mdr.md for more info




## Constants
``` go
const (
    G_datetimeFmt = "2006-01-02:15_04_05"
)
```
``` go
const (
    UpdateDelayMs = 200 // delay in millisec between updates
)
```
keep user entertained while something happens behind the curtain


	see example from mdr_test.go for usage
	use Spinner() if endpoint of progress is unknown/unknowable


## Variables
``` go
var (
    PiX2             = math.Pi * 2.0
    RadiansPerDegree = math.Pi / 180.0
)
```
``` go
var CantCreateRec = errors.New("mdr: cant create Rec256 ")
```
``` go
var (
    G_crcTable *crc64.Table
)
```
``` go
var NormalZtable *datatable.Table
```

## func AbsF64
``` go
func AbsF64(a float64) float64
```
AbsF64 returns the absolute value of a float64


## func AbsInt64
``` go
func AbsInt64(a int64) int64
```
returns the absolute value of an int64


## func BufCRC64
``` go
func BufCRC64(bufr []byte) uint64
```
syntax sugar


## func BufMD5
``` go
func BufMD5(buf []byte) string
```

## func BufSHA256
``` go
func BufSHA256(buf []byte) string
```

## func CommaFmtInt64
``` go
func CommaFmtInt64(n int64) string
```
returns arg as comma inserted number string


	12345 becomes "12,345" (not locale sensitive)
	This is USA format, not internationalized


## func ConfineLoHi
``` go
func ConfineLoHi(lo, i, hi int) int
```

## func Crash
``` go
func Crash(reason string)
```

## func DigestToString
``` go
func DigestToString(digest []byte) string
```

## func Factorial
``` go
func Factorial(n int) int
```
Factorial computes value recursively.


## func FatalError
``` go
func FatalError(err error)
```

## func FileCRC64
``` go
func FileCRC64(fname string) (uint64, error)
```
reads file in 512KB chunks


## func FileGID
``` go
func FileGID(fname string) (int, error)
```
return files groupid number


## func FileIsRegular
``` go
func FileIsRegular(fname string) (bool, error)
```
BUG(mdr) FileIsRegular now OBE (since go 1.1)


## func FileLength
``` go
func FileLength(fname string) (int64, error)
```
Test is Test_008


## func FileLinkCt
``` go
func FileLinkCt(fname string) (int, error)
```

## func FileMD5
``` go
func FileMD5(fname string) (string, error)
```
Compute md5 for given file
Test is Test_005


## func FileNameIsCollection
``` go
func FileNameIsCollection(fname string) (bool, string)
```

## func FileSHA256
``` go
func FileSHA256(fname string) (string, error)
```
Compute sha256 for given file
Test is Test_005


## func FileUID
``` go
func FileUID(fname string) (int, error)
```
return files userid number


## func FlipCoin
``` go
func FlipCoin() bool
```
FlipCoin returns true randomly half the time
Test with mdr_test.go:Test_003


## func GenFlipCoin
``` go
func GenFlipCoin() bool
```
expect to see about 50% head, 50% tails (HasTest)


## func GenRandF64Between
``` go
func GenRandF64Between(lo, hi float64) float64
```
GenRandF64Between endpoints may occur in output


## func GenRandIntBetween
``` go
func GenRandIntBetween(lo, hi int) int
```
GenRandIntBetween endpoints may occur in output


## func GenRandIntBtw
``` go
func GenRandIntBtw(lo, hi int) int
```
RandIntBtw endpoints may occur (HasTest widget).


## func GenRandomNormal
``` go
func GenRandomNormal(mu, stdev float64) float64
```
GenRandomNormal returns a float64 with average mu and standard deviation stdev
depending on mu and stdev picked the values returned could be virtually any float64


## func GenRandomPoisson
``` go
func GenRandomPoisson(lambda float64) int
```
GenRandomPoisson returns an int with Poisson distribution
typically used to determine the number of time units before some event occurs


## func GenRandomUniform
``` go
func GenRandomUniform(r float64) float64
```
GenRandomUniform returns a float64 in range {0 .. r}


## func GenRandomUniformLoHi
``` go
func GenRandomUniformLoHi(low, high float64) float64
```
GenRandomUniformLo returns a float64 in range {low .. high}
high > low is NOT required, not sure if this is right or should panic 8888 ?


## func GenRandomZNormal
``` go
func GenRandomZNormal() float64
```
GenRandomZNormal returns a float64 with average of 0 and standard deviation of 1.0
as implemented the range of values returned will be in [-60..60]


## func GetAllArgs
``` go
func GetAllArgs() []string
```
used as helper with filters to collect args from input stream
limited to RAM for size of return []string


## func GetKey
``` go
func GetKey() string
```
BUG(mdr): GetKey() key is visible during entry


## func HumanTime
``` go
func HumanTime(t time.Duration) (rs string)
```
turns duration into decimal minutes,hours,days as appropriate,
maximum duration is about 290 years
Test_002


## func IPFromUint32
``` go
func IPFromUint32(adr uint32) net.IP
```
beware - presumes knowledge of net.IP internals order


## func InRangeF
``` go
func InRangeF(a, b, c float64) bool
```
true IFF a <= b <= c || a >= b >= c, note a < c not required


## func InRangeI
``` go
func InRangeI(a, b, c int) bool
```

## func InRangeInt64
``` go
func InRangeInt64(a, b, c int64) bool
```
true IFF a <= b <= c || a >= b >= c, note a < c not required


## func Int64FromLSBytes
``` go
func Int64FromLSBytes(b []byte) int64
```
convert an 8 byte array into an int64
see Test_012


## func Int64FromMSBytes
``` go
func Int64FromMSBytes(b []byte) int64
```
convert an 8 byte slice (BigEndian - MSB First) into an int64
see Test_014


## func JobSplit
``` go
func JobSplit(n int, NumCPUs int) []IntPair
```
split n into NumCPUs ranges,


	JobSplit(10,1) -> returns [ {0,9} ]
	JobSplit(10,2) -> returns [ {0,4},{5,9} ]
	JobSplit(10,3) -> returns [ {0,3}, {4,6}, {7,9} ]
	    if not all slices are same length, longer ones will occur first

Test_001
See also ExampleJobSplit()


## func LSBytesFromInt64
``` go
func LSBytesFromInt64(n int64) []byte
```
convert an int64 into an 8 byte array, LSB first (LittleEndian)
test is Test_011


## func LSBytesFromUint32
``` go
func LSBytesFromUint32(u uint32) []byte
```
convert uint32 to [0:4]byte slice in LSB first order
reverse function is ThirtyTwoBit()


## func LeapYear
``` go
func LeapYear(when time.Time) bool
```
true if when is a leap year


## func LoHi
``` go
func LoHi(ary []int) (lo int, hi int)
```
LoHi returns the min and max of an array of ints


## func LoadSHA256Names
``` go
func LoadSHA256Names(fname string) ([]string, error)
```
LoadSHA256Names returns the array of pathnames


## func LoadSHA256asDirMap
``` go
func LoadSHA256asDirMap(fname string) (map[string]*DirNode, error)
```
LoadSHA256asDirMap returns *dirMap[directoryPath]*DirNode


	uses scanner to allow large file.256 to be used as input

Used when its necessary to compute something with directory as input


## func LoadSHA256asList
``` go
func LoadSHA256asList(fname string) ([]string, error)
```
LoadSHA256asList returns lines of file as []string , nil on success


	each line is known to split without error

comments are stripped


## func LoadSHA256asMap
``` go
func LoadSHA256asMap(fname string) (map[string]string, error)
```
LoadSHA256asMap returns map[SHA256]path , nil on success


## func MSBytesFromInt64
``` go
func MSBytesFromInt64(n int64) []byte
```
convert int64 to [0:8]byte slice with MSB first (BigEndian)
reverse function is
test is Test_013


## func MSBytesFromUint32
``` go
func MSBytesFromUint32(u uint32) []byte
```
convert uint32 to [0:4]byte slice in MSB first (aka 'Net') order
reverse function is ThirtyTwoNet()


## func MaxI
``` go
func MaxI(a, b int) int
```

## func MinI
``` go
func MinI(a, b int) int
```

## func PermutedInts
``` go
func PermutedInts(a Ints) []Ints
```
PermutedInts returns all the permutaions of the original array
length of array must be in range of [0..6]


## func Polar
``` go
func Polar(loc Point) (theta, r float64)
```
Polar returns the polar coords (angle and radius) for a Cartesian point.


## func PolarAngle
``` go
func PolarAngle(loc Point) float64
```
PolarAngle returns the angle produced by a line from Point{0,0} to loc.


## func PseudoRandomBlock
``` go
func PseudoRandomBlock(blksize int) []byte
```
returns buffer of specified size filled with random bytes
test funcs are  Test_010 and Benchmark_010


## func Radians
``` go
func Radians(degrees float64) float64
```
radians returns the equivalent to degrees.


## func RangeLoHi64
``` go
func RangeLoHi64(v []float64) (lo, hi float64)
```
Return the lo and hi of given array


## func RangeMinMaxPoint
``` go
func RangeMinMaxPoint(v Points) (minPt, maxPt Point)
```
RangeMinMaxPoint returns bounds of point array.
probably a better name would help... 8888


## func Reverse
``` go
func Reverse(b []byte)
```
reverse a slice of bytes in place
Test with mdr_test.go:Test_006


## func SetTtyRawMode
``` go
func SetTtyRawMode()
```
this works, but see also go.crypto/ssh/terminal


## func SingleCharRead
``` go
func SingleCharRead() byte
```
singleCharRead might return newline if no other input


	otherwise first character if more than one on a line

^D will cause EOF to be printed and ? as the char returned - better choice is ...


## func SixteenBit
``` go
func SixteenBit(n []byte) uint16
```
==================   16 bit functions =====================
convert from little endian two byte slice to int16


## func Spinner
``` go
func Spinner()
```
keep user entertained while something happens behind the curtain


	see example from mdr_test.go for usage
	Choose Spinner() if progress bar can't be used because endpoint not known


## func StarDate
``` go
func StarDate(when time.Time) float64
```
returns decimal year to nearest 52 minutes - usually printed with %9.4f


## func ThirtyTwoBit
``` go
func ThirtyTwoBit(n []byte) uint32
```
convert from LITTLE endian four byte slice to int32
reverse function is LSBytesFromUint32


## func ThirtyTwoNet
``` go
func ThirtyTwoNet(n []byte) uint32
```
ThirtyTwoNet is a synonym for Uint32FromMSBytes


## func UdevRandomBlock
``` go
func UdevRandomBlock(blksize int) []byte
```
returns a buffer of specified size filled with random bits from /dev/urandom
there is a limited supply (about 4096 bits) so read will fail if urandom has not had time
to accumulate enough bits to satisfy the read.
Max readable under optimum conditios ?  actual test => 1020,1932,1020
1 ms sleep not enough(1940fail), 10 not enough (1020fail), 50 not enough 1020fail


## func Uint32FromIP
``` go
func Uint32FromIP(ip net.IP) uint32
```
returns a uint32 from IPv4 so we can use as index to map


	- beware - there are magic numbers here that
	- presume knowledge of net.IP internals order


## func Uint32FromMSBytes
``` go
func Uint32FromMSBytes(b []byte) uint32
```
convert from BIG endian four byte slice to int32
reverse function is MSBytesFromUint32


## func ValidDate
``` go
func ValidDate(year, month, day, hour, minute, second int) bool
```

## func ValidDecChar
``` go
func ValidDecChar(c byte) bool
```

## func ValidDecString
``` go
func ValidDecString(s string) bool
```

## func ValidHexChar
``` go
func ValidHexChar(c byte) bool
```
Test_004


## func ValidHexString
``` go
func ValidHexString(s string) bool
```
Test_004



## type ArchiveType
``` go
type ArchiveType int
```


``` go
const (
    ArchiveNoMatchType ArchiveType = iota
    ArchiveTarType
    ArchiveZipType
    ArchiveArkType
    ArchiveCpioType
)
```






### func WhichArchiveType
``` go
func WhichArchiveType(s string) ArchiveType
```
examine name to determine archive method used
should match magic numbers but that's a different function




## type ByName
``` go
type ByName []FileRec
```










### func (ByName) Len
``` go
func (a ByName) Len() int
```


### func (ByName) Less
``` go
func (a ByName) Less(i, j int) bool
```


### func (ByName) Swap
``` go
func (a ByName) Swap(i, j int)
```


## type CompressType
``` go
type CompressType int
```


``` go
const (
    CompressNoMatchType CompressType = iota
    CompressZipType
    CompressGzipType
    CompressBz2Type
    CompressZcompressType // .Z not common except possibly in Japan
    CompressPackType      // .z deprecated - very rare in last 20 years
    CompressBz1Type       // .bz deprecated - very rare in last 20+ years
)
```






### func WhichCompressType
``` go
func WhichCompressType(s string) CompressType
```
examine name to determine compression method used
should match magic numbers but that's a different function




## type DirNode
``` go
type DirNode struct {
    Pathname       string
    IsSortedByName bool
    Files          []FileRec
    Size           int64
    SHA256         string
}
```










### func (\*DirNode) AddFile
``` go
func (d *DirNode) AddFile(f FileRec)
```


### func (\*DirNode) DeleteFile
``` go
func (d *DirNode) DeleteFile(fname string)
```


### func (\*DirNode) Dump
``` go
func (d *DirNode) Dump()
```


### func (\*DirNode) Finalize
``` go
func (d *DirNode) Finalize() (string, int64)
```
Finalize computes the total size and SHA256 of the leaf



### func (\*DirNode) IndexOf
``` go
func (d *DirNode) IndexOf(fname string) int
```
special case speedup if we check last returned value +1 first?
better handles case where we're stepping through list ???
Returns -1 if not found



### func (\*DirNode) IndexOfBrute
``` go
func (d *DirNode) IndexOfBrute(fname string) int
```
brute force works



### func (\*DirNode) IsSortedProperly
``` go
func (d *DirNode) IsSortedProperly() bool
```


### func (\*DirNode) SortFilesByName
``` go
func (d *DirNode) SortFilesByName()
```


### func (\*DirNode) UnTouchAll
``` go
func (d *DirNode) UnTouchAll()
```


### func (\*DirNode) Update
``` go
func (d *DirNode) Update(fname string) error
```
calls AddFile if file not currently in DirNode
name should be just the base part of path



## type FileRec
``` go
type FileRec struct {
    R Rec256
    //	fname   string // needed to sort on
    Touched bool
}
```










### func (FileRec) Dump
``` go
func (fr FileRec) Dump()
```


## type IntPair
``` go
type IntPair struct {
    X, Y int
}
```










## type Ints
``` go
type Ints []int
```










### func (Ints) ContainsI
``` go
func (a Ints) ContainsI(b int) bool
```


### func (Ints) RotH2T
``` go
func (a Ints) RotH2T() Ints
```
RotT2H rotates the head of an array to tail position


	abcd => bcda



### func (Ints) RotT2H
``` go
func (a Ints) RotT2H() Ints
```
RotT2H rotates the tail of an array to head position


	abcd => dabc



## type Point
``` go
type Point struct {
    X int
    Y int
}
```








### func Cartesian
``` go
func Cartesian(angle, radius float64) Point
```
Cartesian returns Cartesian point from polar point.
NoteBene: NOT the usual [in trig x=cos(a) y=sin(a)]
because y axis is inverted in 'conventional' computer graphics.




### func (\*Point) RotateOnPivotPt
``` go
func (loc *Point) RotateOnPivotPt(p Point, radians float64)
```
RotateOnPivotPt rotate a point using another point as a pivot (mimic mechanical compass drawing).



### func (\*Point) ZRotate
``` go
func (loc *Point) ZRotate(radians float64)
```
ZRotate rotates a point about 0,0 (z axis).



## type Points
``` go
type Points []Point
```








### func CreateBezierPts
``` go
func CreateBezierPts(p1, p2, p3 Point, segments int) Points
```
CreateBezierPts start, control, end pts for quadratic bezier.  segments is
the number of line segments to create, more means smoother curve.




### func (Points) Split
``` go
func (pts Points) Split() (xs []int, ys []int)
```
Split separates the X and Y values into their own arrays



## type ProgStateT
``` go
type ProgStateT struct {
    // contains filtered or unexported fields
}
```








### func OneProgressBar
``` go
func OneProgressBar(goal int64) *ProgStateT
```



### func (\*ProgStateT) Stop
``` go
func (ps *ProgStateT) Stop()
```


### func (\*ProgStateT) Tag
``` go
func (ps *ProgStateT) Tag(t string)
```


### func (\*ProgStateT) Update
``` go
func (ps *ProgStateT) Update(val int64)
```


## type Rec256
``` go
type Rec256 struct {
    Size            int64
    SHA, Date, Name string
}
```








### func Split256
``` go
func Split256(line string) (Rec256, error)
```



### func (Rec256) Dump
``` go
func (r Rec256) Dump()
```


## type VerboseType
``` go
type VerboseType bool
```




``` go
var (
    Verbose VerboseType
)
```






### func (VerboseType) Printf
``` go
func (v VerboseType) Printf(s string, a ...interface{})
```








- - -
DO NOT EDIT ABOVE - Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)

github.com/hotei/mdr imports 
```
	C
	bufio
	bytes
	crypto
	crypto/md5
	crypto/sha256
	errors
	flag
	fmt
	github.com/hotei/datatable
	hash
	hash/crc64
	internal/singleflight
	io
	io/ioutil
	log
	math
	math/rand
	net
	os
	os/exec
	path/filepath
	reflect
	regexp
	regexp/syntax
	runtime
	sort
	strconv
	strings
	sync
	sync/atomic
	syscall
	time
	unicode
	unicode/utf8
	unsafe
```
```
deadcode results:

deadcode: mdr_DirNode.go:15:1: paranoid is unused
```
