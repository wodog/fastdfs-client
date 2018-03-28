package fdfs

import (
	"io"
)

type Client struct {
	tracker_host string
	tracker_port int
}

func (c Client) Upload(file io.Reader) (string, error) {
	tracker := &Tracker{c.tracker_host, c.tracker_port}
	storage, err := tracker.getUploadStorage()
	if err != nil {
		return "", err
	}
	fileId, err := storage.upload(file)
	if err != nil {
		return "", err
	}
	return fileId, nil
}

func (c Client) Download(fileId string, w io.Writer) error {
	tracker := &Tracker{c.tracker_host, c.tracker_port}
	storage, err := tracker.getDownloadStorage(fileId)
	if err != nil {
		return err
	}
	err = storage.download(fileId, w)
	if err != nil {
		return err
	}
	return nil
}
