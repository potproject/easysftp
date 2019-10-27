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

// IOReaderProgress forwards the Read() call
// Addging transferredBytes
type IOReaderProgress struct {
	io.Reader
	TransferredBytes *int64 // Total of bytes transferred
}

func (iorp *IOReaderProgress) Read(p []byte) (int, error) {
	n, err := iorp.Reader.Read(p)
	*iorp.TransferredBytes += int64(n)
	return n, err
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
	return getTransfer(esftp.SFTPClient, localFilepath, remoteFilepath, nil)
}

// GetWithProgress [Experimental] Get with Display Processing Bytes
func (esftp Easysftp) GetWithProgress(localFilepath string, remoteFilepath string, transferred *int64) (int64, error) {
	return getTransfer(esftp.SFTPClient, localFilepath, remoteFilepath, transferred)
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
		localFilepath, _ := getRecursivelyPath(localPath, remotePath, remoteFullFilepath)

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
		_, getErr := getTransfer(esftp.SFTPClient, localFilepath, remoteFullFilepath, nil)
		if getErr != nil {
			return getErr
		}
	}
	return nil
}

// getTransfer Download Transfer execute
func getTransfer(client *sftp.Client, localFilepath string, remoteFilepath string, tfBytes *int64) (int64, error) {
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

	var bytes int64
	var copyErr error
	// withProgress
	if tfBytes != nil {
		remoteFileWithProgress := &IOReaderProgress{Reader: remoteFile, TransferredBytes: tfBytes}
		bytes, copyErr = io.Copy(localFile, remoteFileWithProgress)
	} else {
		bytes, copyErr = io.Copy(localFile, remoteFile)
	}
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
	return putTransfer(esftp.SFTPClient, localFilepath, remoteFilepath, nil)
}

// PutWithProgress [Experimental] Put with Display Processing Bytes
func (esftp Easysftp) PutWithProgress(localFilepath string, remoteFilepath string, transferred *int64) (int64, error) {
	return putTransfer(esftp.SFTPClient, localFilepath, remoteFilepath, transferred)
}

// PutRecursively is Recursively Upload entire directories
func (esftp Easysftp) PutRecursively(localPath string, remotePath string) error {
	localWalkerErr := filepath.Walk(localPath, func(localFullFilepath string, info os.FileInfo, fileErr error) error {
		if fileErr != nil {
			return fileErr
		}
		remoteFilepath, _ := getRecursivelyPath(remotePath, localPath, localFullFilepath)
		// if Not Exist Mkdir
		if info.IsDir() {
			remoteStat, remoteStatErr := esftp.SFTPClient.Stat(remoteFilepath)
			// 存在するかつディレクトリではない場合エラー
			if !os.IsNotExist(remoteStatErr) && !remoteStat.IsDir() {
				return errors.New("Cannot create a directry when that file already exists")
			}
			mode := info.Mode()
			if os.IsNotExist(remoteStatErr) {
				mkErr := esftp.SFTPClient.Mkdir(remoteFilepath)
				if mkErr != nil {
					return mkErr
				}
				chErr := esftp.SFTPClient.Chmod(remoteFilepath, mode)
				if chErr != nil {
					return chErr
				}
			}
			return nil
		}
		_, getErr := putTransfer(esftp.SFTPClient, localFullFilepath, remoteFilepath, nil)
		if getErr != nil {
			return getErr
		}
		return nil
	})
	return localWalkerErr
}

// Upload Transfer execute
func putTransfer(client *sftp.Client, localFilepath string, remoteFilepath string, tfBytes *int64) (int64, error) {
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

	var bytes int64
	var copyErr error
	// withProgress
	if tfBytes != nil {
		localFileWithProgress := &IOReaderProgress{Reader: localFile, TransferredBytes: tfBytes}
		bytes, copyErr = io.Copy(remoteFile, localFileWithProgress)
	} else {
		bytes, copyErr = io.Copy(remoteFile, localFile)
	}

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

func getRecursivelyPath(localPath string, remotePath string, remoteFullFilepath string) (string, error) {
	rel, err := filepath.Rel(filepath.Clean(remotePath), remoteFullFilepath)
	if err != nil {
		return "", err
	}
	localFilepath := filepath.Join(localPath, rel)
	return filepath.ToSlash(localFilepath), nil
}

// Quit alias Close()
func (esftp Easysftp) Quit() (errors []error) {
	return esftp.Close()
}
