package mdr

import (
	"github.com/hotei/datatable"
	// also uses standard lib go 1.2 pkgs below
	"fmt"
	"math"
	"math/rand"
)

var NormalZtable *datatable.Table

func init() {
	// table contains : z value, area (ie. probability) to left of z value
	// upper half of table only, negative values = (1.0 - upper value)
	NormalZtable = new(datatable.Table)
	NormalZtable.Name = "normal_Z_table"
	NormalZtable.Data =
		[]datatable.DblPair{{0.0, 0.5000},
			{0.25, 0.5987},
			{0.5, 0.6915},
			{0.75, 0.7734},
			{1.0, 0.8413},
			{1.5, 0.9332},
			{2.0, 0.9722},
			{2.5, 0.9938},
			{3.0, 0.9987},
			{3.49, 0.9998},
			{4.00, 0.99999},
			{5.00, 0.999995},
			{6.00, 0.9999966},
			{60.00, 1.00},
		}
}

// return float64 in range {0 .. r}
func GenRandomUniform(r float64) float64 {
	return rand.Float64() * r
}

func GenRandomUniformLoHi(low, high float64) float64 {
	r := high - low
	return low + GenRandomUniform(r)
}

func GenRandomZNormal() float64 {
	rnd := 0.5 + (rand.Float64() / 2.0)
	rv, err := NormalZtable.ReverseEval(rnd)
	if err != nil {
		Crash(fmt.Sprintf("reverse eval of normalZtable failed with err %v", err))
	}
	if len(rv) > 1 {
		fmt.Printf("got more than one return value\n")
	}
	if FlipCoin() {
		rv[0] = -rv[0]
	}
	return rv[0]
}

func GenRandomNormal(mu, stdev float64) float64 {
	rnd := GenRandomZNormal()
	dev := stdev * rnd
	if FlipCoin() {
		dev = -dev
	}
	rv := mu + dev
	return rv
}

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
