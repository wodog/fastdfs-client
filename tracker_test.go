package fdfs

import (
	"fmt"
	"os"
	"testing"
)

func TestUpload(t *testing.T) {
	// tracker := &Tracker{"zpbeer.com", 22122}
	// storage := tracker.getStorage()
	// fmt.Println(storage)

	client := &Client{
		tracker_host: "zpbeer.com",
		tracker_port: 22122,
	}
	file, _ := os.Open("w")
	fileId := client.Upload(file)
	fmt.Println(fileId)
}
