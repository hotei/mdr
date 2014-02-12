// mdr_cryptoHelp.go (c) 2012-2014 David Rook - see LICENSE.md

package mdr

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

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

// this works, but see also go.crypto/ssh/terminal
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
