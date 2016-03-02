package ioredis

import (
	"fmt"
	"os"
	"testing"
)

var client Client

func init() {
	conn, err := client.Connect()

	if err != nil {
		fmt.Fprintf(os.Stdout, "Redis connect error %s", err)
	}

	client.Conn = conn
}

func TestBasic(t *testing.T) {
	_, err := client.Set("key", "hello world!")

	if err != nil {
		t.Fatal("Redis SET error")
	}

	value, err := client.Get("key")

	if err != nil {
		t.Fatal("Redis GET error")
	}

	if string(string(value.([]byte))) != "hello world!" {
		t.Fatal("Redis key value error")
	}
}

func BenchmarkGet(b *testing.B) {
	for i := 1; i < b.N; i++ {
		client.Set(string(i), "hello world!")
	}
}

func BenchmarkSet(b *testing.B) {
	for i := 1; i < b.N; i++ {
		client.Set("1", "hello world!")
	}
}
