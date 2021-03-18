// mdr_progress.go (c) 2012-2015 David Rook - released under Simplified BSD License

package mdr

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// keep user entertained while something happens behind the curtain
//  see example from mdr_test.go for usage
//  use Spinner() if endpoint of progress is unknown/unknowable

var (
	// these would normally be C static vars inside Spinner()
	UpdateDelayMs = 1000 // delay in millisec between updates
	spinCt        int8
	lastCall      time.Time
)

// order is important so preserve it
// we use &= instead of mod 8 so length of spinchars must be 8
const spinchars string = "|/-\\|/-\\ "

// keep user entertained while something happens behind the curtain
//  see example from mdr_test.go for usage
//  Choose Spinner() if progress bar can't be used because endpoint not known
func Spinner() {
	// do not update spinner faster than 10 times per second
	now := time.Now()
	if now.Sub(lastCall) < (time.Millisecond * 100) {
		return
	}
	lastCall = now
	fmt.Fprintf(os.Stderr, "%s\r", spinchars[spinCt:spinCt+1])
	spinCt++
	spinCt &= 0x7
}

type ProgStateT struct {
	id         int
	Label      string
	lastBarlen int64
	val        int64
	goal       int64
	endtag     string
	startedAt  time.Time
	ShowCount  bool
	ShowTime   bool
}

var barWidth = int64(100)

func progUpdater(progState *ProgStateT) {
	//fmt.Fprintf(os.Stderr, "Progress bars:\n")
	//     "                                                                                                   1"
	//     "         1         2         3         4         5         6         7         8         9         0"
	//     "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	//mybar := "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
	//nobar := "....................................................................................................."
	//nobar  = "-----------------------------------------------------------------------------------------------------"
	nobar := strings.Repeat("o", 100)
	mybar := strings.Repeat("+", 100)
	var (
		newlen  int64
		longest int
	)
	if barWidth > int64(len(mybar)) {
		barWidth = int64(len(mybar))
	}
	fmt.Fprintf(os.Stderr, "%s\r", nobar)
	for {
		val := progState.val
		if val < 0 { // found stop signal
			return
		}
		goal := progState.goal
		if goal <= 0 {
			return
		}
		if val > goal { // don't exceed 100% completion
			val = goal
		}
		newlen = (val * int64(barWidth)) / goal
		s := mybar[:newlen] + nobar[:barWidth-newlen]
		unitsLeft := goal - val
		out := fmt.Sprintf("%s %s %s ", progState.Label, s, progState.endtag)
		if progState.ShowCount {
			out += fmt.Sprintf("%9d items left => ", unitsLeft)
		}
		if progState.ShowTime {
			elapsed := time.Since(progState.startedAt)
			elapsedSec := elapsed.Seconds()
			//fmt.Printf("elapsedSec = %g\n",elapsedSec)
			currentRate := float64(progState.val) / elapsedSec // elapsed sec / items done so far
			secToGo := float64(unitsLeft) / currentRate
			out += fmt.Sprintf("%5d seconds to go ", int(secToGo))
		}
		if len(out) > longest {
			longest = len(out)
		}
		fmt.Fprintf(os.Stderr, "%s\r", strings.Repeat(" ", longest))
		fmt.Fprintf(os.Stderr, "%s\r", out)

		if newlen != progState.lastBarlen {
			progState.lastBarlen = newlen
		}

		time.Sleep(time.Millisecond * time.Duration(UpdateDelayMs))
	}
}

func (ps *ProgStateT) Update(val int64) {
	ps.val = val
}

func (ps *ProgStateT) Stop() {
	time.Sleep(time.Second)
	if Verbose {
		fmt.Fprintf(os.Stderr, "\nstopping %s: %s\n\n", ps.Label, ps.endtag)
	}
	ps.val = -1
}

func (ps *ProgStateT) Tag(t string) {
	t = strings.Trim(t, "\n\r\t ")
	ps.endtag = t
}

func OneProgressBar(goal int64) *ProgStateT {
	var p ProgStateT
	p.goal = goal
	p.startedAt = time.Now()
	go progUpdater(&p)
	return &p
}
