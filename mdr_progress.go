package mdr

import (
	"fmt"
	"os"
)

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

// keep user entertained while something happens behind the curtain
//  see example from mdr_test.go for usage
//  see also Spinner() if endpoint of progress is unknown/unknowable
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
