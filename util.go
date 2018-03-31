package fdfs

import (
	"encoding/binary"
	"strings"
)

func lengthByte(length uint64) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, length)
	return bs
}

func clearZero(s string) string {
	return strings.Replace(s, "\u0000", "", -1)
}
