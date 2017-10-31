package main

import (
    "fmt"
    "strings"
    "log"
    "net"
    "github.com/CunTianXing/go_app/go-grpc/v4/api"
    "google.golang.org/grpc"
    "golang.org/x/net/context"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/metadata"
)

// private type for Context keys
type contextKey int

const (
    clientIDKey contextKey = iota
)

// authenticateAgent check the client credentials
func authenticateClient(ctx context.Context, s *api.Server) (string, error) {
    if md, ok := metadata.FromIncomingContext(ctx); ok {
        fmt.Printf("ctx data: %+v\n",md)
        clientLogin := strings.Join(md["login"], "")
        clientPassword := strings.Join(md["password"],"")

        if clientLogin != "xingcuntian" {
            return "", fmt.Errorf("unknown user %s", clientLogin)
        }
        if clientPassword != "xingcuntian" {
            return "", fmt.Errorf("bad password %s", clientPassword)
        }
        log.Printf("authenticated client: %s", clientLogin)
        return "42", nil
    }
    return "", fmt.Errorf("missing credentials")
}

// unaryInterceptor calls authenticateClient with current context
func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    s, ok := info.Server.(*api.Server)
    if !ok {
        return nil, fmt.Errorf("unable to cast server")
    }
    clientID, err := authenticateClient(ctx, s)
    if err != nil {
        return nil, err
    }
    ctx = context.WithValue(ctx, clientIDKey, clientID)
    return handler(ctx,req)
}
func main() {

    lis, err := net.Listen("tcp",fmt.Sprintf("%s:%d","127.0.0.1",7777)) 
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
    opts := []grpc.ServerOption{grpc.Creds(creds),grpc.UnaryInterceptor(unaryInterceptor)} 
    //create a gRPC server object
    grpcServer := grpc.NewServer(opts...)
    // attach the Ping service to the server
    api.RegisterPingServer(grpcServer, &s)
    // start the server    
    if err := grpcServer.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %s", err)

    }
}

