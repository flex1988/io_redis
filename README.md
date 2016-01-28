## io_redis

    io_redis is a very simple client for redis.

## Installation

    `go get github.com/flex1988/io_redis`

## Example

    ``` go
    package main

    import (
      "fmt"
      "github.com/flex1988/io_redis"
      "os"
    )

    func main() {
      var client ioredis.Client

      client.Connect()

      _, err := client.Send("SET", "mykey", "hello world!")

      if err != nil {
      }

      str, _ := client.Send("GET", "mykey")

      fmt.Fprint(os.Stdout, "The value of the key mykey is ", string(str.([]byte)))
    }
    ```
## Test

    go test

## There's still a lot things todo
