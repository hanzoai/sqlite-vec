// Package sqlitevec compiles the sqlite-vec extension into the binary and
// registers it for every SQLite connection the process opens.
//
// This is the cgo companion to hanzoai/csqlite: wrapper.c textually includes
// ../sqlite-vec.c (the repo root keeps the ONE copy of the C — sqlite-vec.c
// itself #include's its -rescore/-diskann siblings, so compiling the root
// directory as a Go package would compile those fragments a second time), and the sqlite3_* symbols resolve at final link against the
// SQLite that csqlite already carries — one SQLite in the binary, one extension
// object, no shared library and no vendored second copy of either.
//
// SQLITE_CORE is deliberate: it makes sqlite-vec.h include sqlite3.h (committed
// here, copied byte-for-byte from csqlite's sqlite3-binding.h so the header can
// never drift from the library that links) instead of the loadable-extension
// shim. Upstream's own Go bindings compile the same way but leave sqlite3.h to
// whatever the SYSTEM has installed — a build that works only when
// libsqlite3-dev happens to be present, and against whatever version it
// happens to be. Committing the matching header removes that ambient
// dependency.
//
// Ported from asg017/sqlite-vec-go-bindings (MIT/Apache-2.0, same author and
// licence as sqlite-vec itself).
package sqlitevec

/*
#cgo CFLAGS: -I${SRCDIR}/.. -DSQLITE_CORE
#cgo LDFLAGS: -lm
#include "sqlite-vec.h"
*/
import "C"

import (
	"bytes"
	"encoding/binary"
)

// Auto registers sqlite-vec on every SQLite connection this process opens from
// now on, via sqlite3_auto_extension. Call it once, before sql.Open.
func Auto() {
	C.sqlite3_auto_extension((*[0]byte)(C.sqlite3_vec_init))
}

// Cancel undoes Auto for connections opened after the call.
func Cancel() {
	C.sqlite3_cancel_auto_extension((*[0]byte)(C.sqlite3_vec_init))
}

// SerializeFloat32 encodes a float32 vector as the little-endian BLOB
// sqlite-vec expects for float[N] columns and MATCH parameters.
func SerializeFloat32(vector []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, vector); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
