package fdfs

import (
	"os"
	"testing"
)

// func TestUpload(t *testing.T) {
// 	client := New()
// 	client.AddTracker("zpbeer.com:22122")
// 	file, err := os.Open("README.md")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fileID, err := client.Upload(file)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(fileID)
// }

func TestDownload(t *testing.T) {
	client := New()
	err := client.AddTracker("zpbeer.com:22122")
	if err != nil {
		panic(err)
	}
	err = client.Download("group1/M00/00/00/eBg-z1q-1CWAOOR4AAACUPg8FDI1021631", os.Stdout)
	if err != nil {
		panic(err)
	}
}
