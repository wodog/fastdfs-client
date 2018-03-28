package fdfs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
)

type Storage struct {
	group string
	host  string
	port  int
	index byte
}

func (s *Storage) upload(file io.Reader) string {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		panic(err)
	}
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	respHeader := &Header{
		length:  uint64(len(bs)) + 15,
		command: STORAGE_PROTO_CMD_UPLOAD_FILE,
		status:  0,
	}
	conn.Write(respHeader.encode())

	buffer := &bytes.Buffer{}
	buffer.WriteByte(s.index)
	buffer.Write(lengthByte(uint64(len(bs))))
	buffer.Write(lengthByte(0))
	conn.Write(buffer.Bytes())
	conn.Write(bs)

	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		panic(err)
	}
	storageResp := b[:n]

	// 检查header
	storageRespHeader := &Header{
		buf: storageResp[:10],
	}
	storageRespHeader.decode()
	if storageRespHeader.status != 0 {
		panic(err)
	}

	storageRespBody := storageResp[10:]
	group := string(storageRespBody[:16])
	path := string(storageRespBody[16:])
	return group + "/" + path
}
