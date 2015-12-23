// Copyright 2015 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
Package proj is a Go wrapper around the Proj4 C library.

It provides coordinate reference system definition and transformation functionalities.

No shared library is required at runtime as the C code is ntegrated in the package.
The PROJ_LIB environment variable must be set or the
SetFinder function must be called in order to identify the location of the share
folder.
*/
package proj

// #cgo windows CFLAGS: -DHAVE_LOCALECONV
// #include "proj_api.h"
import "C"

import (
	"errors"
	"runtime"
	"unsafe"

	"github.com/xeonx/geom"
)

//Proj represents a coordinate reference system.
//It is not safe for concurent use
type Proj struct {
	proj C.projPJ
	ctx  C.projCtx
}

//InitPlus initializes a new projection from a proj4 plus string (eg. "+init=epsg:4326" )
func InitPlus(definition string) (*Proj, error) {

	ctx := C.pj_ctx_alloc()
	if ctx == nil {
		errnoRef := C.pj_get_errno_ref()
		if errnoRef == nil {
			return nil, errors.New("unknown error on pj_ctx_alloc")
		}
		return nil, errors.New(C.GoString(C.pj_strerrno(*errnoRef)))
	}

	c := C.CString(definition)
	defer C.free(unsafe.Pointer(c))
	proj := C.pj_init_plus_ctx(ctx, c)
	if proj == nil {
		errno := C.pj_ctx_get_errno(ctx)
		return nil, errors.New(C.GoString(C.pj_strerrno(errno)))
	}

	p := &Proj{
		proj: proj,
		ctx:  ctx,
	}
	runtime.SetFinalizer(p, closeProj)
	return p, nil
}

func closeProj(p *Proj) {
	p.Close()
}

//Close deallocates the projection immediately. Otherwise, it will be deallocated on garbage collection.
func (p *Proj) Close() {
	if p.proj != nil {
		C.pj_free(p.proj)
		p.proj = nil
	}
	if p.ctx != nil {
		C.pj_ctx_free(p.ctx)
		p.ctx = nil
	}
}

//IsLatLong returns whether the projection is geographic
func (p *Proj) IsLatLong() bool {
	return C.pj_is_latlong(p.proj) != 0
}

//IsGeoCent returns whether the projection is geocentric
func (p *Proj) IsGeoCent() bool {
	return C.pj_is_geocent(p.proj) != 0
}

//GetDef returns an initialization string suitable for use with InitPlus
func (p *Proj) GetDef() string {
	return C.GoString(C.pj_get_def(p.proj, 0))
}

//TransformRaw transforms the x/y/z points from the source coordinate system to the destination coordinate system.
//zs can be nil or must have the same length as xs and ys.
func TransformRaw(src *Proj, dst *Proj, xs []float64, ys []float64, zs []float64) error {

	if len(xs) != len(ys) {
		return errors.New("Incoherent slice size between x and y.")
	}
	if zs != nil && len(xs) != len(zs) {
		return errors.New("Incoherent slice size between x and z.")
	}

	if len(xs) == 0 {
		return nil
	}

	var errno C.int

	if zs != nil {

		errno = C.pj_transform(src.proj, dst.proj,
			C.long(len(xs)),
			1,
			(*C.double)(unsafe.Pointer(&xs[0])),
			(*C.double)(unsafe.Pointer(&ys[0])),
			(*C.double)(unsafe.Pointer(&zs[0])))
	} else {

		errno = C.pj_transform(src.proj, dst.proj,
			C.long(len(xs)),
			1,
			(*C.double)(unsafe.Pointer(&xs[0])),
			(*C.double)(unsafe.Pointer(&ys[0])),
			nil)
	}

	if errno != 0 {
		return errors.New(C.GoString(C.pj_strerrno(errno)))
	}

	return nil

}

var geomPointSize = (C.int)(unsafe.Sizeof(geom.Point{}) / unsafe.Sizeof(float64(0.0)))

//TransformPoints transforms the points inplace from the source coordinate system to the destination coordinate system.
func TransformPoints(src *Proj, dst *Proj, points []geom.Point) error {

	errno := C.pj_transform(src.proj, dst.proj,
		C.long(len(points)),
		geomPointSize,
		(*C.double)(unsafe.Pointer(&points[0].X)),
		(*C.double)(unsafe.Pointer(&points[0].Y)),
		nil)

	if errno != 0 {
		return errors.New(C.GoString(C.pj_strerrno(errno)))
	}

	return nil

}

//Transformation projects coordinates from a source to a destination
type Transformation struct {
	src *Proj
	dst *Proj
}

//NewTransformation initializes a new transformation with src and dst
func NewTransformation(src, dst *Proj) (Transformation, error) {
	return Transformation{src, dst}, nil
}

//TransformRaw transforms the x/y/z points from the source coordinate system to the destination coordinate system.
//zs can be nil or must have the same length as xs and ys.
func (t Transformation) TransformRaw(xs, ys, zs []float64) error {
	return TransformRaw(t.src, t.dst, xs, ys, zs)
}

//TransformPoints transforms the points inplace from the source coordinate system to the destination coordinate system.
func (t Transformation) TransformPoints(points []geom.Point) error {
	return TransformPoints(t.src, t.dst, points)
}
