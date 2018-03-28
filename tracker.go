package fdfs

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
)

type Tracker struct {
	host string
	port int
}

func (t Tracker) getUploadStorage() (*Storage, error) {
	address := fmt.Sprintf("%s:%d", t.host, t.port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	trackerReqHeader := &Header{
		length:  0,
		command: TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE,
		status:  0,
	}
	conn.Write(trackerReqHeader.encode())

	bs := make([]byte, 1024)
	n, err := conn.Read(bs)

	if err != nil {
		return nil, err
	}
	trackerResp := bs[:n]

	// 解析header
	trackerRespHeader := &Header{
		buf: trackerResp[:10],
	}
	trackerRespHeader.decode()
	if trackerRespHeader.length != 40 || trackerRespHeader.status != 0 {
		return nil, errors.New("[tracker]状态码错误")
	}

	trackerRespBody := trackerResp[10:]
	group := string(trackerRespBody[:16])
	host := string(trackerRespBody[16 : 16+15])
	port := int(binary.BigEndian.Uint64(trackerRespBody[16+15 : 16+15+8]))
	index := trackerRespBody[16+15+8 : 16+15+8+1][0]

	return &Storage{
		group: group,
		host:  host,
		port:  port,
		index: index,
	}, nil
}

func (t Tracker) getDownloadStorage(fileId string) (*Storage, error) {
	address := fmt.Sprintf("%s:%d", t.host, t.port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	ss := strings.SplitN(fileId, "/", 2)
	group := ss[0]
	path := ss[1]
	trackerReqHeader := &Header{
		length:  uint64(FDFS_GROUP_NAME_MAX_LEN + len(path)),
		command: TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE,
		status:  0,
	}

	conn.Write(trackerReqHeader.encode())
	r := strings.NewReader(group)
	b := make([]byte, 16)
	r.Read(b)
	conn.Write(b)
	conn.Write([]byte(path))

	b = make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		return nil, err
	}
	trackerResp := b[:n]

	// 解析header
	trackerRespHeader := &Header{
		buf: trackerResp[:10],
	}
	trackerRespHeader.decode()
	if trackerRespHeader.length != 39 || trackerRespHeader.status != 0 {
		return nil, errors.New("[tracker]状态码错误")
	}

	trackerRespBody := trackerResp[10:]
	group = string(trackerRespBody[:16])
	host := string(trackerRespBody[16 : 16+15])
	port := int(binary.BigEndian.Uint64(trackerRespBody[16+15 : 16+15+8]))

	return &Storage{
		group: group,
		host:  host,
		port:  port,
	}, nil
}
