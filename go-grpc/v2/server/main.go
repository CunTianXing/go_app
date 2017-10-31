package main

import (
    "fmt"
    "log"
    "net"
    "github.com/CunTianXing/go_app/go-grpc/v2/api"
    "google.golang.org/grpc"
)

func main() {
    
    lis, err := net.Listen("tcp",fmt.Sprintf(":%d",7777))
    if err != nil {
        log.Fatalf("failed to listen: %v",err)
    }
    // create a server instance
    s := api.Server{}
    //create a gRPC server object
    grpcServer := grpc.NewServer()
    // attach the Ping service to the server
    api.RegisterPingServer(grpcServer, &s)
    // start the server
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %s", err)
    }
}

