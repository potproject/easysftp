package easysftp

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// File Local and remote filepath
type File struct {
	LocalFilepath  string
	RemoteFilepath string
}

// Easysftp Stored Client val
type Easysftp struct {
	SSHClient  *ssh.Client
	SFTPClient *sftp.Client
}

// Connect SSH Connection
func Connect(username string, host string, port uint16, keyPath string) (esftp Easysftp, err error) {
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

	conn, err := ssh.Dial("tcp", host+":"+strconv.Itoa(int(port)), clientConfig)
	if err != nil {
		return
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return
	}
	esftp = Easysftp{
		SSHClient:  conn,
		SFTPClient: client,
	}
	return
}

// NewClient SSH Using Conection
func NewClient(conn *ssh.Client) (esftp Easysftp, err error) {
	client, err := sftp.NewClient(conn)
	if err != nil {
		return
	}
	esftp = Easysftp{
		SSHClient:  conn,
		SFTPClient: client,
	}
	return
}

// Get is Single File Download
func (esftp Easysftp) Get(localFilepath string, remoteFilepath string) (int64, error) {
	return getTransfer(esftp.SFTPClient, localFilepath, remoteFilepath)
}

// GetRecursively is Recursively Download entire directories
func (esftp Easysftp) GetRecursively(localPath string, remotePath string) error {
	remoteWalker := esftp.SFTPClient.Walk(remotePath)
	if remoteWalker == nil {
		return errors.New("SFTP Walker Error")
	}
	for remoteWalker.Step() {
		err := remoteWalker.Err()
		if err != nil {
			return err
		}
		remoteFullFilepath := remoteWalker.Path()
		localFilepath, _ := getLocalFilepath(localPath, remotePath, remoteFullFilepath)

		// if Not Exist Mkdir
		if remoteWalker.Stat().IsDir() {
			localStat, localStatErr := os.Stat(localFilepath)
			// 存在するかつディレクトリではない場合エラー
			if !os.IsNotExist(localStatErr) && !localStat.IsDir() {
				return errors.New("Cannot create a directry when that file already exists")
			}
			mode := remoteWalker.Stat().Mode()
			if os.IsNotExist(localStatErr) {
				mkErr := os.Mkdir(localFilepath, mode)
				if mkErr != nil {
					return mkErr
				}
			}
			continue
		}
		_, getErr := getTransfer(esftp.SFTPClient, localFilepath, remoteFullFilepath)
		if getErr != nil {
			return getErr
		}
	}
	return nil
}

func getLocalFilepath(localPath string, remotePath string, remoteFullFilepath string) (string, error) {
	rel, err := filepath.Rel(filepath.Clean(remotePath), remoteFullFilepath)
	if err != nil {
		return "", err
	}
	localFilepath := filepath.Join(localPath, rel)
	return localFilepath, nil
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
func (esftp Easysftp) Put(localFilepath string, remoteFilepath string) (int64, error) {
	return putTransfer(esftp.SFTPClient, localFilepath, remoteFilepath)
}

// UploadMultiple is Multiple File Upload
func (esftp Easysftp) putMultiple(files []File) (int64s []int64, errors []error) {
	for _, file := range files {
		i, e := putTransfer(esftp.SFTPClient, file.LocalFilepath, file.RemoteFilepath)
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

// Close Connection ALL Connection Close
func (esftp Easysftp) Close() (errors []error) {
	if esftp.SFTPClient != nil {
		sftpErr := esftp.SFTPClient.Close()
		if sftpErr != nil {
			errors = append(errors, sftpErr)
		}
	}
	if esftp.SSHClient != nil {
		sshErr := esftp.SSHClient.Close()
		if sshErr != nil {
			errors = append(errors, sshErr)
		}
	}
	return
}

// Quit alias Close()
func (esftp Easysftp) Quit() (errors []error) {
	return esftp.Close()
}
