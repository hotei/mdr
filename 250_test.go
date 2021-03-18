// sutil_test.go (c) 2019 David Rook all rights reserved

/*
Useful at times...

		t.Fatalf("this test failed,  run next test func now")
		t.Errorf("this test failed, but keep testing this func")

		log.Printf(" is not a valid regexp, err: %s", err.Error())
		log.Fatalf(" is not a valid regexp, err: %s", err.Error())

		fmt.Printf(`go test -bench="*." to run all benchmarks\n`)
		fmt.Printf(`to run a single test E use go test -test.run="Test_E"\n`)

*/

package mdr

import (
	"fmt"
	//"log"
	"testing"
)

// this is a template for testing
func Test_250(t *testing.T) {
	fmt.Printf("Test_250 template \n")
	errCt := 0

	// main body of tests go here
	computed := "target"
	expected := "target"
	if computed != expected {
		errCt++
		t.Errorf("Test_250() failed, expected %s got %s", expected, computed)
	}

	// next is a trivial example of "table-driven" testing
	// the idea is to exercise a function - len() in this case - over a range of
	// values and compare the result to an expected value
	type testRack struct {
		given  string
		expect bool
	}

	var targets []testRack = []testRack{
		{
			given:  "above",
			expect: true,
		}, {
			given:  "beyond",
			expect: false,
		},
	}

	for _, val := range targets {
		result := len(val.given) <= 5
		if result != val.expect {
			t.Errorf("Test_00() len() failed here. Expected %v, got %v",
				val.expect, result)
			errCt++
		}
	}
}

func Test_251(t *testing.T) {
	fmt.Printf("Test_251 template \n")
	errCt := 0

	// main body of tests go here

	type testRack struct {
		given  string
		expect bool
	}

	var targets []testRack = []testRack{
		{
			given:  "above",
			expect: true,
		}, {
			given:  "beyond",
			expect: false,
		},
	}

	for _, val := range targets {
		result := len(val.given) <= 5
		if result != val.expect {
			t.Errorf("Test_251() len() failed here. Expected %v, got %v",
				val.expect, result)
			errCt++
		}
	}
}
