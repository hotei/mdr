package mdr

import (
	"fmt"
	"log"
	"os"
	"time"
)

// these would normally be C static vars
var (
	spinCt int8
	lastCall time.time
	)
const spinchars string = "|/-\\|/-\\ "

// keep user entertained while something happens behind the curtain
//  see example from mdr_test.go for usage
//  Choose this if progress bar can't be used because endpoint not known
func Spinner() {
	now := time.Now()
	if now.Sub(lastCall) < (time.Millisecond * 100) { 
		return	
	}
	lastCall = now
	fmt.Fprintf(os.Stderr, "%s\r", spinchars[spinCt:spinCt+1])
	spinCt++
	spinCt &= 0x7 // mod 8 which is length of spinchars by design
}

// keep user entertained while something happens behind the curtain
//  see example from mdr_test.go for usage
//  see also Spinner() if endpoint of progress is unknown/unknowable
func ProgressBar(barWidth int, p chan int64, alldone int64) {
	if alldone <= 0 {
		log.Fatalf("nothing to do - alldone[%d] is <= 0\n", alldone)
	}
	fmt.Fprintf(os.Stderr, "Progress bar:\n")

	//"1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	bar := "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
	nobar := "....................................................................................................."
	barlen := 0
	if barWidth > len(bar) {
		barWidth = len(bar)
	}
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
