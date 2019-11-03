package easysftp

// This test is Production test. Not UnitTest
// Setting environment values actually sftp server.
// ENVIRONMENT:
// ENV["EASYSFTP_TEST_USERNAME"] 	: sftp server Username
// ENV["EASYSFTP_TEST_HOST"] 		: sftp server Hostname
// ENV["EASYSFTP_TEST_PORT"] 		: sftp server SSH PORT
// ENV["EASYSFTP_TEST_FILEPATH"] 	: sftp server rsa OpenSSH key FilePath

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestInit(t *testing.T) {
	if os.Getenv("EASYSFTP_TEST_USERNAME") == "" {
		t.Error("EASYSFTP_TEST_USERNAME is not set environment")
	}
	if os.Getenv("EASYSFTP_TEST_HOST") == "" {
		t.Error("EASYSFTP_TEST_HOST is not set environment")
	}
	if os.Getenv("EASYSFTP_TEST_PORT") == "" {
		t.Error("EASYSFTP_TEST_PORT is not set environment")
	}
	if os.Getenv("EASYSFTP_TEST_FILEPATH") == "" {
		t.Error("EASYSFTP_TEST_FILEPATH is not set environment")
	}
}
func TestConnect(t *testing.T) {
	username := os.Getenv("EASYSFTP_TEST_USERNAME")
	host := os.Getenv("EASYSFTP_TEST_HOST")
	port, _ := strconv.Atoi(os.Getenv("EASYSFTP_TEST_PORT"))
	keyPath := os.Getenv("EASYSFTP_TEST_FILEPATH")
	esftp, err := Connect(username, host, uint16(port), keyPath)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer esftp.Close()
}

func TestGetRecursively(t *testing.T) {
	username := os.Getenv("EASYSFTP_TEST_USERNAME")
	host := os.Getenv("EASYSFTP_TEST_HOST")
	port, _ := strconv.Atoi(os.Getenv("EASYSFTP_TEST_PORT"))
	keyPath := os.Getenv("EASYSFTP_TEST_FILEPATH")
	esftp, err := Connect(username, host, uint16(port), keyPath)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer esftp.Close()

	sess, err := esftp.SSHClient.NewSession()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer sess.Close()

	err = sess.Run("mkdir -p /tmp/test && echo 'TestGetRecursively' > /tmp/test/test.txt && mkdir -p /tmp/test/test2 && echo 'TestGetRecursively2' > /tmp/test/test2/test2.txt")
	if err != nil {
		t.Error(err.Error())
		return
	}

	downloadError := esftp.GetRecursively("./testdirectory", "/tmp/test")
	if downloadError != nil {
		t.Error(downloadError.Error())
		return
	}

	if err = os.RemoveAll("./testdirectory"); err != nil {
		t.Error(err.Error())
		return
	}

	sess2, err := esftp.SSHClient.NewSession()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer sess2.Close()

	err = sess2.Run("rm -rf /tmp/test")
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestGet(t *testing.T) {
	username := os.Getenv("EASYSFTP_TEST_USERNAME")
	host := os.Getenv("EASYSFTP_TEST_HOST")
	port, _ := strconv.Atoi(os.Getenv("EASYSFTP_TEST_PORT"))
	keyPath := os.Getenv("EASYSFTP_TEST_FILEPATH")
	esftp, err := Connect(username, host, uint16(port), keyPath)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer esftp.Close()

	sess, err := esftp.SSHClient.NewSession()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer sess.Close()

	err = sess.Run("echo 'TestGet' > /tmp/test.txt")
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, downloadError := esftp.Get("./test.txt", "/tmp/test.txt")
	if downloadError != nil {
		t.Error(downloadError.Error())
		return
	}

	var i int64
	var a int64
	b, downloadWPError := esftp.GetWithProgress("./test.txt", "/tmp/test.txt", &i, &a)
	if i != b {
		t.Error("not matched transferred Bytes i:", i, "b:", b)
	}
	if i != a {
		t.Error("not matched transferred Bytes i:", i, "a:", a)
	}

	if downloadWPError != nil {
		t.Error(downloadWPError.Error())
		return
	}

	if err = os.Remove("./test.txt"); err != nil {
		t.Error(err.Error())
		return
	}

	sess2, err := esftp.SSHClient.NewSession()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer sess2.Close()

	err = sess2.Run("rm -rf /tmp/test.txt")
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestPut(t *testing.T) {
	username := os.Getenv("EASYSFTP_TEST_USERNAME")
	host := os.Getenv("EASYSFTP_TEST_HOST")
	port, _ := strconv.Atoi(os.Getenv("EASYSFTP_TEST_PORT"))
	keyPath := os.Getenv("EASYSFTP_TEST_FILEPATH")
	esftp, err := Connect(username, host, uint16(port), keyPath)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer esftp.Close()

	file, err := os.OpenFile("./test.txt", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Fprintln(file, "TestPut")
	file.Close()

	_, uploadError := esftp.Put("./test.txt", "/tmp/test.txt")
	if uploadError != nil {
		t.Error(uploadError.Error())
		return
	}

	var i int64
	var a int64
	b, uploadWPError := esftp.PutWithProgress("./test.txt", "/tmp/test.txt", &i, &a)
	if i != b {
		t.Error("not matched transferred Bytes i:", i, "b:", b)
	}
	if i != a {
		t.Error("not matched transferred Bytes i:", i, "a:", a)
	}
	if uploadWPError != nil {
		t.Error(uploadWPError.Error())
		return
	}

	if err = os.Remove("./test.txt"); err != nil {
		t.Error(err.Error())
		return
	}

	sess, err := esftp.SSHClient.NewSession()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer sess.Close()

	err = sess.Run("rm -rf /tmp/test.txt")
	if err != nil {
		t.Error(err.Error())
		return
	}
}

func TestPutRecursively(t *testing.T) {
	username := os.Getenv("EASYSFTP_TEST_USERNAME")
	host := os.Getenv("EASYSFTP_TEST_HOST")
	port, _ := strconv.Atoi(os.Getenv("EASYSFTP_TEST_PORT"))
	keyPath := os.Getenv("EASYSFTP_TEST_FILEPATH")
	esftp, err := Connect(username, host, uint16(port), keyPath)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer esftp.Close()

	err = os.MkdirAll("./testdirectory", 0777)
	if err != nil {
		t.Error(err.Error())
		return
	}

	file, err := os.OpenFile("./testdirectory/test.txt", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Fprintln(file, "TestPutRecursively")
	file.Close()

	err = os.MkdirAll("./testdirectory/test2", 0777)
	if err != nil {
		t.Error(err.Error())
		return
	}

	file2, err := os.OpenFile("./testdirectory/test2/test2.txt", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Fprintln(file2, "TestPutRecursively2")
	file2.Close()

	uploadError := esftp.PutRecursively("./testdirectory", "/tmp/test")
	if uploadError != nil {
		t.Error(uploadError.Error())
		return
	}

	if err = os.RemoveAll("./testdirectory"); err != nil {
		t.Error(err.Error())
		return
	}

	sess, err := esftp.SSHClient.NewSession()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer sess.Close()

	err = sess.Run("rm -rf /tmp/test")
	if err != nil {
		t.Error(err.Error())
		return
	}
}
