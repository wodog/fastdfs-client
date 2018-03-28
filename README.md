## fastdfs-client

go版的fastdfs客户端，实现了上传，下载功能

#### Install

```
go get github.com/wodog/fastdfs-client
```

#### Usage

```
package main

import "github.com/wodog/fastdfs-client"

func main() {
  client := &fdfs.Client{
    tracker_host: "xxx.com",
    tracker_port: 22122,
  }

  // upload file
  file, _ := os.Open(fileName)
  fileId := client.Upload(file)

  // download file
  client.Download(fileId, os.Stdout)
}
```

#### Reference

[协议参考](http://weakyon.com/2014/09/01/analysis-of-source-code-for-fastdfs.html)
