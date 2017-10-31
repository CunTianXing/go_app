package main

import (
    "flag"
    "net/http"
    "github.com/golang/glog"
    "github.com/grpc-ecosystem/grpc-gateway/runtime"
    "golang.org/x/net/context"
    "google.golang.org/grpc"
    hello "github.com/CunTianXing/go_app/go-micro/grpc_gateway_demo_v1/proto/hello"
)

var (
    endpoint = flag.String("endpoint","localhost:9090","go.micro.srv.greeter address")
)

func run() error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithInsecure()}
    err := hello.RegisterSayHandlerFromEndpoint(ctx,mux,*endpoint,opts)
    if err != nil {
        return err
    }
    return http.ListenAndServe(":8080",mux)
}

func main() {
    flag.Parse()
    defer glog.Flush()
    if err := run(); err != nil {
        glog.Fatal(err)
    }
}
