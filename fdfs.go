package fdfs

import (
	"errors"
	"io"
	"net"
	"time"
)

var timeout = 10 * time.Second

type client struct {
	trackers []*tracker
}

// New client
func New() *client {
	return &client{}
}

func (c *client) AddTracker(trackerAddr string) error {
	host, port, err := net.SplitHostPort(trackerAddr)
	if err != nil {
		return err
	}
	c.trackers = append(c.trackers, &tracker{
		host,
		port,
	})
	return nil
}

func (c *client) SetTimeout(t time.Duration) {
	timeout = t
}

func (c *client) Upload(file io.Reader) (string, error) {
	tracker, err := c.getTracker()
	if err != nil {
		return "", err
	}
	storage, err := tracker.getUploadStorage()
	if err != nil {
		return "", err
	}
	fileID, err := storage.upload(file)
	if err != nil {
		return "", err
	}
	return fileID, nil
}

func (c *client) Download(fileID string, w io.Writer) error {
	tracker, err := c.getTracker()
	if err != nil {
		return err
	}
	storage, err := tracker.getDownloadStorage(fileID)
	if err != nil {
		return err
	}
	err = storage.download(fileID, w)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) getTracker() (*tracker, error) {
	if len(c.trackers) == 0 {
		return nil, errors.New("没有添加tracker")
	}
	tracker := c.trackers[0]
	return tracker, nil
}
