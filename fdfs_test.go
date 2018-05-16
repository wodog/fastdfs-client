package fdfs

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var fileID string
var fileName = "README.md"
var trackerAddr = "yxsm-test:22122"
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

			success = t.Run("open", open)
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
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("upload:", fileID)
}

func open(t *testing.T) {
	r, err := client.Open(fileID)
	if err != nil {
		t.Fatal(err)
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("open:", b)
}

func delete(t *testing.T) {
	err := client.Delete(fileID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("delete:", "success")
}

func info(t *testing.T) {
	m, err := client.Info(fileID)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("info:", m)
}
