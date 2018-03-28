package fdfs

import (
	"fmt"
	"io"
)

type Client struct {
	tracker_host string
	tracker_port int
}

func (c Client) Upload(file io.Reader) string {
	tracker := &Tracker{c.tracker_host, c.tracker_port}
	storage := tracker.getStorage(TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE)
	fileId := storage.upload(file)
	return fileId
}

func (c Client) Download() {
	fmt.Println("TODO")
}
