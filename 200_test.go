// 200_test.go

package mdr

/*  maybe make a table of reverse data- with right-left reversed.
then we can still use interp even with the recent bounds checking
*/
import (
	"fmt"
	//	"errors"
	"testing"
)

func Test_200(t *testing.T) {
	fmt.Printf("Test_200 template\n")
	if false {
		t.Errorf("print fail, but keep testing")
	}
	if false {
		t.Fatalf("print fail and keep testing")
	}
	Verbose.Printf("Pass - test_200\n")
}

func Test_201(t *testing.T) {
	fmt.Printf("Test_201 tables\n")

	mytable := new(Table)
	mytable.Name = "test table"
	mytable.Data = []DblPair{
		{0.0, 1737.0},
		{100.0, 1384.0},
		{200.0, 1092.0},
		{300.0, 800.0},
		{400.0, 600.0},
		{500.0, 497.0},
		{900.0, 0.0},
	}
	var rv float64
	var err error
	err = mytable.Setup()
	if err != nil {
		t.Errorf("Eval() failed with %v\n", err)
	}

	ForceRangeFlag = false
	if Verbose {
		mytable.Dump()
	}
	rv, err = mytable.Eval(0.0)
	if err != nil {
		t.Errorf("Eval() failed with %v\n", err)
		return
	}
	if rv != 1737.0 {
		t.Errorf("failed for lowest DblPair")
	}

	rv, err = mytable.Eval(900.0)
	if err != nil {
		t.Errorf("Eval() failed with %v\n", err)
		return
	}
	if rv != 0.0 {
		t.Errorf("failed for highest DblPair")
	}

	rv, err = mytable.Eval(200)
	if err != nil {
		t.Errorf("Eval() failed with %v\n", err)
		return
	}
	if rv != 1092.0 {
		t.Errorf("failed for match one end of a DblPair")
	}

	rv, err = mytable.Eval(350)
	if err != nil {
		t.Errorf("Eval() failed with %v\n", err)
		return
	}
	if rv != 700.0 {
		t.Errorf("failed for middle of a DblPair")
	}

	_, err = mytable.Eval(-350)
	if err != OutOfRangeError {
		t.Errorf("Eval() failed to detect out of range %v\n", err)
		return
	}

	_, err = mytable.Eval(3500)
	if err != OutOfRangeError {
		t.Errorf("Eval() failed to detect out of range %v\n", err)
		return
	}

	ForceRangeFlag = true
	rv, err = mytable.Eval(-35)
	if err != nil {
		t.Errorf("Eval() failed with %v\n", err)
		return
	}
	if rv != 1737 {
		t.Errorf("failed to force when outside the low end of all DblPairs")
	}

	rv, err = mytable.Eval(3500)
	if err != nil {
		t.Errorf("Eval() failed with %v\n", err)
		return
	}
	if rv != 0.0 {
		t.Errorf("failed to force when outside the high end of all DblPairs")
	}
	Verbose.Printf("Pass - test_201\n")
}

func Test_202(t *testing.T) {
	fmt.Printf("Test_202 table setup \n")
	Verbose.Printf("table.Setup() should get errors trying to use invalid tables \n")
	mytable := new(Table)
	mytable.Name = "test table"
	var err error

	//////////////////
	mytable.Data = []DblPair{
		{0.0, 1737.0},
		{100.0, 1384.0},
		{100.0, 1092.0},
		{300.0, 800.0},
		{400.0, 600.0},
		{500.0, 497.0},
		{900.0, 0.0},
	}
	err = mytable.Setup()
	if err != InvalidTableError {
		t.Errorf("Setup() failed to detect an invalid (duplicate line) table %v\n", err)
	}
	//////////////////
	mytable.Data = []DblPair{
		{0.0, 1737.0},
	}
	err = mytable.Setup()
	if err == nil {
		t.Errorf("Setup() failed to detect an invalid (short) table %v\n", err)
	}
	//////////////////
	mytable.Data = []DblPair{
		{100.0, 1384.0},
		{0.0, 1737.0},
	}
	err = mytable.Setup()
	if err == nil {
		t.Errorf("Setup() failed to detect an invalid (incorrectly ordered) table %v\n", err)
	}

	Verbose.Printf("Pass - test_202\n")

}

func Test_203(t *testing.T) {
	fmt.Printf("Test_203 table ReverseEval \n")
	mytable := new(Table)
	mytable.Name = "test table"
	var rv []float64
	var err error

	//Verbose = VerboseType(true)

	//////////////////
	mytable.Data = []DblPair{
		{0.0, 1737.0},
		{100.0, 1384.0},
		{200.0, 1092.0},
		{300.0, 800.0},
		{400.0, 600.0},
		{500.0, 500.0},
		{1000.0, 0.0},
	}
	err = mytable.Setup()
	if err != nil {
		t.Errorf("Setup() failed with %v\n", err)
	}

	rv, err = mytable.ReverseEval(1737.0)
	if err != nil {
		t.Errorf("ReverseEval() failed with %v\n", err)
		return
	}
	if len(rv) != 1 {
		t.Errorf("ReverseEval() first failed, return array has wrong length")
	}
	if rv[0] != 0.0 {
		t.Errorf("ReverseEval() first failed expected [0.0] got %v", rv)
	}
	//////////////////
	rv, err = mytable.ReverseEval(0.0)
	if err != nil {
		t.Errorf("eval failed with %v\n", err)
		return
	}
	if len(rv) != 1 {
		t.Errorf("ReverseEval() first failed, return array has wrong length")
	}
	if rv[0] != 1000.0 {
		t.Errorf("ReverseEval() first failed expected [1000.0] got %v", rv)
	}
	//////////////////
	rv, err = mytable.ReverseEval(1800.0)
	if err != OutOfRangeError {
		t.Errorf("eval failed to detect out of range (high) with %v\n", err)
		return
	}
	//////////////////
	rv, err = mytable.ReverseEval(-1800.0)
	if err != OutOfRangeError {
		t.Errorf("eval failed to detect out of range (low) with %v\n", err)
		return
	}
	//////////////////
	rv, err = mytable.ReverseEval(250.0)
	if err != nil {
		t.Errorf("eval failed with %v\n", err)
		return
	}
	if len(rv) != 1 {
		t.Errorf("ReverseEval() first failed, return array has wrong length")
	}
	if rv[0] != 750.0 {
		t.Errorf("ReverseEval() first failed expected [750.0] got %v", rv)
	}
	//////////////////
	rv, err = mytable.ReverseEval(500.0)
	if err != nil {
		t.Errorf("eval failed with %v\n", err)
		return
	}
	if len(rv) != 1 {
		t.Errorf("ReverseEval() first failed, return array has wrong length")
	}
	if rv[0] != 500.0 {
		t.Errorf("ReverseEval() first failed expected [500.0] got %v", rv)
	}
	//////////////////
	rv, err = mytable.ReverseEval(0.0)
	if err != nil {
		t.Errorf("eval failed with %v\n", err)
		return
	}
	if rv[0] != 1000.00 {
		t.Errorf("ReverseEval() first failed, expected [1000.0] got %v", rv)
	}
	//////////////////
	rv, err = mytable.ReverseEval(250)
	if err != nil {
		t.Errorf("eval failed with %v\n", err)
		return
	}
	if rv[0] != 750.00 {
		t.Errorf("ReverseEval() first failed, expected [750.0] got %v", rv)
	}

	Verbose.Printf("Pass - test_203\n")

}
