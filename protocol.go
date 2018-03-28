package fdfs

import (
	"bytes"
	"encoding/binary"
)

const (
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE               = 102
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE = 101
	TRACKER_QUERY_STORAGE_STORE_BODY_LEN                    = 40
	FDFS_GROUP_NAME_MAX_LEN                                 = 16
	STORAGE_PROTO_CMD_DOWNLOAD_FILE                         = 14
	STORAGE_PROTO_CMD_UPLOAD_FILE                           = 11
)

type Header struct {
	length  uint64
	command byte
	status  byte
	buf     []byte
}

func (h *Header) encode() []byte {
	buffer := &bytes.Buffer{}
	buffer.Write(lengthByte(h.length))
	buffer.WriteByte(byte(h.command))
	buffer.WriteByte(byte(h.status))
	return buffer.Bytes()
}

func (h *Header) decode() {
	h.length = binary.BigEndian.Uint64(h.buf[0:8])
	h.command = h.buf[8:9][0]
	h.status = h.buf[9:10][0]
}

func lengthByte(length uint64) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, length)
	return bs
}
