package fdfs

//
// import (
// 	"bytes"
// 	"encoding/binary"
// 	"fmt"
// 	"io/ioutil"
// 	"net"
// 	"os"
// 	"time"
// )
//
// func main() {
//
// 	file, err := os.Open("w")
// 	if err != nil {
// 		panic(err)
// 	}
// 	fileStat, err := file.Stat()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(fileStat.Size())
//
// 	// bytes.NewBuffer(make([]byte, 25))
// 	//
// 	// b := make([]byte, 8)
// 	// binary.BigEndian.PutUint64(b, 256)
// 	// fmt.Println(b)
// 	// // convInt := binary.BigEndian.Uint32([]byte{0, 0, 0, 0, 0, 0x59, 0xd8})
// 	// // fmt.Println(convInt)
//
// 	conn, err := net.Dial("tcp", "zpbeer.com:22122")
// 	if err != nil {
// 		panic(err)
// 	}
// 	// buffer := bytes.NewBuffer(make([]byte, 0))
// 	buffer := &bytes.Buffer{}
// 	length := make([]byte, 8)
// 	binary.BigEndian.PutUint64(length, 0)
// 	buffer.Write(length)  // length
// 	buffer.WriteByte(101) // cmd
// 	buffer.WriteByte(0)   // status
// 	conn.Write(buffer.Bytes())
//
// 	trackerResp := make([]byte, 1024)
// 	n, err := conn.Read(trackerResp)
// 	if err != nil {
// 		panic(err)
// 	}
// 	trackerResp = trackerResp[:n]
// 	fmt.Println(trackerResp)
// 	// trackerRespHeader := trackerResp[0:10]
// 	trackerRespBody := trackerResp[10:]
// 	group := string(trackerRespBody[:16])
// 	ip := string(trackerRespBody[16 : 16+15])
// 	port := binary.BigEndian.Uint64(trackerRespBody[16+15 : 16+15+8])
// 	index := int(trackerRespBody[16+15+8 : 16+15+8+1][0])
// 	fmt.Println(group)
// 	fmt.Println(ip)
// 	fmt.Println(port)
// 	fmt.Println(index)
//
// 	conn, err = net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	buffer = &bytes.Buffer{}
// 	length = make([]byte, 8)
// 	binary.BigEndian.PutUint64(length, uint64(fileStat.Size()+15))
// 	buffer.Write(length)
// 	buffer.WriteByte(11)
// 	buffer.WriteByte(0)
// 	buffer.WriteByte(byte(index))
// 	length = make([]byte, 8)
// 	binary.BigEndian.PutUint64(length, uint64(fileStat.Size()))
// 	buffer.Write(length)
// 	length = make([]byte, 8)
// 	binary.BigEndian.PutUint64(length, 0)
// 	buffer.Write(length)
//
// 	conn.Write(buffer.Bytes())
//
// 	time.Sleep(20 * time.Second)
//
// 	b, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		panic(err)
// 	}
// 	conn.Write(b)
//
// 	b = make([]byte, 1024)
// 	n, err = conn.Read(b)
// 	if err != nil {
// 		panic(err)
// 	}
// 	storageResp := b[:n]
// 	storageRespBody := storageResp[10:]
// 	fmt.Println(string(storageRespBody[:16]))
// 	fmt.Println(string(storageRespBody[16:]))
// }
