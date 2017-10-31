package main

import (
    "fmt"
    "log"
    "net"
    "github.com/CunTianXing/go_app/go-grpc/v3/api"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

func main() {
    
    lis, err := net.Listen("tcp",fmt.Sprintf("%s:%d","localhost",7777))
    if err != nil {
        log.Fatalf("failed to listen: %v",err)
    }

    // create a server instance
    s := api.Server{}

    //Create the TLS credentials
    creds, err := credentials.NewServerTLSFromFile("cert/server.crt","cert/server.key")
    if err != nil {
        log.Fatalf("could not load TLS keys: %s", err)
    }
    // Create an array of gRPC options with the credentials
    opts := []grpc.ServerOption{grpc.Creds(creds)}
    //create a gRPC server object
    grpcServer := grpc.NewServer(opts...)
    // attach the Ping service to the server
    api.RegisterPingServer(grpcServer, &s)
    // start the server    
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %s", err)

    }
}

