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
const (
	UpdateDelayMs = 200 // delay in millisec between updates
)

var (
	// these would normally be C static vars inside Spinner()
	spinCt   int8
	lastCall time.Time
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
	label      string
	lastBarlen int64
	val        int64
	goal       int64
	endtag     string
}

var barWidth = int64(100)

func progUpdater(progState *ProgStateT) {
	fmt.Fprintf(os.Stderr, "Progress bars:\n")
	//     "                                                                                                   1"
	//     "         1         2         3         4         5         6         7         8         9         0"
	//     "1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	bar := "+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++"
	nobar := "....................................................................................................."

	var newlen int64
	if barWidth > int64(len(bar)) {
		barWidth = int64(len(bar))
	}
	fmt.Fprintf(os.Stderr, "%s\r", nobar)
	for {
		val := progState.val
		if val < 0 { // found stop signal
			return
		}
		goal := progState.goal
		newlen = (val * int64(barWidth)) / goal
		s := bar[:newlen] + nobar[:barWidth-newlen]
		fmt.Fprintf(os.Stderr, "%s: %s %s\r", progState.label, s, progState.endtag)
		if newlen != progState.lastBarlen {
			progState.lastBarlen = newlen
		}

		time.Sleep(time.Millisecond * UpdateDelayMs)
	}
}

func (ps *ProgStateT) Update(val int64) {
	ps.val = val
}

func (ps *ProgStateT) Stop() {
	time.Sleep(time.Second)
	fmt.Fprintf(os.Stderr, "\nstopping %s: %s\n", ps.label, ps.endtag)
	ps.val = -1
}

func (ps *ProgStateT) Tag(t string) {
	t = strings.Trim(t, "\n\r\t ")
	ps.endtag = t
}

func OneProgressBar(goal int64) *ProgStateT {
	var p ProgStateT
	p.goal = goal
	p.label = "pbar"
	fmt.Printf("Created bar as: %v\n", p)
	go progUpdater(&p)
	return &p
}
