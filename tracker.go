package fdfs

import (
	"encoding/binary"
	"fmt"
	"net"
)

type Tracker struct {
	host string
	port int
}

func (t Tracker) getStorage(command byte) *Storage {
	address := fmt.Sprintf("%s:%d", t.host, t.port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}

	trackerReqHeader := &Header{
		length:  0,
		command: command,
		status:  0,
	}
	conn.Write(trackerReqHeader.encode())

	bs := make([]byte, 1024)
	n, err := conn.Read(bs)

	if err != nil {
		panic(err)
	}
	trackerResp := bs[:n]

	// 解析header
	trackerRespHeader := &Header{
		buf: trackerResp[:10],
	}
	trackerRespHeader.decode()
	if trackerRespHeader.length != 40 || trackerRespHeader.status != 0 {
		panic("内部错误")
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
	}
}
