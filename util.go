package fdfs

import (
	"encoding/binary"
	"strings"
)

func lengthByte(length uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, length)
	return b
}

func clearZero(s string) string {
	return strings.Replace(s, "\u0000", "", -1)
}
