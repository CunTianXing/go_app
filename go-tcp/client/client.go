package main

import (
    "bufio"
    "fmt"
    "log"
    "net"
    "os"
)

func main() {
    conn, err := net.Dial("tcp","127.0.0.1:1433")
    if err != nil {
        log.Fatal(err)
    }
    for {
        reader := bufio.NewReader(os.Stdin)
        fmt.Print("Send messages to TCP server: ")
        text, _ := reader.ReadString('\n')
        fmt.Fprintf(conn,text+"\n")
        message, _ := bufio.NewReader(conn).ReadString('\n')
        fmt.Print("Reply from server: " + message)
    }
}
