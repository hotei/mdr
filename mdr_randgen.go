// mdr_randgen.go (c) 2010-2020 David Rook - all rights reserved

package mdr

import (
	"fmt"
	"math"
	"math/rand"
)

var NormalZtable *Table

func init() {
	Verbose.Printf("mdr.randgen.go init() entry\n")
	defer Verbose.Printf("mdr.randgen.go init() exit\n")
	// table contains : z value, area (ie. probability) to left of z value
	// upper half of table only, negative values = (1.0 - upper value)
	// table values could be expanded if it makes sense to do so
	// the 60 value is a wag
	NormalZtable = new(Table)
	NormalZtable.Name = "normal_Z_table"
	NormalZtable.Data =
		[]DblPair{
			{0.0, 0.5000},
			{0.1, 0.5398},
			{0.2, 0.5793},
			{0.25, 0.5987},
			{0.5, 0.6915},
			{0.75, 0.7734},
			{1.0, 0.8413},
			{1.5, 0.9332},
			{2.0, 0.9722},
			{2.4, 0.9938},
			{3.0, 0.9987},
			{3.4, 0.9997},
			{3.49, 0.9998}, // last value in table IV of Walpole & Meyers "Prob & Stat for Engrs & Sci"
			{4.00, 0.99999},
			{5.00, 0.999999},
			{6.00, 0.9999999},
			{8.00, 0.999999999},
			{25.00, 1.00},
		}
}

// GenRandomZNormal returns a float64 with average of 0 and standard deviation of 1.0
// as implemented the range of values returned will be in [-60..60]
func GenRandomZNormal() float64 {
	if true {
		return rand.NormFloat64()
	} else {
		rnd := 0.5 + (rand.Float64() / 2.0) // should return [.5 .. 1.0]
		rv, err := NormalZtable.ReverseEval(rnd)
		if err != nil {
			Crash(fmt.Sprintf("reverse eval of normalZtable failed with err %v", err))
		}
		if len(rv) > 1 {
			Crash(fmt.Sprintf("got more than one return value\n"))
		}
		if FlipCoin() {
			rv[0] = -rv[0]
		}
		return rv[0]
	}
}

// RandIntBtw endpoints may occur (HasTest widget).
func GenRandIntBtw(lo, hi int) int {
	if lo > hi {
		lo, hi = hi, lo
	}
	dif := (hi - lo) + 1
	t := rand.Int31() % int32(dif)
	return int(t) + lo
}

// GenRandomNormal returns a float64 with average mu and standard deviation stdev
// depending on mu and stdev picked the values returned could be virtually any float64
func GenRandomNormal(mu, stdev float64) float64 {
	rnd := GenRandomZNormal()
	dev := stdev * rnd
	if FlipCoin() {
		dev = -dev
	}
	rv := mu + dev
	return rv
}

// GenRandomPoisson returns an int with Poisson distribution
// typically used to determine the number of time units before some event occurs
func GenRandomPoisson(lambda float64) int {
	L := 1.0 / math.Exp(lambda)
	k := 0
	p := 1.0
	for {
		k += 1
		u := rand.Float64()
		p *= u
		if p <= L {
			break
		}
	}
	return k
}

// expect to see about 50% head, 50% tails (HasTest)
func GenFlipCoin() bool {
	return rand.Int31n(2) == 0
}

// GenRandomUniform returns a float64 in range {0 .. r}
func GenRandomUniform(r float64) float64 {
	return rand.Float64() * r
}

// GenRandomUniformLo returns a float64 in range {low .. high}
// high > low is NOT required, not sure if this is right or should panic 8888 ?
func GenRandomUniformLoHi(low, high float64) float64 {
	if low > high {
		low, high = high, low
	}
	myrange := high - low
	return low + GenRandomUniform(myrange)
}

// GenRandIntBetween endpoints may occur in output
func GenRandIntBetween(lo, hi int) int {
	if lo > hi {
		lo, hi = hi, lo
	}
	dif := (hi - lo) + 1
	t := rand.Int31() % int32(dif)
	return int(t) + lo
}

// GenRandF64Between endpoints may occur in output
func GenRandF64Between(lo, hi float64) float64 {
	if lo > hi {
		lo, hi = hi, lo
	}
	dif := hi - lo
	return rand.Float64()*dif + lo
}
