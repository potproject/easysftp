# potproject/easysftp
SFTP connection easy. [pkg/sftp](https://github.com/pkg/sftp) wrapper Golang Library

## Usage

```sh
import "github.com/potproject/easysftp"
```

## Example
```go
package main

import (
	"log"

	"github.com/potproject/easysftp"
)

func main() {
	// Connect Example
	// [SFTP Command] $ sftp USERNAME@example.hostname.local -oPort=22 -i ~/.ssh/id_rsa
	esftpSession, err := easysftp.Connect("USERNAME", "example.hostname.local", 22, "~/.ssh/id_rsa")

	// Alternative: Using *ssh.Client
	// esftpSession, err := easysftp.NewClient(conn)

	if err != nil {
		log.Fatalln(err.Error())
	}
	defer esftpSession.Close()

	// SFTP File Get
	// [SFTP Command] sftp> get /tmp/remotefile.txt /tmp/localfile.txt
	downloadBytes, downloadError := esftpSession.Get("/tmp/localfile.txt", "/tmp/remotefile.txt")
	if downloadError != nil {
		log.Fatalln("Download Error:", err.Error())
	}
	log.Println("Download OK:", downloadBytes)

	// SFTP File Put
	// [SFTP Command] sftp> put /tmp/localfile.txt /tmp/remotefile.txt
	uploadBytes, uploadError := esftpSession.Put("/tmp/localfile.txt", "/tmp/remotefile.txt")
	if uploadError != nil {
		log.Fatalln("Upload Error:", err.Error())
	}
	log.Println("Upload OK:", uploadBytes)
}
```

## LICENSE
MIT
