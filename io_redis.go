package ioredis

import (
	"bufio"
	"net"
	"strconv"
	"strings"
)

var defaultAddr = "127.0.0.1:6379"

type Client struct {
	Addr     string
	PoolSize int
	conn     net.Conn
	pool     chan net.Conn
}

type RedisError string

func (err RedisError) Error() string {
	return "Redis Error: " + string(err)
}

func (c *Client) connect() (conn net.Conn, err error) {
	var addr = defaultAddr
	if c.Addr != "" {
		addr = c.Addr
	}
	conn, err = net.Dial("tcp", addr)
	return
}

func (c *Client) rawSend(conn net.Conn, cmd []byte) (interface{}, error) {
	_, err := conn.Write(cmd)

	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)

	data, err := readResponse(reader)

	if err != nil {
		return nil, err
	}
	return data, nil
}

func readResponse(reader *bufio.Reader) (interface{}, error) {
	var line string
	var err error

	for {
		line, err = reader.ReadString('\n')

		if len(line) == 0 || err != nil {
			return line, err
		}
		line = strings.TrimSpace(line)

		if len(line) > 0 {
			break
		}
	}

	if line[0] == '+' {
		return strings.TrimSpace(line[1:]), nil
	} else if line[0] == '$' {
		return strconv.Atoi(strings.TrimSpace(line[1:]))
	}

	return readBulk(reader, line)
}

func readBulk(reader *bufio.Reader, head string) ([]byte, error) {
	var data []byte
	var err error

	if head == "" {
		head, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
	}

	switch head[0] {
	case ':':
		data = []byte(strings.TrimSpace(head[1:]))
	}
	return data, err
}
