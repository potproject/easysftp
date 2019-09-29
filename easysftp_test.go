package easysftp

// This test is Production test. Not UnitTest
// Setting environment values actually sftp server.
// ENVIRONMENT:
// ENV["EASYSFTP_TEST_USERNAME"] 	: sftp server Username
// ENV["EASYSFTP_TEST_HOST"] 		: sftp server Hostname
// ENV["EASYSFTP_TEST_PORT"] 		: sftp server SSH PORT
// ENV["EASYSFTP_TEST_FILEPATH"] 	: sftp server rsa OpenSSH key FilePath

import (
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
	conn, client, err := connect(username, host, uint16(port), keyPath)
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer conn.Close()
	defer client.Close()
}
