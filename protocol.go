package fdfs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

const (
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE               = 102
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE = 101
	TRACKER_QUERY_STORAGE_STORE_BODY_LEN                    = 40
	STORAGE_PROTO_CMD_QUERY_FILE_INFO                       = 22
	FDFS_GROUP_NAME_MAX_LEN                                 = 16
	STORAGE_PROTO_CMD_DOWNLOAD_FILE                         = 14
	STORAGE_PROTO_CMD_DELETE_FILE                           = 12
	STORAGE_PROTO_CMD_UPLOAD_FILE                           = 11
)

type protocol struct {
	header
	io.ReadWriter
}

type header struct {
	length  uint64
	command uint8
	status  uint8
}

func newProtocol(h header, rw net.Conn) *protocol {
	return &protocol{
		h,
		rw,
	}
}

func (h *header) encode() []byte {
	buffer := &bytes.Buffer{}
	buffer.Write(lengthByte(h.length))
	buffer.WriteByte(h.command)
	buffer.WriteByte(h.status)
	return buffer.Bytes()
}

func (h *header) decode(b []byte) {
	h.length = binary.BigEndian.Uint64(b[:8])
	h.command = b[8]
	h.status = b[9]
}

func (p *protocol) body() ([]byte, error) {
	b := make([]byte, p.length)
	_, err := io.ReadFull(p, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (p *protocol) request(reqBody []byte) error {
	p.Write(p.encode())
	p.Write(reqBody)

	b := make([]byte, 10)
	_, err := io.ReadFull(p, b)
	if err != nil {
		return err
	}
	p.decode(b)
	if p.status != 0 {
		return errors.New("[tracker]状态码错误")
	}
	return nil
}
