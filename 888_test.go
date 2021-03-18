// 888_test.go

// go test '-bench=.'   # expression matches everything, runs all benchmarks
// go test -run="Test_000" to run just one function

package mdr

import (
	"fmt"
	//"testing"
	//"time"
)

/////////////////////////  E X A M P L E S  ////////////////////////////
func workerFunction(w IntPair) {
	fmt.Printf("work on items from %d through %d\n", w.X, w.Y)
}

func ExampleJobSplit() {
	nCPU := 4
	totalWork := 100
	jobrange := JobSplit(totalWork, nCPU)
	for i := 0; i < nCPU; i++ {
		go workerFunction(jobrange[i])
	}
}

/*
func ExampleProgressBar() {
	var (
		status    int64
		endNumber int64 = 100
	)

	progChan := make(chan int64, 2)
	go ProgressBar(50, progChan, endNumber) // start the display handler
	progChan <- 0                           // make first progress display visible
	for {
		//      ... do something to advance status towards endNumber ...
		time.Sleep(time.Second)
		status += 10
		progChan <- status
		if status >= endNumber {
			break
		}
	}
	progChan <- -1 // close up shop
}
*/
