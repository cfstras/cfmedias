package glib

/*
#cgo pkg-config: libgpod-1.0
#include "gpod/itdb.h"
#include "stdlib.h"
*/
import "C"
import (
	"github.com/cfstras/cfmedias/errrs"
	"unsafe"
)

func CStr(str string) *C.gchar {
	return (*C.gchar)(C.CString(str))
}

func Free(str *C.gchar) {
	C.free(unsafe.Pointer(str))
}

func Str(str *C.gchar) string {
	return C.GoString((*C.char)(str))
}

func Err(err *C.GError) error {
	str := Str(err.message)
	return errrs.New(str)
}
