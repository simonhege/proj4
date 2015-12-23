// Copyright 2015 Simon HEGE. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package proj

// #cgo windows CFLAGS: -DHAVE_LOCALECONV
// #include "proj_api.h"
// extern char *go_finder_wrapper(char *name);
import "C"

import (
	"os"
	"path/filepath"
	"unsafe"
)

var searchPaths []string
var finderResults map[string]*C.char

//export goFinder
func goFinder(cname *C.char) *C.char {
	name := C.GoString(cname)
	path, ok := finderResults[name]
	if !ok {
		for _, p := range searchPaths {
			p = filepath.Join(p, name)
			_, err := os.Stat(p)
			if err == nil {
				path = C.CString(p)
				break
			}
		}
		// cache result, even if it is nil
		finderResults[name] = path
	}
	return path
}

// SetFinder add one or more directories to search for proj definition files.
// Multiple calls overwrite the previous search paths.
func SetFinder(paths []string) {
	finderResults = make(map[string]*C.char)
	searchPaths = paths
	C.pj_set_finder((*[0]byte)(unsafe.Pointer(C.go_finder_wrapper)))
}
