// Copyright 2015 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proj

import (
	"math"
	"testing"

	"github.com/xeonx/geom"
)

func TestNew(t *testing.T) {
	p, err := InitPlus("+init=epsg:4326")
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("projection is nil")
	}

	p, err = InitPlus("+init=epsg:999999")
	if err == nil {
		t.Fatal("no error for unknown projection")
	}

	p, err = InitPlus("")
	if err == nil {
		t.Fatal("no error for empty projection")
	}

	p, err = InitPlus(" +proj=utm +zone=32 +ellps=GRS80 +towgs84=0,0,0,0,0,0,0 +units=m +no_defs ")
	if err != nil {
		t.Fatal(err)
	}
	if p == nil {
		t.Fatal("projection is nil")
	}

	p, err = InitPlus(" +proj=utm +zone=99 +ellps=GRS80 +towgs84=0,0,0,0,0,0,0 +units=m +no_defs ")
	if err == nil {
		t.Fatal("no error for invalid projection")
	}
}

type transformTest struct {
	srcProj string
	dstProj string
	src     geom.Point
	dst     geom.Point
	err     error
}

const deg2Rad = math.Pi / 180.0
const rad2Deg = 180.0 / math.Pi

var transformData = []transformTest{
	{"+init=epsg:4326", "+init=epsg:25832", geom.Point{X: 8.15 * deg2Rad, Y: 53.2 * deg2Rad}, geom.Point{X: 443220.719, Y: 5894856.508}, nil},
}

func TestTransformRaw(t *testing.T) {

	for _, testCase := range transformData {

		//Initialise projections
		p1, err := InitPlus(testCase.srcProj)
		if err != nil {
			t.Fatal(err)
		}
		defer p1.Close()
		p2, err := InitPlus(testCase.dstProj)
		if err != nil {
			t.Fatal(err)
		}
		defer p2.Close()

		//Create raw array
		iCount := 32
		var xs, ys []float64
		for i := 0; i < iCount; i++ {
			xs = append(xs, testCase.src.X)
			ys = append(ys, testCase.src.Y)
		}

		if err := TransformRaw(p1, p2, xs, ys, nil); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < iCount; i++ {
			if math.Abs(xs[i]-testCase.dst.X) > 0.01 {
				t.Error(xs)
			}
			if math.Abs(ys[i]-testCase.dst.Y) > 0.01 {
				t.Error(ys)
			}
		}

		if err := TransformRaw(p2, p1, xs, ys, nil); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < iCount; i++ {
			if math.Abs(xs[i]-testCase.src.X) > 0.01 {
				t.Error(xs)
			}
			if math.Abs(ys[i]-testCase.src.Y) > 0.01 {
				t.Error(ys)
			}
		}
	}
}

func TestTransformRawZ(t *testing.T) {

	for _, testCase := range transformData {

		//Initialise projections
		p1, err := InitPlus(testCase.srcProj)
		if err != nil {
			t.Fatal(err)
		}
		defer p1.Close()
		p2, err := InitPlus(testCase.dstProj)
		if err != nil {
			t.Fatal(err)
		}
		defer p2.Close()

		//Create raw array
		iCount := 32
		var xs, ys, zs []float64
		for i := 0; i < iCount; i++ {
			xs = append(xs, testCase.src.X)
			ys = append(ys, testCase.src.Y)
			zs = append(zs, 12.)
		}

		if err := TransformRaw(p1, p2, xs, ys, zs); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < iCount; i++ {
			if math.Abs(xs[i]-testCase.dst.X) > 0.01 {
				t.Error(xs)
			}
			if math.Abs(ys[i]-testCase.dst.Y) > 0.01 {
				t.Error(ys)
			}
			if math.Abs(zs[i]-12.) > 0.01 {
				t.Error(zs)
			}
		}

		if err := TransformRaw(p2, p1, xs, ys, zs); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < iCount; i++ {
			if math.Abs(xs[i]-testCase.src.X) > 0.01 {
				t.Error(xs)
			}
			if math.Abs(ys[i]-testCase.src.Y) > 0.01 {
				t.Error(ys)
			}
			if math.Abs(zs[i]-12.) > 0.01 {
				t.Error(zs)
			}
		}
	}
}
func TestTransformPoint(t *testing.T) {

	for _, testCase := range transformData {

		//Initialise projections
		p1, err := InitPlus(testCase.srcProj)
		if err != nil {
			t.Fatal(err)
		}
		defer p1.Close()
		p2, err := InitPlus(testCase.dstProj)
		if err != nil {
			t.Fatal(err)
		}
		defer p2.Close()

		//Create raw array
		iCount := 32
		var points []geom.Point
		for i := 0; i < iCount; i++ {
			points = append(points, testCase.src)
		}

		if err := TransformPoints(p1, p2, points); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < iCount; i++ {
			if math.Abs(points[i].X-testCase.dst.X) > 0.01 {
				t.Error(points)
			}
			if math.Abs(points[i].Y-testCase.dst.Y) > 0.01 {
				t.Error(points)
			}
		}

		if err := TransformPoints(p2, p1, points); err != nil {
			t.Fatal(err)
		}

		for i := 0; i < iCount; i++ {
			if math.Abs(points[i].X-testCase.src.X) > 0.01 {
				t.Error(points)
			}
			if math.Abs(points[i].Y-testCase.src.Y) > 0.01 {
				t.Error(points)
			}
		}
	}
}
