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
	conn, client, err := Connect(username, host, uint16(port), keyPath)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer conn.Close()
	defer client.Close()
}

func TestGet(t *testing.T) {
	username := os.Getenv("EASYSFTP_TEST_USERNAME")
	host := os.Getenv("EASYSFTP_TEST_HOST")
	port, _ := strconv.Atoi(os.Getenv("EASYSFTP_TEST_PORT"))
	keyPath := os.Getenv("EASYSFTP_TEST_FILEPATH")
	conn, client, err := Connect(username, host, uint16(port), keyPath)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer conn.Close()
	defer client.Close()

	sess, err := conn.NewSession()
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

	_, downloadError := Get(client, "./test.txt", "/tmp/test.txt")
	if downloadError != nil {
		t.Error(downloadError.Error())
		return
	}

	if err = os.Remove("./test.txt"); err != nil {
		t.Error(err.Error())
		return
	}

	sess2, err := conn.NewSession()
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
	conn, client, err := Connect(username, host, uint16(port), keyPath)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer conn.Close()
	defer client.Close()

	file, err := os.OpenFile("./test.txt", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		t.Error(err.Error())
		return
	}

	fmt.Fprintln(file, "TestPut")
	file.Close()

	_, uploadError := Put(client, "./test.txt", "/tmp/test.txt")
	if uploadError != nil {
		t.Error(uploadError.Error())
		return
	}

	if err = os.Remove("./test.txt"); err != nil {
		t.Error(err.Error())
		return
	}

	sess, err := conn.NewSession()
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
