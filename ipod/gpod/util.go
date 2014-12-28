package gpod

/*
#include "stdlib.h"
#include "glib.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

func cstr(str string) *C.gchar {
	return (*C.gchar)(C.CString(str))
}

func free(str *C.gchar) {
	C.free(unsafe.Pointer(str))
}

func str(str *C.gchar) string {
	return C.GoString((*C.char)(str))
}

func err(err *C.GError) error {
	str := str(err.message)
	return errors.New(str)
}
