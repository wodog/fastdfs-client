package fdfs

import (
	"io"
)

type Client struct {
	tracker_host string
	tracker_port int
}

func (c Client) Upload(file io.Reader) string {
	tracker := &Tracker{c.tracker_host, c.tracker_port}
	storage := tracker.getUploadStorage()
	fileId := storage.upload(file)
	return fileId
}

func (c Client) Download(fileId string, w io.Writer) {
	tracker := &Tracker{c.tracker_host, c.tracker_port}
	storage := tracker.getDownloadStorage(fileId)
	storage.download(fileId, w)
}
