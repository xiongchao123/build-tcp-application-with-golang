package util

import (
	"unsafe"
	"encoding/binary"
)

const N int = int(unsafe.Sizeof(0))

func GetOrder() binary.ByteOrder{
	x := 0x1234
	p := unsafe.Pointer(&x)
	p2 := (*[N]byte)(p)
	if p2[0] == 0 {
		return binary.BigEndian
	} else {
		return binary.LittleEndian
	}
}