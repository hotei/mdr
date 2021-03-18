// mdr_point.go

package mdr

import (
	"fmt"
	"log"
)

// ============================================================================  Pointi
type Pointi struct {
	X int
	Y int
}

// ============================================================================  PointiList
type PointiList []Pointi

// Split separates the X and Y values into their own arrays
func (pts PointiList) Split() (xs []int, ys []int) {
	for i := 0; i < len(pts); i++ {
		xs = append(xs, pts[i].X)
		ys = append(ys, pts[i].Y)
	}
	return xs, ys
}

// RangeMinMaxPoint returns bounds of point array.
// probably a better name would help... 8888
func RangeMinMaxPoint(v PointiList) (minPt, maxPt Pointi) {
	vlen := len(v)
	if vlen <= 0 {
		// warn?
		return
	}
	minPt, maxPt = v[0], v[0]
	for i := 1; i < vlen; i++ {
		if v[i].X < minPt.X {
			minPt.X = v[i].X
		}
		if v[i].Y < minPt.Y {
			minPt.Y = v[i].Y
		}
		if v[i].X > maxPt.X {
			maxPt.X = v[i].X
		}
		if v[i].Y > maxPt.Y {
			maxPt.Y = v[i].Y
		}
	}
	return minPt, maxPt
}

// ============================================================================  Pointe
type Pointe struct {
	X float64
	Y float64
}

func (pt Pointe) String() string {
	return fmt.Sprintf("{ X:%13.6f, Y:%13.6f }", pt.X, pt.Y)
}

func (pt Pointe) ValidGPS2d() bool {
	return InRangeF64(-180.0, pt.X, 180.0) && InRangeF64(-90.0, pt.Y, 90.0)
}

func (pt Pointe) GPS2d() GPS2dT {
	var rv GPS2dT
	rv.Lon = pt.X
	rv.Lat = pt.Y
	return rv
}

func (pt Pointe) KmTo(pt2 Pointe) float64 {
	return KmBetweenGC(pt.GPS2d(), pt2.GPS2d())
}

// ============================================================================  PointeList
type PointeList []Pointe

// ============================================================================  BBox
type BBox struct {
	MinX, MinY, MaxX, MaxY float64
}

func (bb BBox) String() string { // LwrLeft, UprRight
	return fmt.Sprintf("{ MinX:%12.6f, MinY:%12.6f, MaxX:%12.6f, MaxY:%12.6f }",
		bb.MinX, bb.MinY, bb.MaxX, bb.MaxY)
}

func (bb *BBox) ExpandByPts(pts []Pointe) {
	if len(pts) <= 0 {
		return
	}
	if bb.MinX+bb.MaxX+bb.MinY+bb.MaxY == 0 {
		bb.MinX = pts[0].X
		bb.MaxX = pts[0].X
		bb.MinY = pts[0].Y
		bb.MaxY = pts[0].Y
	}
	for i := 0; i < len(pts); i++ {
		pt := pts[i]
		if pt.X < bb.MinX {
			bb.MinX = pt.X
		}
		if pt.X > bb.MaxX {
			bb.MaxX = pt.X
		}
		if pt.Y < bb.MinY {
			bb.MinY = pt.Y
		}
		if pt.Y > bb.MaxY {
			bb.MaxY = pt.Y
		}
	}
}

func (bb *BBox) ExpandByPolyLine(parts PolyLine) {
	if parts.NumPoints() <= 0 {
		return
	}
	if bb.MinX+bb.MaxX+bb.MinY+bb.MaxY == 0 {
		bb.MinX = parts[0][0].X
		bb.MaxX = parts[0][0].X
		bb.MinY = parts[0][0].Y
		bb.MaxY = parts[0][0].Y
	}
	for j := 0; j < len(parts); j++ {
		bb.ExpandByPts(parts[j])
	}
}

func (bb BBox) Centroid() Pointe {
	var (
		c          = bb.Corners()
		xsum, ysum float64
		n          float64
	)
	for j := 0; j < len(c); j++ {
		xsum += c[j].X
		ysum += c[j].Y
		n++
	}
	return Pointe{xsum / n, ysum / n}
}

// Corners() returns all 4 corners though only two are needed to define a box.
// Easier to draw if you have all 4, and they are in drawable order.
func (bb BBox) Corners() []Pointe {
	// must retain order LL,UL,UR,LR for draw to work
	var rv = []Pointe{
		Pointe{X: bb.MinX, Y: bb.MinY}, // LL
		Pointe{X: bb.MinX, Y: bb.MaxY}, // UL
		Pointe{X: bb.MaxX, Y: bb.MaxY}, // UR
		Pointe{X: bb.MaxX, Y: bb.MinY}, // LR
	}
	return rv
}

func (bb BBox) Contains(pt Pointe) bool {
	if Paranoid && !bb.Sane() {
		log.Fatalf("bad bounding box found")
	}
	x := pt.X
	y := pt.Y
	if x < bb.MinX {
		//fmt.Printf("False  ContainsBBox( %s )  Pt ( %s )\n", bb, pt)
		return false
	}
	if y < bb.MinY {
		//fmt.Printf("False  ContainsBBox( %s )  Pt ( %s )\n", bb, pt)
		return false
	}
	if x > bb.MaxX {
		//fmt.Printf("False  ContainsBBox( %s )  Pt ( %s )\n", bb, pt)
		return false
	}
	if y > bb.MaxY {
		//fmt.Printf("False  ContainsBBox( %s )  Pt ( %s )\n", bb, pt)
		return false
	}
	return true
}

// Intersects() returns true if box areas overlap at any point.
// Failed when one box completely encloses the other so add second test
// Should return true if one encloses the other
func (bb BBox) Intersects(xx BBox) bool {
	if Paranoid && (!bb.Sane() || !xx.Sane()) {
		log.Fatalf("bad bounding box found")
	}
	for _, val := range xx.Corners() {
		if bb.Contains(val) {
			//Verbose.Printf("pt %s is inside %s \n", val, bb)
			return true
		}
	}
	for _, val := range bb.Corners() {
		if xx.Contains(val) {
			//Verbose.Printf("pt %s is inside %s \n", val, xx)
			return true
		}
	}
	return false
}

// BUG(mdr): !Sane or just not init???
// BBox.Sane() makes basic checks, usually employed only if Paranoid turned on
//   see BBox.Contains() above as example
func (bb BBox) Sane() bool {
	// all zero if not initialized, usually not a good thing
	if bb.MinX+bb.MaxX+bb.MinY+bb.MaxY == 0 {
		return false
	}
	if bb.MinX > bb.MaxX {
		return false
	}
	if bb.MaxX < bb.MinX {
		return false
	}
	if bb.MinY > bb.MaxY {
		return false
	}
	if bb.MaxY < bb.MinY {
		return false
	}
	return true
}
