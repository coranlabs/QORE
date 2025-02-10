package main

import (
    "fmt"
    "unsafe"
)

// #include <stdlib.h>
import "C"



//export FreeMemory
func FreeMemory(pointer *int64) {
    C.free(unsafe.Pointer(pointer))
}
