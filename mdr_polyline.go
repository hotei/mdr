// polyline.go (c) 2019 David Rook - all rights reserved

/*
While the polyline can be used to depict GPS coordinates that is NOT a requirement.
Use []GPSloc aka GPSlist if you need that (3D) capability.
*/
package mdr

import (
	//"fmt"
	"log"
)

//========================================================================
type PolyLine [][]Pointe

// BBox() returns the smallest box which encloses all the points of the PolyLine
func (pl PolyLine) BBox() BBox {
	var rv BBox
	if Paranoid && len(pl) <= 0 {
		log.Fatalf("asked to create an empty bounds box\n")
	}
	rv.MinX = pl[0][0].X
	rv.MinY = pl[0][0].Y
	rv.MaxX = pl[0][0].X
	rv.MaxY = pl[0][0].Y
	for i := 0; i < len(pl); i++ {
		rv.ExpandByPts(pl[i])
	}
	return rv
}

func (pl *PolyLine) ExtendBy(part []Pointe) {
	*pl = append(*pl, part)
	return
}

func (pl *PolyLine) ExtendByPolyLine(pl2 PolyLine) {
	for i := 0; i < len(pl2); i++ {
		*pl = append(*pl, pl2[i])
	}
	return
}

func (pl PolyLine) AddNewPart(part []Pointe) PolyLine {
	pl = append(pl, part)
	return pl
}

// Centroid() returns the "center of gravity" of the shape
func (pl PolyLine) Centroid() Pointe {
	var (
		rv         Pointe
		xsum, ysum float64
		n          float64
	)
	if len(pl) <= 0 {
		return rv
		// BUG(mdr): or is this an error ?
	}
	if len(pl[0]) == 1 {
		return pl[0][0]
	}
	for j := 0; j < len(pl); j++ {
		for i := 0; i < len(pl[j]); i++ {
			xsum += pl[j][i].X
			ysum += pl[j][i].Y
			n++
		}
	}
	return Pointe{xsum / n, ysum / n}
}

func (pl PolyLine) NumParts() int {
	return len(pl)
}

// TODO(mdr): Variadic version ?
func Join(a, b PolyLine) PolyLine {
	var rv PolyLine
	if a != nil {
		(&rv).ExtendByPolyLine(a)
	}
	if b != nil {
		(&rv).ExtendByPolyLine(b)
	}
	return rv
}

// NumPts() returns the number of points in all parts.
func (pl PolyLine) NumPoints() int {
	var rv int

	for i := 0; i < len(pl); i++ {
		rv += len(pl[i])
	}
	return rv
}

func MakePoly(pts GPS2dList) PolyLine {
	var part []Pointe
	for i := 0; i < len(pts); i++ {
		pt := pts[i].Point()
		part = append(part, pt)
	}
	return PolyLine{part}
}

/*
// Extend returns a copy of the original after adding a part containing the new points.
func (pl PolyLine) Extend(gl []GPSloc) PolyLine {
	var (
		pts []Pointe
	)
	for i := 0; i < len(gl); i++ {
		gl2 := gl[i].GPS2d()
		pts = append(pts, Pointe{X: gl2.Lon, Y: gl2.Lat})
	}
	pl = append(pl, pts)
	return pl
}

// Extend a PolyLine in place.
func (pl *PolyLine) Extend2(gl []GPSloc) {
	var (
		pts []Point
	)
	for i := 0; i < len(gl); i++ {
		gl2 := gl[i].GPS2d()
		pts = append(pts, Point{X: gl2.Lon, Y: gl2.Lat})
	}
	*pl = append(*pl, pts)

}
*/
