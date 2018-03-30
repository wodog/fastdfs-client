package fdfs

import (
	"os"
	"testing"
)

// func TestUpload(t *testing.T) {
// 	client := New()
// 	client.AddTracker("zpbeer.com", 22122)
// 	file, err := os.Open("README.md")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fileId, err := client.Upload(file)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(fileId)
// }

func TestDownload(t *testing.T) {
	client := New()
	err := client.AddTracker("zpbeer.com:22122")
	if err != nil {
		panic(err)
	}
	err = client.Download("group1/M00/00/00/eBg-z1q7wNOAeQUsAAACLtpjbEY0261044", os.Stdout)
	if err != nil {
		panic(err)
	}
}
