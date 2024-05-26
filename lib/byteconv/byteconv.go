package byteconv

import "unsafe"

func Bytes(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func String(slice []byte) string {
	return unsafe.String(unsafe.SliceData(slice), len(slice))
}
