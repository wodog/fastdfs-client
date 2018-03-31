package fdfs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strings"
)

type storage struct {
	host  string
	port  string
	group string
	index byte
}

func (s *storage) upload(file io.Reader) (string, error) {
	b, err := ioutil.ReadAll(file)
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
	buf.Write(lengthByte(uint64(len(b))))
	buf.Write(make([]byte, 6))
	buf.Write(b)
	h := header{
		uint64(buf.Len()),
		STORAGE_PROTO_CMD_UPLOAD_FILE,
		0,
	}
	p := newProtocol(h, conn)
	b, err = p.request(buf.Bytes())
	if err != nil {
		return "", err
	}

	group := clearZero(string(b[:16]))
	path := clearZero(string(b[16:]))
	return group + "/" + path, nil
}

func (s *storage) download(fileID string, w io.Writer) error {
	ss := strings.SplitN(fileID, "/", 2)
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
	h := header{
		uint64(buf.Len()),
		STORAGE_PROTO_CMD_DOWNLOAD_FILE,
		0,
	}
	p := newProtocol(h, conn)
	b, err = p.request(buf.Bytes())
	if err != nil {
		return err
	}

	w.Write(b)
	return nil
}
