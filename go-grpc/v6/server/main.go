package main

import (
    "fmt"
    "strings"
    "log"
    "net"
    "net/http"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "github.com/CunTianXing/go_app/go-grpc/v6/api"
    "google.golang.org/grpc"
    "golang.org/x/net/context"
    "golang.org/x/net/trace"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/grpclog"
)

// private type for Context keys
type contextKey int

const (
    clientIDKey contextKey = iota
)

func credMatcher(headerName string) (mdName string, ok bool) {
    if headerName == "Login" || headerName == "Password" {
        return headerName, true
    }
    return "", false
}

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

func startGRPCServer(address, certFile, keyFile string) error {

    lis, err := net.Listen("tcp",address) 
    if err != nil {
       return fmt.Errorf("failed to listen: %v",err)
    }

    // create a server instance
    s := api.Server{}

    //Create the TLS credentials
    creds, err := credentials.NewServerTLSFromFile("cert/server.crt","cert/server.key")
    if err != nil {
        return fmt.Errorf("could not load TLS keys: %s", err)
    }
    // Create an array of gRPC options with the credentials
    opts := []grpc.ServerOption{grpc.Creds(creds),grpc.UnaryInterceptor(unaryInterceptor)} 
    //create a gRPC server object
    grpcServer := grpc.NewServer(opts...)
    // attach the Ping service to the server
    api.RegisterPingServer(grpcServer, &s)
    //go startTraceServer(":50051")
    // start the server    
    if err := grpcServer.Serve(lis); err != nil {
        return fmt.Errorf("failed to serve: %s", err)
    }
    return nil
  }

  func startTraceServer(address string) {
      trace.AuthRequest = func(req *http.Request) (any,sensitive bool) {
          return true, true
      }
      go http.ListenAndServe(address,nil)
      grpclog.Println("Trace Listen on "+address)
  }

  func startRESTServer(address, grpcAddress, certFile string) error {
      ctx := context.Background()
      ctx, cancel := context.WithCancel(ctx)
      defer cancel()

      mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(credMatcher))
      creds, err := credentials.NewClientTLSFromFile(certFile,"")
      if err != nil {
          return fmt.Errorf("could not load TLS certificate: %s", err)
      }

      opts :=[]grpc.DialOption{grpc.WithTransportCredentials(creds)}

      err = api.RegisterPingHandlerFromEndpoint(ctx,mux,grpcAddress,opts)
      if err != nil {
          return fmt.Errorf("could not register service Ping: %s", err)
      }
      log.Printf("starting HTTP/1.1 REST server on %s", address)
      http.ListenAndServe(address,mux)
      return nil
  }

  func main() {
      grpcAddress := fmt.Sprintf("%s:%d","localhost",7777)
      restAddress := fmt.Sprintf("%s:%d","localhost",7778)
      certFile := "cert/server.crt"
      keyFile := "cert//server.key"
      go startTraceServer(":50051")
      go func() {
        err := startGRPCServer(grpcAddress, certFile, keyFile)
        if err != nil {
            log.Fatalf("failed to start gRPC server: %s", err)
        }
      }()

      go func() {
        err := startRESTServer(restAddress, grpcAddress, certFile)
        if err != nil {
            log.Fatalf("failed to start gRPC server: %s", err)
        }
      }()
      log.Printf("Entering infinite loop")
      select {}
  }























