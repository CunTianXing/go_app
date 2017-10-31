package main

import (
    "log"
    "github.com/CunTianXing/go_app/go-grpc/v3/api"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

func main() {
    var conn *grpc.ClientConn
    // Create the  client TLS credentials
    creds, err := credentials.NewClientTLSFromFile("cert/server.crt","localhost")
    if err != nil {
        log.Fatalf("could not load tls cert: %s", err)
    }
    conn, err = grpc.Dial("localhost:7777", grpc.WithTransportCredentials(creds))
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
