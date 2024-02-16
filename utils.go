package taglib

// #cgo pkg-config: taglib
// #cgo LDFLAGS: -ltag_c
// #include <stdlib.h>
// #include <tag_c.h>
import "C"
import (
	"unsafe"
)

func convertAndFree(cstr *C.char) string {
	defer C.free(unsafe.Pointer(cstr))
	return C.GoString(cstr)
}

func toGoStringArray(cArray **C.char) []string {
	var goArray []string

	elem := cArray

	for elem != nil && *elem != nil {
		goArray = append(goArray, C.GoString(*elem))

		elem = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(elem)) + unsafe.Sizeof(*elem)))
	}

	return goArray
}
