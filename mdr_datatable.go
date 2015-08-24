// datatable.go  (c) 2011-2014 David Rook - all rights reserved

package mdr

import (
	"errors"
	"fmt"
	"log"
)


//  ----------------------------------------
//  |           d a t a t a b l e          |
//  ----------------------------------------
//

type DblPair struct {
	Left, Right float64
}

type Table struct {
	Name            string
	Data            []DblPair
	RevData         []DblPair
	high, low       float64
	revHigh, revLow float64
	isSetup         bool
}

var (
	ForceRangeFlag    bool = true // forces out of range values to fit
	CantHappenError        = errors.New("Can't-happen")
	DivideByZeroError      = errors.New("Divide By Zero")
	OutOfRangeError        = errors.New("Value outside expected range")
	InvalidTableError      = errors.New("Not a valid table")
)

func init() {
	Verbose = false
}

func (d *Table) Setup() error {
	if len(d.Name) <= 0 {
		fmt.Printf("table really should be named\n")
	}
	length := len(d.Data)
	if length < 2 {
		fmt.Printf("table not valid, length must be >= 2 but we found %d\n", length)
		return InvalidTableError
	}
	last := length - 1

	d.high = d.Data[0].Left
	d.low = d.high
	for i := 0; i < length; i++ {
		datum := d.Data[i].Left
		if datum > d.high {
			d.high = datum
		}
		if datum < d.low {
			d.low = datum
		}
		// table left values must increase each time in a valid table
		if i > 0 {
			if d.Data[i].Left <= d.Data[i-1].Left {
				fmt.Printf("table not valid, left values must increase with each row\n")
				return InvalidTableError
			}
		}
	}
	for i := last; i >= 0; i-- {
		d.RevData = append(d.RevData, DblPair{Left: d.Data[i].Right, Right: d.Data[i].Left})
	}
	d.revHigh = d.RevData[0].Left
	d.revLow = d.revHigh
	for i := 0; i < length; i++ {
		datum := d.RevData[i].Left
		if datum > d.revHigh {
			d.revHigh = datum
		}
		if datum < d.revLow {
			d.revLow = datum
		}
	}
	Verbose.Printf("RevData = %v\n", d.RevData)
	d.isSetup = true
	return nil
}

func (d *Table) Dump() {
	if !d.isSetup {
		log.Printf("Table %s was not setup, doing it now...\n", d.Name)
		d.Setup()
	}
	fmt.Printf("\nTable Name = %s\n", d.Name)
	fmt.Printf("Data = %v\n", d.Data)
	fmt.Printf("Reverse Data = %v\n", d.RevData)
	fmt.Printf("Data Length = %d\n\n", len(d.Data))
}

func inrange(a, b, c float64) bool {
	if a > c {
		a, c = c, a
	}
	if b < a {
		return false
	}
	if b > c {
		return false
	}
	return true
}

/*
 * x1,y1 and x2,y2 define a line
 * given v use linear interpolation to find its y value
 * x1 <= v <= x2
 */
func interp(x1, y1, x2, y2, v float64) (rc float64, err error) {
	Verbose.Printf("interp: x1[%v] y1[%v] x2[%v] y2[%v] v %v\n", x1, y1, x2, y2, v)
	if x2 == x1 {
		return 0.0, DivideByZeroError
	}
	rc = ((y2-y1)/(x2-x1))*(v-x1) + y1
	Verbose.Printf("interp returned %v %v\n", rc, err)
	return rc, nil
}

func (d *Table) Eval(v float64) (rc float64, err error) {
	if !d.isSetup {
		log.Printf("Table %s was not setup, doing it now...\n", d.Name)
		d.Setup()
	}
	length := len(d.Data)
	last := length - 1
	Verbose.Printf("%s : evaluate(%g) with %d items in lookup\n", d.Name, v, length)
	if !inrange(d.Data[0].Left, v, d.Data[last].Left) {
		if ForceRangeFlag {
			if v < d.Data[0].Left {
				Verbose.Printf(" low side is %g\n", d.Data[0].Left)
				Verbose.Printf("forced low end value match \n")
				return d.Data[0].Right, nil
			}

			if v > d.Data[last].Left {
				Verbose.Printf(" high side is %g\n", d.Data[last].Left)
				Verbose.Printf("forced high end value match \n")
				return d.Data[last].Right, nil
			}

		} else {
			return 0.0, OutOfRangeError
		}
	}

	Verbose.Printf(" low side is %g\n", d.Data[0].Left)
	if v == d.Data[0].Left {
		Verbose.Printf("extreme low end value match \n")
		return d.Data[0].Right, nil
	}

	Verbose.Printf(" high side is %g\n", d.Data[last].Left)
	if v == d.Data[last].Left {
		Verbose.Printf("extreme high end value match \n")
		return d.Data[last].Right, nil
	}

	n := 1 // work table from bottom up, starting with second row
	for {
		if n > last {
			break // cant happen
		}
		if d.Data[n].Left > v { // must be in this range
			rc, err = interp(d.Data[n-1].Left, d.Data[n-1].Right, d.Data[n].Left, d.Data[n].Right, v)
			if err == nil {
				return rc, nil
			}
			if (err == OutOfRangeError) && ForceRangeFlag {
				return rc, nil
			}
			return rc, err
		}
		n++
	}
	// since we've already tested high & low we should never get here, but need a return anyway
	fmt.Printf("Cant happen\n")
	return 0.0, CantHappenError // 0.0 and failed
}

// returns multiple hits where appropriate
func (d *Table) ReverseEval(val float64) (rc []float64, err error) {
	if !d.isSetup {
		log.Printf("Table %s was not setup, doing it now...\n", d.Name)
		d.Setup()
	}
	var rv float64
	rc = make([]float64, 0, 10)
	length := len(d.Data)
	last := length - 1

	Verbose.Printf(" low side is %g high side is %g\n", d.RevData[0].Left, d.RevData[last].Left)
	Verbose.Printf("d.RevLow %v d.RevHigh %v\n", d.revLow, d.revHigh)
	if !inrange(d.revLow, val, d.revHigh) {
		return nil, OutOfRangeError
	}
	n := last
	//    TODO:? out of range values force to endpoints bug/feature?
	// I can see an argument for NOT doing out of bounds stuff

	for n > 0 {
		Verbose.Printf("n = %d\n", n)
		if inrange(d.RevData[n-1].Left, val, d.RevData[n].Left) {
			rv, err = interp(d.RevData[n-1].Left, d.RevData[n-1].Right, d.RevData[n].Left, d.RevData[n].Right, val)
			if err == nil {
				rc = append(rc, rv)
			} else {
				return nil, err // divide by zero should never happen with valid table
			}
		}
		// keep looking
		n--
	}
	if len(rc) <= 0 {
		return nil, nil
	}
	// BUG(mdr) remove duplicates if necessary, these happen in adjacent spots when a value is on
	// the end of a pair, can happen elsewhere as well, but not treating those yet.  This method is
	// quick, but inadequate for long run.  Need insertion sort or map for final solution.
	// currently just document that the resulting array may not have unique values
	Verbose.Printf("rv before 'uniq' = %v\n", rc)
	var newrc = []float64{rc[0]}
	for i := 1; i < len(rc); i++ {
		if rc[i] != rc[i-1] {
			newrc = append(newrc, rc[i])
		}
	}
	return newrc, nil
}
