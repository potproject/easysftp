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

// Connect SSH Connection
func Connect(username string, host string, port uint16, keyPath string) (conn *ssh.Client, client *sftp.Client, err error) {
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

// Get is Single File Download
func Get(client *sftp.Client, localFilepath string, remoteFilepath string) (int64, error) {
	return getTransfer(client, localFilepath, remoteFilepath)
}

// GetMultiple is Multiple Files Download
func GetMultiple(client *sftp.Client, files []File) (int64s []int64, errors []error) {
	for _, file := range files {
		i, e := getTransfer(client, file.LocalFilepath, file.RemoteFilepath)
		int64s = append(int64s, i)
		errors = append(errors, e)
	}
	return
}

// getTransfer Download Transfer execute
func getTransfer(client *sftp.Client, localFilepath string, remoteFilepath string) (int64, error) {
	localFile, localFileErr := os.Create(localFilepath)
	if localFileErr != nil {
		return 0, errors.New("localFileErr: " + localFileErr.Error())
	}
	defer localFile.Close()

	remoteFile, remoteFileErr := client.Open(remoteFilepath)
	if remoteFileErr != nil {
		return 0, errors.New("remoteFileErr: " + remoteFileErr.Error())
	}
	defer remoteFile.Close()

	bytes, copyErr := io.Copy(localFile, remoteFile)
	if copyErr != nil {
		return 0, errors.New("copyErr: " + copyErr.Error())
	}

	syncErr := localFile.Sync()
	if syncErr != nil {
		return 0, errors.New("syncErr: " + syncErr.Error())
	}
	return bytes, nil
}

// Put is Single File Upload
func Put(client *sftp.Client, localFilepath string, remoteFilepath string) (int64, error) {
	return putTransfer(client, localFilepath, remoteFilepath)
}

// UploadMultiple is Multiple File Upload
func putMultiple(client *sftp.Client, files []File) (int64s []int64, errors []error) {
	for _, file := range files {
		i, e := putTransfer(client, file.LocalFilepath, file.RemoteFilepath)
		int64s = append(int64s, i)
		errors = append(errors, e)
	}
	return
}

// Upload Transfer execute
func putTransfer(client *sftp.Client, localFilepath string, remoteFilepath string) (int64, error) {
	remoteFile, remoteFileErr := client.Create(remoteFilepath)
	if remoteFileErr != nil {
		return 0, errors.New("remoteFileErr: " + remoteFileErr.Error())
	}
	defer remoteFile.Close()

	localFile, localFileErr := os.Open(localFilepath)
	if localFileErr != nil {
		return 0, errors.New("localFileErr: " + localFileErr.Error())
	}
	defer localFile.Close()

	bytes, copyErr := io.Copy(remoteFile, localFile)
	if copyErr != nil {
		return 0, errors.New("copyErr: " + copyErr.Error())
	}
	return bytes, nil
}
