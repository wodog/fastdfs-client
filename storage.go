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
	host  string
	port  string
	group string
	index byte
}

func (s *Storage) upload(file io.Reader) (string, error) {
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", s.host, s.port), timeout)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	buf := &bytes.Buffer{}
	buf.WriteByte(s.index)
	buf.Write(lengthByte(uint64(len(bs))))
	buf.Write(make([]byte, 6))
	buf.Write(bs)
	respHeader := &Header{
		length:  uint64(buf.Len()),
		command: STORAGE_PROTO_CMD_UPLOAD_FILE,
		status:  0,
	}
	conn.Write(respHeader.encode())
	conn.Write(buf.Bytes())

	b := make([]byte, 10)
	_, err = io.ReadFull(conn, b)
	if err != nil {
		return "", err
	}
	header := &Header{
		buf: b,
	}
	header.decode()
	if header.status != 0 {
		return "", errors.New("[storage]状态码错误")
	}
	b = make([]byte, header.length)
	_, err = io.ReadFull(conn, b)
	if err != nil {
		return "", err
	}
	group := clearZero(string(b[:16]))
	path := clearZero(string(b[16:]))
	return group + "/" + path, nil
}

func (s *Storage) download(fileId string, w io.Writer) error {
	ss := strings.SplitN(fileId, "/", 2)
	group := ss[0]
	path := ss[1]

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", s.host, s.port), timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	buf := &bytes.Buffer{}
	buf.Write(lengthByte(0))
	buf.Write(lengthByte(0))
	b := make([]byte, 16)
	copy(b, group)
	buf.Write(b)
	buf.WriteString(path)

	header := &Header{
		length:  uint64(buf.Len()),
		command: STORAGE_PROTO_CMD_DOWNLOAD_FILE,
		status:  0,
	}
	conn.Write(header.encode())
	conn.Write(buf.Bytes())

	b = make([]byte, 10)
	_, err = io.ReadFull(conn, b)
	if err != nil {
		return err
	}
	header = &Header{
		buf: b,
	}
	header.decode()
	if header.status != 0 {
		return errors.New("[storage]状态码错误")
	}
	b = make([]byte, header.length)
	_, err = io.ReadFull(conn, b)
	if err != nil {
		return err
	}
	w.Write(b)
	return nil
}
