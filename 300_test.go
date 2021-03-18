// 300_test.go

package mdr

import (
	"fmt"
	"testing"
	"time"
)

func Test_300(t *testing.T) {
	testName := "Test_300"
	runStart := time.Now()
	fmt.Printf("%s template\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()
	if false {
		t.Errorf("print fail, but keep testing")
	}
	if false {
		t.Fatalf("print fail and keep testing")
	}
}

func Test_301(t *testing.T) {
	testName := "Test_301"
	runStart := time.Now()
	fmt.Printf("%s float16 (Garmin)\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()
	x := Float16To64(681)
	Verbose.Printf("float16 for 681  = 2.660156\n")
	Verbose.Printf("float16 for 681  =%10.6f\n", x)
	if CloseToF64(x, 2.660156, 0.0001) == false {
		t.Errorf("Float16To64(681) not 2.660156")
	}
}

func Test_302(t *testing.T) {
	testName := "Test_302"
	runStart := time.Now()
	fmt.Printf("%s Float 32/64 Extraction\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()

	testTrue := []struct { // these should result in true
		test     string
		expected float64
	}{
		{test: `"-1234.5670"`, expected: -1234.5670},   // no comma
		{test: `"-1234.5671",`, expected: -1234.5671},  // post comma
		{test: `,"-1234.5672"`, expected: -1234.5672},  // pre comma
		{test: `,"-1234.5673",`, expected: -1234.5673}, // pre & post comma
		{test: `"-123.567"`, expected: -123.567},
	}

	for _, val := range testTrue {
		x := float64(ExtractFloat32(val.test))
		if CloseToF64(x, val.expected, 0.0001) == false {
			t.Errorf("ExtractFloat32 failed on %s", val.test)
		}
	}

	for _, val := range testTrue {
		x := ExtractFloat64(val.test)
		if CloseToF64(x, val.expected, 0.0001) == false {
			t.Errorf("ExtractFloat64 failed on %s", val.test)
		}
	}
}

func Test_303(t *testing.T) {
	testName := "Test_303"
	runStart := time.Now()
	fmt.Printf("%s Uint 8/16/32/64 Extraction\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()

	testTrue := []struct { // these should result in true
		test     string
		expected int64
	}{
		{test: `"1234"`, expected: 1234},   // no comma
		{test: `"1235",`, expected: 1235},  // post comma
		{test: `,"1236"`, expected: 1236},  // pre comma
		{test: `,"1237",`, expected: 1237}, // pre & post comma
		{test: `"123"`, expected: 123},
	}

	for _, val := range testTrue {
		x := ExtractUint16(val.test)
		if x != uint16(val.expected) {
			t.Errorf("ExtractUint16 failed on %s", val.test)
		}
	}

	for _, val := range testTrue {
		x := ExtractUint32(val.test)
		if x != uint32(val.expected) {
			t.Errorf("ExtractUint32 failed on %s", val.test)
		}
	}

	for _, val := range testTrue {
		x := ExtractUint64(val.test)
		if x != uint64(val.expected) {
			t.Errorf("ExtractUint64 failed on %s", val.test)
		}
	}
}

func Test_304(t *testing.T) {
	testName := "Test_304"
	runStart := time.Now()
	fmt.Printf("%s Int 8/16/32/64 Extraction\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()

	testTrue := []struct { // these should result in true
		test     string
		expected int64
	}{
		{test: `"-1234"`, expected: -1234},   // no comma
		{test: `"-1235",`, expected: -1235},  // post comma
		{test: `,"-1236"`, expected: -1236},  // pre comma
		{test: `,"-1237",`, expected: -1237}, // pre & post comma
		{test: `"-123"`, expected: -123},
	}

	for _, val := range testTrue {
		x := ExtractInt16(val.test)
		if x != int16(val.expected) {
			t.Errorf("ExtractInt16 failed on %s", val.test)
		}
	}

	for _, val := range testTrue {
		x := ExtractInt32(val.test)
		if x != int32(val.expected) {
			t.Errorf("ExtractInt32 failed on %s", val.test)
		}
	}

	for _, val := range testTrue {
		x := ExtractInt64(val.test)
		if x != int64(val.expected) {
			t.Errorf("ExtractInt64 failed on %s", val.test)
		}
	}
}

func Test_305(t *testing.T) {
	testName := "Test_305 Centigrade"
	runStart := time.Now()
	fmt.Printf("%s temperature F->C, C->F\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()

	testTrue := []struct { // these should result in true
		test     int16
		expected int16
	}{
		{test: 32, expected: 0},    // convert F -> C
		{test: 212, expected: 100}, // convert F -> C
	}

	for _, val := range testTrue {
		x := Centigrade(val.test)
		if x != val.expected {
			t.Errorf("Centigrade failed on %d (got %d)", val.test, x)
		}
	}

	testTrue = []struct { // these should result in true
		test     int16
		expected int16
	}{
		{test: -40, expected: -40}, // convert C -> F
		{test: 0, expected: 32},    // convert C -> F
		{test: 100, expected: 212}, // convert C -> F
	}

	for _, val := range testTrue {
		x := Fahrenheit(val.test)
		if x != val.expected {
			t.Errorf("Fahrenheit failed on %d (got %d) expected %d", val.test, x, val.expected)
		}
	}

}

func Test_306(t *testing.T) {
	testName := "Test_306 PolyLine NumParts()"
	runStart := time.Now()
	fmt.Printf("%s\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()

	var (
		a PolyLine = PolyLine{
			[]Pointe{{X: 2.0, Y: 2.0}, {X: 1.0, Y: 1.0}, {X: 3.0, Y: 3.0}},
		}
		b PolyLine = PolyLine{
			[]Pointe{{X: 1.0, Y: 1.0}, {X: 2.0, Y: 2.0}},
			[]Pointe{{X: 1.0, Y: 1.0}, {X: 2.0, Y: 2.0}},
			[]Pointe{{X: 1.0, Y: 1.0}, {X: 2.0, Y: 2.0}},
		}
		c PolyLine = PolyLine{
			[]Pointe{{X: 1.0, Y: 1.0}},
		}
		d   PolyLine
		pts []Pointe = []Pointe{{X: 1.0, Y: 1.0}, {X: 2.0, Y: 2.0}}
	)
	//Verbose = true
	Verbose.Printf("a has %d parts\n", a.NumParts())
	Verbose.Printf("b has %d parts\n", b.NumParts())
	Verbose.Printf("c has %d parts\n", c.NumParts())
	d.ExtendBy(pts)
	e := d.AddNewPart(pts)
	Verbose.Printf("d has %d parts\n", d.NumParts())
	Verbose.Printf("\n")
	Verbose.Printf("a has %d points\n", a.NumPoints())
	Verbose.Printf("b has %d points\n", b.NumPoints())
	Verbose.Printf("c has %d points\n", c.NumPoints())
	Verbose.Printf("d has %d points\n", d.NumPoints())
	Verbose.Printf("e has %d points\n", e.NumPoints())

	testTrue := []struct { // these should result in true
		test     PolyLine
		expected int
	}{
		{test: a, expected: 1},
		{test: b, expected: 3},
		{test: c, expected: 1},
		{test: d, expected: 1},
	}

	for _, val := range testTrue {
		x := (val.test).NumParts()
		if x != val.expected {
			t.Errorf("NumParts() failed on %v (got %d) expected %d", val.test, x, val.expected)
		}
	}

	testTrue = []struct { // these should result in true
		test     PolyLine
		expected int
	}{
		{test: a, expected: 3},
		{test: b, expected: 6},
		{test: c, expected: 1},
		{test: d, expected: 2},
		{test: e, expected: 4},
	}

	for _, val := range testTrue {
		x := (val.test).NumPoints()
		if x != val.expected {
			t.Errorf("NumPoints() failed on %v (got %d) expected %d", val.test, x, val.expected)
		}
	}

}

func Test_307(t *testing.T) {
	testName := "Test_307 ---"
	runStart := time.Now()
	fmt.Printf("%s\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()
}

func Test_308(t *testing.T) {
	testName := "Test_308 ---"
	runStart := time.Now()
	fmt.Printf("%s\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()
}

// Test_309() check PolyLine Centroid() method
func Test_309(t *testing.T) {
	testName := "Test_309 ExtendPolyLine"
	runStart := time.Now()
	fmt.Printf("%s\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()
	var (
		gpstest1 GPS2dList = GPS2dList{
			{Lat: 10, Lon: 10},
			{Lat: 15, Lon: 15},
		}
		gpstest2 GPS2dList = GPS2dList{
			//			{Lat:10,Lon:10,Up:10},
			//			{Lat:15,Lon:15,Up:10},
			{Lat: 36.810202, Lon: -77.025878},
			{Lat: 36.803840, Lon: -76.862869},
			{Lat: 36.814619, Lon: -76.902701},
		}
		pl PolyLine
		p2 PolyLine
		p3 PolyLine
	)
	Verbose.Printf("Centroid of empty list = %v\n", pl.Centroid())
	pl = append(pl, []Pointe{Pointe{X: 1.0, Y: 1.0}})
	Verbose.Printf("Centroid of single point [1,1] = %v\n", pl.Centroid())
	pl[0] = append(pl[0], Pointe{X: 100.0, Y: 100.0})
	Verbose.Printf("Centroid of line point [1,1],[100,100] = %v\n", pl.Centroid())
	Verbose.Printf("Expect [50.5,50.5]\n\n")
	Verbose.Printf("%v\n", gpstest1)
	Verbose.Printf("Centroid of gpslist is %v\n", pl.Centroid())
	gpsPoly1 := gpstest1.PolyLine()
	gpsPoly2 := gpstest2.PolyLine()

	p2.ExtendByPolyLine(gpsPoly1)
	Verbose.Printf("len p2 now = %d\n", len(p2))
	Verbose.Printf("Centroid of p2      is %v\n", p2.Centroid())

	p3.ExtendByPolyLine(gpsPoly2)
}

func Test_310(t *testing.T) {
	testName := "Test_310 BBox.Contains()"
	runStart := time.Now()
	fmt.Printf("%s\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()
	var ChainBox BBox = BBox{
		MinX: -77.5,
		MinY: 36.0,
		MaxX: -76,
		MaxY: 37.5,
	}

	var pt Pointe = Pointe{X: -76.5, Y: 37.0}
	Verbose.Printf("point in ChainBox (true) %v\n", ChainBox.Contains(pt))
	pt.X = 100.0
	Verbose.Printf("point in ChainBox (false) %v\n", ChainBox.Contains(pt))
}

func Test_311(t *testing.T) {
	testName := "Test_311 BBox from PolyLine"
	runStart := time.Now()
	fmt.Printf("%s\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()

	var ChainBox BBox = BBox{
		MinX: -77.5,
		MinY: 36.0,
		MaxX: -76,
		MaxY: 37.5,
	}

	var ChainList GPS2dList = GPS2dList{
		GPS2dT{Lat: 37.50, Lon: -77.50}, // UL
		GPS2dT{Lat: 37.50, Lon: -76.0},  // UR
		GPS2dT{Lat: 36.0, Lon: -76.0},   // LR
		GPS2dT{Lat: 36.0, Lon: -77.50},  // LL
		GPS2dT{Lat: 37.50, Lon: -77.50}, // UL used again to close box
	}
	newbox := ChainList.PolyLine().BBox()
	ok := ChainBox == newbox
	Verbose.Printf("BBox from GPS2dList test (true) = %v\n", ok)
	if !ok {
		t.Errorf("Failed -> should be equal\n")
		fmt.Printf("ChainBox %s\n", ChainBox)
		fmt.Printf("  newbox %s\n", newbox)
	}
}

func Test_312(t *testing.T) {
	testName := "Test_312"
	runStart := time.Now()
	fmt.Printf("%s BBox intersection\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()

	var (
		A []Pointe = []Pointe{
			Pointe{X: 10, Y: 10},
			Pointe{X: 50, Y: 50},
			Pointe{X: 100, Y: 100},
		}
		aBox BBox
		X    []Pointe = []Pointe{
			Pointe{X: 1, Y: 1},
			Pointe{X: 50, Y: 50},
		}
		xBox BBox
		Y    []Pointe = []Pointe{
			Pointe{X: 100, Y: 100},
			Pointe{X: 500, Y: 500},
		}
		yBox BBox
		Z    []Pointe = []Pointe{
			Pointe{X: 1, Y: 1},
			Pointe{X: 5, Y: 5},
		}
		zBox BBox
	)
	aBox.ExpandByPts(A)
	Verbose.Printf("aBox = %s\n", aBox)
	xBox.ExpandByPts(X)
	Verbose.Printf("xBox = %s\n", xBox)
	yBox.ExpandByPts(Y)
	Verbose.Printf("yBox = %s\n", yBox)
	zBox.ExpandByPts(Z)
	Verbose.Printf("zBox = %s\n", zBox)

	value := aBox.Intersects(xBox)
	expected := true
	if value != expected {
		t.Errorf("Intersects() failed")
	}

	value = aBox.Intersects(zBox)
	expected = false
	if value != expected {
		t.Errorf("Intersects() failed")
	}

	value = aBox.Intersects(yBox)
	expected = true
	if value != expected {
		t.Errorf("Intersects() failed")
	}

	value = xBox.Intersects(zBox)
	expected = true
	if value != expected {
		t.Errorf("Intersects() failed")
	}
}

func Test_313(t *testing.T) {
	testName := "Test_313"
	runStart := time.Now()
	fmt.Printf("%s Expand BBox\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()
	var (
		A PolyLine = PolyLine{
			[]Pointe{
				Pointe{X: 10, Y: 10},
				Pointe{X: 50, Y: 50},
				Pointe{X: 100, Y: 100},
			},
			[]Pointe{
				Pointe{X: 1, Y: 1},
				Pointe{X: 50, Y: 50},
			},
		}
		bbox     BBox
		expected BBox = BBox{MinX: 1, MinY: 1, MaxX: 100, MaxY: 100}
	)
	Verbose.Printf("bbox = %s\n", bbox)
	Verbose.Printf("PolyLine = %s\n", A)
	bbox.ExpandByPolyLine(A)
	Verbose.Printf("bbox = %s\n", bbox)
	if bbox != expected {
		t.Errorf("not expected value\n")
	}
}

// Simple example of commutative failure, value2 failed though value1 is ok...
//  definitely should NOT matter what order the args appear
//BUG(mdr) corner generation seemed ok...
// failed when one box is contained by another fixed now I believe
func Test_314(t *testing.T) {
	testName := "Test_314"
	runStart := time.Now()
	fmt.Printf("%s BBox Intersections\n", testName)
	defer func() {
		Verbose.Printf("%s took %v\n", testName, time.Since(runStart))
	}()
	var (
		mybox   BBox = BBox{MinX: -76.932480, MinY: 36.549660, MaxX: -76.407532, MaxY: 36.921856}
		railbox BBox = BBox{MinX: -76.917970, MinY: 36.682432, MaxX: -76.916911, MaxY: 36.684338}
	)

	if !mybox.Sane() || !railbox.Sane() {
		t.Fatalf("Intersects() failed")
	}
	expected := true
	value1 := mybox.Intersects(railbox)
	if value1 != expected {
		t.Errorf("Intersects() failed")
	}
	value2 := railbox.Intersects(mybox)
	if value2 != expected {
		t.Errorf("Intersects() failed")
	}
}
