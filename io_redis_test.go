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

	client.Conn = conn
}

func TestBasic(t *testing.T) {
	_, err := client.Send("SET", "key", "hello world!")

	if err != nil {
		t.Fatal("Redis SET error")
	}

	value, err := client.Send("GET", "key")

	if err != nil {
		t.Fatal("Redis GET error")
	}

	if string(string(value.([]byte))) != "hello world!" {
		t.Fatal("Redis key value error")
	}
}
