package main

import (
    "log"
    "github.com/CunTianXing/go_app/go-grpc/v6/api"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
)

// Authentication holds the login/password
type Authentication struct {
    Login     string
    Password  string
}

// GetRequestMetadata gets the current request metadata
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
     return map[string]string{
         "login":    a.Login,
         "password": a.Password,
    }, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security
func (a *Authentication) RequireTransportSecurity() bool {
    return true
}

func main() {
    var conn *grpc.ClientConn
    // Create the  client TLS credentials
    creds, err := credentials.NewClientTLSFromFile("cert/server.crt","")
    if err != nil {
        log.Fatalf("could not load tls cert: %s", err)
    }
    // Setup the login/pass
    auth := Authentication{
        Login:    "xingcuntian",
        Password: "xingcuntian",
    }
    conn, err = grpc.Dial("localhost:7777",grpc.WithTransportCredentials(creds),grpc.WithPerRPCCredentials(&auth)) 
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
