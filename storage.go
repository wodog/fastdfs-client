package fdfs

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
)

type Storage struct {
	group string
	host  string
	port  int
	index byte
}

func (s *Storage) upload(file io.Reader) (string, error) {
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return "", err
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
	buffer.Write(make([]byte, 6))
	conn.Write(buffer.Bytes())
	conn.Write(bs)

	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		return "", err
	}
	storageResp := b[:n]

	// 检查header
	storageRespHeader := &Header{
		buf: storageResp[:10],
	}
	storageRespHeader.decode()
	if storageRespHeader.status != 0 {
		return "", errors.New("[storage]状态码错误")
	}

	storageRespBody := storageResp[10:]
	group := string(storageRespBody[:16])
	path := string(storageRespBody[16:])
	return group + "/" + path, nil
}

func (s *Storage) download(fileId string, w io.Writer) error {
	ss := strings.SplitN(fileId, "/", 2)
	group := ss[0]
	path := ss[1]

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return err
	}
	respHeader := &Header{
		length:  8 + 8 + 16 + uint64(len(path)),
		command: STORAGE_PROTO_CMD_DOWNLOAD_FILE,
		status:  0,
	}
	conn.Write(respHeader.encode())
	conn.Write(lengthByte(0))
	conn.Write(lengthByte(0))
	r1 := strings.NewReader(group)
	b1 := make([]byte, 16)
	r1.Read(b1)
	conn.Write(b1)
	conn.Write([]byte(path))

	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		return err
	}
	storageResp := b[:n]
	// 检查header
	storageRespHeader := &Header{
		buf: storageResp[:10],
	}
	storageRespHeader.decode()
	if storageRespHeader.status != 0 {
		return errors.New("[storage]状态码错误")
	}
	storageRespBody := storageResp[10:]
	w.Write(storageRespBody)
	return nil
}
