package mdr

import (
	"fmt"
	"log"
	"os"
	"time"
)

// these would normally be C static vars
var (
	spinCt   int8
	lastCall time.Time
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
	//     "         1                   2                   3                   4                   5         6"
	//     "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
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

var progStates []ProgState

type ProgState struct {
	id         int
	lastBarlen int64
	val        int64
	goal       int64
	tag        string
}

var barWidth = int64(100)

// must produce output at least 1 per second even if flooded with input
// tokenbucket ? or
/*
	now := time.Now()
	if now.Sub(lastCall) < (time.Millisecond * 100) {
		return
	}
	lastCall = now


func progPrinter() {
	for {
		for i := 0; i < len(progStates); i++ {
			newlen := (progStates[i].val * int64(barWidth)) / progStates[i].goal
			fmt.Printf("newlen[%d] = %d\n", i, newlen)
			if newlen != progStates[i].lastBarlen {
				newlen = progStates[i].lastBarlen
			}
		}
		// print the bars
		for i := 0; i < len(progStates); i++ {
			fmt.Printf("id[%d] = %d\n", i, progStates[i].val)
		}
		time.Sleep(time.Millisecond * 1000)
	}
}
*/

func progUpdater() {
	fmt.Fprintf(os.Stderr, "Progress bars:\n")
	//     "                                                                                                   1"
	//     "         1         2         3         4         5         6         7         8         9         0"
	//     "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	bar := "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
	nobar := "....................................................................................................."
	//barlen := 0
	var newlen int64
	if barWidth > int64(len(bar)) {
		barWidth = int64(len(bar))
	}
	fmt.Fprintf(os.Stderr, "%s\r", nobar)
	for {
		progLen := len(progStates)
		for i := 0; i < progLen; i++ {
			val := progStates[i].val
			// if val < 0  then delete that entry - if all deleted then return
			if val < 0 {
				if progLen <= 1 {
					return
				} else { // swap current and last states
					progStates[i] = progStates[progLen-1]
					// trim the last one off
					progStates = progStates[:progLen-1]
					progLen--
					fmt.Fprintf(os.Stderr, "\n") // only right if one bar
					continue
				}
			}
			goal := progStates[i].goal
			newlen = (val * int64(barWidth)) / goal
			s := bar[:newlen] + nobar[:barWidth-newlen]
			fmt.Fprintf(os.Stderr, "%s %s\r", s, progStates[i].tag)
			if newlen != progStates[i].lastBarlen {
				progStates[i].lastBarlen = newlen
			}
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func (ps *ProgState) Update(val int64) {
	progStates[ps.id].val = val
}

func (ps *ProgState) Stop() {
	progStates[ps.id].val = -1
}

func (ps *ProgState) Tag(t string) {
	progStates[ps.id].tag = t
}

func NewProgressBar(goal int64) *ProgState {
	if len(progStates) <= 0 {
		go progUpdater()
	}
	if len(progStates) >= 1 {
		return nil
	}
	var p ProgState
	p.id = len(progStates)
	p.goal = goal
	progStates = append(progStates, p)
	fmt.Printf("Created bar[%d]\n", p.id)
	return &p
}
