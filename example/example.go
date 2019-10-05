package example

import (
	"log"

	"github.com/potproject/easysftp"
)

func main() {
	// Connect
	// [SFTP Command] $ sftp USERNAME@example.hostname.local -oPort=22 -i ~/.ssh/id_rsa
	conn, client, err := easysftp.Connect("USERNAME", "example.hostname.local", 22, "~/.ssh/id_rsa")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()
	defer client.Close()

	// SFTP File Get
	// [SFTP Command] sftp> get /tmp/remotefile.txt /tmp/localfile.txt
	downloadBytes, downloadError := easysftp.Get(client, "/tmp/localfile.txt", "/tmp/remotefile.txt")
	if downloadError != nil {
		log.Fatalln("Download Error:", err.Error())
	}
	log.Println("Downlaod OK:", downloadBytes)

	// SFTP File Put
	// [SFTP Command] sftp> put /tmp/localfile.txt /tmp/remotefile.txt
	uploadBytes, uploadError := easysftp.Put(client, "/tmp/localfile.txt", "/tmp/remotefile.txt")
	if uploadError != nil {
		log.Fatalln("Uplaod Error:", err.Error())
	}
	log.Println("Upload OK:", uploadBytes)
}
