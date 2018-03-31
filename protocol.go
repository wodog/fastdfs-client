package fdfs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE               = 102
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE = 101
	TRACKER_QUERY_STORAGE_STORE_BODY_LEN                    = 40
	FDFS_GROUP_NAME_MAX_LEN                                 = 16
	STORAGE_PROTO_CMD_DOWNLOAD_FILE                         = 14
	STORAGE_PROTO_CMD_UPLOAD_FILE                           = 11
)

type protocol struct {
	header
	io.ReadWriter
}

type header struct {
	length  uint64
	command byte
	status  byte
}

func newProtocol(h header, rw io.ReadWriter) *protocol {
	return &protocol{
		h,
		rw,
	}
}

func (h *header) encode() []byte {
	buffer := &bytes.Buffer{}
	buffer.Write(lengthByte(h.length))
	buffer.WriteByte(byte(h.command))
	buffer.WriteByte(byte(h.status))
	return buffer.Bytes()
}

func (h *header) decode(b []byte) {
	h.length = binary.BigEndian.Uint64(b[0:8])
	h.command = b[8:9][0]
	h.status = b[9:10][0]
}

func (p *protocol) request(reqBody []byte) (resBody []byte, err error) {
	p.Write(p.encode())
	p.Write(reqBody)

	b := make([]byte, 10)
	_, err = io.ReadFull(p, b)
	if err != nil {
		return nil, err
	}
	p.decode(b)
	if p.status != 0 {
		return nil, errors.New("[tracker]状态码错误")
	}

	b = make([]byte, p.length)
	_, err = io.ReadFull(p, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
