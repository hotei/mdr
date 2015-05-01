// progTest.go

package main

import (
	"fmt"
	"os"
	"time"

	// local pkgs
	"github.com/hotei/mdr"
)

func main() {
	fmt.Printf("Test_Progress \n")
	goal := int64(600) //
	barA := mdr.OneProgressBar(goal)
	for i := int64(0); i < goal; i++ {
		// note that bar doesn't update for every loop, just every 200 ms
		barA.Update(i)
		barA.Tag(fmt.Sprintf("%d of %d have been tested\n", i, goal))
		time.Sleep(time.Millisecond * 10)
	}
	barA.Update(goal)
	barA.Tag(fmt.Sprintf("%d of %d have been tested", goal, goal))
	barA.Stop()

	goal = 300
	barB := mdr.OneProgressBar(goal)
	for i := int64(0); i < goal; i++ {
		// note that bar doesn't update for every loop, just every 200 ms
		barB.Update(i)
		barB.Tag(fmt.Sprintf("%d of %d have been done", i, goal))
		time.Sleep(time.Millisecond * 10)
	}
	barB.Update(goal)
	barB.Tag(fmt.Sprintf("%d of %d have been done\n", goal, goal))
	barB.Stop()
	fmt.Printf("Pass - Test_Progress()\n")
	if false {
		os.Exit(0)
	}

}
