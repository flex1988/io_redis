package ioredis

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

var defaultAddr = "127.0.0.1:6379"

type Client struct {
	Addr     string
	PoolSize int
	Conn     net.Conn
	Pool     chan net.Conn
}

type RedisError string

func (err RedisError) Error() string {
	return "Redis Error: " + string(err)
}

func (c *Client) Connect() (conn net.Conn, err error) {
	var addr = defaultAddr
	if c.Addr != "" {
		addr = c.Addr
	}
	conn, err = net.Dial("tcp", addr)
	return
}

func (c *Client) Send(cmd string, args ...string) (interface{}, error) {

	data, err := c.rawSend(c.Conn, cmd, args...)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) rawSend(conn net.Conn, cmd string, args ...string) (interface{}, error) {
	_, err := conn.Write(command(cmd, args...))

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

func command(cmd string, args ...string) []byte {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "*%d\r\n$%d\r\n%s\r\n", len(args)+1, len(cmd), cmd)

	for _, v := range args {
		fmt.Fprintf(&buf, "$%d\r\n%s\r\n", len(v), v)
	}

	return buf.Bytes()
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

	switch line[0] {
	case '+':
		return strings.TrimSpace(line[1:]), nil
	case ':':
		return strconv.Atoi(strings.TrimSpace(line[1:]))
	case '-':
		return nil, RedisError("ERR:" + strings.TrimSpace(line[5:]))
	case '$':
		return readBulk(reader, line)

		//TODO *
	default:
		return nil, RedisError("INVALID DATA TYPE")
	}
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

	if head[0] == '$' {
		size, err := strconv.Atoi(head[1:])

		if err != nil {
			return nil, err
		}

		if size == 0 {
			data = nil
		}

		str := io.LimitReader(reader, int64(size))

		data, err = ioutil.ReadAll(str)

		if err != nil {
			return nil, err
		}
	}
	return data, err
}
