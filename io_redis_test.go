package ioredis

import (
	"fmt"
	"os"
	"testing"
)

var client Client

func init() {
	conn, err := client.connect()

	if err != nil {
		fmt.Fprintf(os.Stdout, "Redis connect error %s", err)
	}

	client.conn = conn
}

func TestBasic(t *testing.T) {
	_, err := client.rawSend(client.conn, []byte("*2\r\n$3\r\nget\r\n$5\r\nhello\r\n"))

	if err != nil {
		t.Fatalf("Redis connect error %s", err)
	}
}
