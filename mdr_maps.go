// mdr_maps.go

/*
// BUG(mdr): need some more work on TurnDirection()
func TurnDirection(a, b GPSloc) float64 {
	dx := a.Lat - b.Lat
	dy := a.Long - b.Long
	heading := math.Atan2(dy, dx) * DegPerRad
	heading += 180
	if heading >= 360 {
		heading -= 360
	}
	Verbose.Printf("turnDirection = %6.1f degrees \n", heading)
	return heading
}
*/

package mdr

import (
	"fmt"
	"math"
)

type HeadingDeg float64

type dirT struct {
	start float64
	dir   string
}

var DegPerRad = 180.0 / math.Pi

var dirList []dirT = []dirT{
	{0, "North"},
	{22.5 + 45.0*0.0, "NorthEast"},
	{22.5 + 45.0*1.0, "East"},
	{22.5 + 45.0*2.0, "SouthEast"},
	{22.5 + 45.0*3.0, "South"},
	{22.5 + 45.0*4.0, "SouthWest"},
	{22.5 + 45.0*5.0, "West"},
	{22.5 + 45.0*6.0, "NorthWest"},
	{22.5 + 45.0*7.0, "North"},
}

func (h HeadingDeg) String() string {
	if h < 0 {
		h += 360.0
	}
	rv := "North"
	for i := 0; i < len(dirList); i++ {
		if h < HeadingDeg(dirList[i].start) {
			return dirList[i].dir
		}
	}
	return rv
}

//=============================================================================  GPS2d

type GPS2dT struct {
	Lon float64 // degrees
	Lat float64 // degrees
}

func (g GPS2dT) Point() Pointe {
	var rv Pointe
	rv.X = g.Lon
	rv.Y = g.Lat
	return rv
}

func (g GPS2dT) ValidGPS2d() bool {
	return InRangeF64(-180.0, g.Lon, 180.0) && InRangeF64(-90.0, g.Lat, 90.0)
}

// KmBetweenGC calculates the shortest path between two GPS coordinates. Assumes
// both points are at sealevel.
//  NOTE! Earth is not a sphere. Haversine has more error in E-W calculations than N-S.
//
// Distance of 1 degree differs ( E.W. vs N.S. )
//  Error is 6:40,000 E.W. vs 0:40000 N.S.
//  https://en.wikipedia.org/wiki/Earth_radius source for GaiaKm
//
// Uses haversine method to calculate great circle distances
func KmBetweenGC(a, b GPS2dT) float64 {
	x1 := Radians(a.Lon)
	y1 := Radians(a.Lat)

	x2 := Radians(b.Lon)
	y2 := Radians(b.Lat)

	dY := y1 - y2
	dX := x1 - x2

	c := math.Pow(math.Sin(dY/2), 2) + math.Cos(y1)*math.Cos(y2)*
		math.Pow(math.Sin(dX/2), 2)

	d := 2 * math.Atan2(math.Sqrt(c), math.Sqrt(1-c))

	return d * GaiaKm
}

func (pt GPS2dT) String() string {
	return fmt.Sprintf("{ Lon:%13.6f, Lat:%13.6f }", pt.Lon, pt.Lat)
}

//=============================================================================  GPS2dList
type GPS2dList []GPS2dT

func (g2 GPS2dList) BBox() BBox {
	var rv BBox
	var newPts []Pointe

	for _, val := range g2 {
		newPts = append(newPts, val.Point())
	}
	(&rv).ExpandByPts(newPts)
	return rv
}

func (g2 GPS2dList) PolyLine() PolyLine {
	//var rv PolyLine
	var part []Pointe
	for i := 0; i < len(g2); i++ {
		part = append(part, Pointe{X: g2[i].Lon, Y: g2[i].Lat})
	}
	return PolyLine{part}
}

func (x GPS2dT) KmTo(y GPS2dT) float64 {

	xx := Radians(x.Lon)
	xy := Radians(x.Lat)
	a := math.Sin(xx) * math.Cos(xy)
	b := math.Cos(xx) * math.Cos(xy)
	c := math.Sin(xy)
	//fmt.Printf("Position X(x,y,z) = %g %g %g\n", a,b,c)

	yx := Radians(y.Lon)
	yy := Radians(y.Lat)
	d := math.Sin(yx) * math.Cos(yy)
	e := math.Cos(yx) * math.Cos(yy)
	f := math.Sin(yy)
	//fmt.Printf("Position Y(x,y,z) = %g %g %g\n", d,e,f)

	km := a*d + b*e + c*f

	if km < -1.0 {
		km = -1
	}
	if km > 1.0 {
		km = 1.0
	}
	//log.Printf("d = %g\n", d)
	km = (math.Acos(km) / math.Pi) * 20000.0
	return km
}

func KmBetween(a, b GPS2dT) float64 {
	return a.KmTo(b)
}

// ============================================================================  GPS3d
type GPS3dT struct {
	Lon float64 // degrees
	Lat float64 // degrees
	Up  float64 // meters
}
