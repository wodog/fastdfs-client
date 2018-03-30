package fdfs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Tracker struct {
	host string
	port string
}

func (t Tracker) getUploadStorage() (*Storage, error) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", t.host, t.port), timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	trackerReqHeader := &Header{
		length:  0,
		command: TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE,
		status:  0,
	}
	conn.Write(trackerReqHeader.encode())

	b := make([]byte, 10)
	_, err = io.ReadFull(conn, b)
	if err != nil {
		return nil, err
	}
	header := &Header{
		buf: b,
	}
	header.decode()
	if header.status != 0 {
		return nil, errors.New("[tracker]状态码错误")
	}
	b = make([]byte, header.length)
	_, err = io.ReadFull(conn, b)
	if err != nil {
		return nil, err
	}
	group := clearZero(string(b[:16]))
	host := clearZero(string(b[16 : 16+15]))
	port := strconv.Itoa(int(binary.BigEndian.Uint64(b[16+15 : 16+15+8])))
	index := b[16+15+8 : 16+15+8+1][0]

	return &Storage{
		group: group,
		host:  host,
		port:  port,
		index: index,
	}, nil
}

func (t Tracker) getDownloadStorage(fileId string) (*Storage, error) {
	ss := strings.SplitN(fileId, "/", 2)
	group := ss[0]
	path := ss[1]

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", t.host, t.port), timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	buf := &bytes.Buffer{}
	b := make([]byte, 16)
	copy(b, group)
	buf.Write(b)
	buf.WriteString(path)
	header := &Header{
		length:  uint64(buf.Len()),
		command: TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE,
		status:  0,
	}
	conn.Write(header.encode())
	conn.Write(buf.Bytes())

	b = make([]byte, 10)
	_, err = io.ReadFull(conn, b)
	if err != nil {
		return nil, err
	}
	header = &Header{
		buf: b,
	}
	header.decode()
	if header.status != 0 {
		return nil, errors.New("[tracker]状态码错误")
	}
	b = make([]byte, header.length)
	_, err = io.ReadFull(conn, b)
	if err != nil {
		return nil, err
	}
	group = clearZero(string(b[:16]))
	host := clearZero(string(b[16 : 16+15]))
	port := strconv.FormatUint(binary.BigEndian.Uint64(b[16+15:16+15+8]), 10)

	return &Storage{
		group: group,
		host:  host,
		port:  port,
	}, nil
}
