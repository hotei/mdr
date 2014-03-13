package mdr

import (
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"hash"
	"hash/crc64"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	G_crcTable *crc64.Table
)

func init() {
	G_crcTable = crc64.MakeTable(crc64.ECMA)
}

func DigestToString(digest []byte) string {
	rv := ""
	for _, hex := range digest {
		rv = rv + fmt.Sprintf("%02x", hex)
	}
	return string(rv)
}

// syntax sugar
func BufCRC64(bufr []byte) uint64 {
	return crc64.Checksum(bufr[:], G_crcTable)
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

var CantCreateRec = errors.New("mdr: cant create Rec256 ")

type Rec256 struct {
	Size            int64
	SHA, Date, Name string
}

func (r Rec256) Dump() {
	fmt.Printf("%12d | %64s | %s | %s \n",
		r.Size,
		r.SHA,
		r.Date,
		r.Name)
}

func Split256(line string) (Rec256, error) {
	var rec Rec256
	var err error
	x := strings.Split(line, "|")
	if len(x) != 4 {
		fmt.Printf("!ERR--->  found other than 4 parts during split of %q\n", line)
		return rec, CantCreateRec
	}
	rec.Size, err = strconv.ParseInt(strings.Trim(x[0], " \n\t\r"), 10, 64)
	if err != nil {
		fmt.Printf("!ERR--->  Cant parse to int64: %s\n", x[0])
		return rec, CantCreateRec
	}
	rec.SHA = strings.Trim(x[1], " \n\t\r")
	if !ValidHexString(rec.SHA) {
		fmt.Printf("!ERR--->  %q is not valid hex \n", rec.SHA)
		return rec, CantCreateRec
	}
	rec.Date = strings.Trim(x[2], " \n\t\r")
	rec.Name = strings.Trim(x[3], " \n\t\r")

	return rec, nil
}
