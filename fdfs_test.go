package fdfs

import (
	"fmt"
	"log"
	"os"
	"testing"
)

var fileID string
var fileName = "README.md"
var trackerAddr = "zpbeer.com:22122"
var client = Default()

func init() {
	err := client.AddTracker(trackerAddr)
	if err != nil {
		panic(err)
	}
}

func TestUploadDownloadDelete(t *testing.T) {
	for i := 0; i < 1; i++ {
		func() {
			success := t.Run("upload", upload)
			if !success {
				return
			}

			success = t.Run("download", download)
			if success {
				t.Run("info", info)
			}

			t.Run("delete", delete)
		}()
	}
}

func upload(t *testing.T) {
	file, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}

	fileID, err = client.Upload(file)
	fmt.Println(fileID)
	if err != nil {
		t.Fatal(err)
	}
}

func download(t *testing.T) {
	err := client.Download(fileID, os.Stdout)
	if err != nil {
		t.Fatal(err)
	}
}

func delete(t *testing.T) {
	err := client.Delete(fileID)
	if err != nil {
		t.Fatal(err)
	}
}

func info(t *testing.T) {
	m, err := client.Info(fileID)
	log.Println(m)
	if err != nil {
		t.Fatal(err)
	}
}
