package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "strings"
)

func main() {
    fmt.Println("TCP server....")

    ln, err := net.Listen("tcp",":1433")
    if err != nil {
        log.Fatal(err)
    }

    conn, err := ln.Accept()
    if err != nil {
        log.Fatal(err)
    }

    for {
        message, _ := bufio.NewReader(conn).ReadString('\n')
        fmt.Print("Message Received: ", string(message))
        newMessage := strings.ToTitle(message)
        conn.Write([]byte(newMessage + "\n"))
    }
}
