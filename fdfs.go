package fdfs

import (
	"errors"
	"io"
	"time"
)

var timeout = 10 * time.Second

type Client struct {
	trackers []*Tracker
}

func New() *Client {
	return &Client{}
}

func (c *Client) AddTracker(host string, port int) {
	c.trackers = append(c.trackers, &Tracker{
		host,
		port,
	})
}

func (c *Client) SetTimeout(t time.Duration) {
	timeout = t
}

func (c *Client) Upload(file io.Reader) (string, error) {
	tracker, err := c.getTracker()
	if err != nil {
		return "", err
	}
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

func (c *Client) Download(fileId string, w io.Writer) error {
	tracker, err := c.getTracker()
	if err != nil {
		return err
	}
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

func (c *Client) getTracker() (*Tracker, error) {
	if len(c.trackers) == 0 {
		return nil, errors.New("没有添加tracker")
	}
	tracker := c.trackers[0]
	return tracker, nil
}
