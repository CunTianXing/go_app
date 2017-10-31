package main

import (
    "log"
    "github.com/CunTianXing/go_app/go-grpc/v1/api"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
)

func main() {
    var conn *grpc.ClientConn

    conn, err := grpc.Dial(":7777", grpc.WithInsecure())
    if err != nil {
        log.Fatalf("did not connect: %s", err)
    }
    defer conn.Close()

    c := api.NewPingClient(conn)
    response, err := c.SayHello(context.Background(), &api.PingMessage{Message:"Gary"})
    if err != nil {
        log.Fatalf("Error when calling sayHello:%s",err)
    }
    log.Printf("Response from server: %s",response.Message)
}
