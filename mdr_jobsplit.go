package mdr

import (
	"fmt"
	"os"
)

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
