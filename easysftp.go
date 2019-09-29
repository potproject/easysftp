package easysftp

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// File Local and remote filepath
type File struct {
	LocalFilepath  string
	RemoteFilepath string
}

// Report execute bytes and Error Report
type Report struct {
	Bytes int64
	Error error
}

func main() {
	return
}

// Download is Single File Download
func Download(username string, host string, port uint16, keyPath string, localFilepath string, remoteFilepath string) Report {
	conn, client, connErr := connect(username, host, port, keyPath)
	if connErr != nil {
		return Report{Bytes: 0, Error: errors.New("connErr: " + connErr.Error())}
	}
	defer conn.Close()
	defer client.Close()
	return downloadTransfer(client, localFilepath, remoteFilepath)
}

// DownloadMultiple is Multiple Files Download
func DownloadMultiple(username string, host string, port uint16, keyPath string, files []File) ([]Report, error) {
	var reports []Report
	conn, client, connErr := connect(username, host, port, keyPath)
	if connErr != nil {
		return reports, errors.New("connErr: " + connErr.Error())
	}
	defer conn.Close()
	defer client.Close()
	for _, file := range files {
		report := downloadTransfer(client, file.LocalFilepath, file.RemoteFilepath)
		reports = append(reports, report)
	}
	return reports, nil
}

// Download Transfer execute
func downloadTransfer(client *sftp.Client, localFilepath string, remoteFilepath string) Report {
	localFile, localFileErr := os.Create(localFilepath)
	if localFileErr != nil {
		return Report{Bytes: 0, Error: errors.New("localFileErr: " + localFileErr.Error())}
	}
	defer localFile.Close()

	remoteFile, remoteFileErr := client.Open(remoteFilepath)
	if remoteFileErr != nil {
		return Report{Bytes: 0, Error: errors.New("remoteFileErr: " + remoteFileErr.Error())}
	}

	bytes, copyErr := io.Copy(localFile, remoteFile)
	if copyErr != nil {
		return Report{Bytes: 0, Error: errors.New("copyErr: " + copyErr.Error())}
	}

	syncErr := localFile.Sync()
	if syncErr != nil {
		return Report{Bytes: 0, Error: errors.New("syncErr: " + syncErr.Error())}
	}
	return Report{Bytes: bytes, Error: nil}
}

// Upload is Single File Upload
func Upload(username string, host string, port uint16, keyPath string, localFilepath string, remoteFilepath string) Report {
	conn, client, connErr := connect(username, host, port, keyPath)
	if connErr != nil {
		return Report{Bytes: 0, Error: errors.New("connErr: " + connErr.Error())}
	}
	defer conn.Close()
	defer client.Close()

	return uploadTransfer(client, localFilepath, remoteFilepath)
}

// UploadMultiple is Multiple File Upload
func UploadMultiple(username string, host string, port uint16, keyPath string, files []File) ([]Report, error) {
	var reports []Report
	conn, client, connErr := connect(username, host, port, keyPath)
	if connErr != nil {
		return reports, errors.New("connErr: " + connErr.Error())
	}
	defer conn.Close()
	defer client.Close()
	for _, file := range files {
		report := uploadTransfer(client, file.LocalFilepath, file.RemoteFilepath)
		reports = append(reports, report)
	}
	return reports, nil
}

// Upload Transfer execute
func uploadTransfer(client *sftp.Client, localFilepath string, remoteFilepath string) Report {
	remoteFile, remoteFileErr := client.Create(remoteFilepath)
	if remoteFileErr != nil {
		return Report{Bytes: 0, Error: errors.New("remoteFileErr: " + remoteFileErr.Error())}
	}
	defer remoteFile.Close()

	localFile, localFileErr := os.Open(localFilepath)
	if localFileErr != nil {
		return Report{Bytes: 0, Error: errors.New("localFileErr: " + localFileErr.Error())}
	}

	bytes, copyErr := io.Copy(remoteFile, localFile)
	if copyErr != nil {
		return Report{Bytes: 0, Error: errors.New("copyErr: " + copyErr.Error())}
	}
	return Report{Bytes: bytes, Error: nil}
}

// SSH Connection
func connect(username string, host string, port uint16, keyPath string) (conn *ssh.Client, client *sftp.Client, err error) {
	privateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		return
	}
	clientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err = ssh.Dial("tcp", host+":"+strconv.Itoa(int(port)), clientConfig)
	if err != nil {
		return
	}

	client, err = sftp.NewClient(conn)
	return
}
