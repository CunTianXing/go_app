package main

import (
    "log"
    "time"
    hello "github.com/CunTianXing/go_app/go-micro/grpc_demo_v1/srv/proto/hello"
    "github.com/micro/go-grpc"
    "github.com/micro/go-micro"
    "golang.org/x/net/context"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
    log.Print("Received say.hello request")
    rsp.Msg = "Hello " + req.Name
    return nil
}

func main() {
    service := grpc.NewService(
        micro.Name("go.micro.srv.greeter"),
        micro.RegisterTTL(time.Second*30),
        micro.RegisterInterval(time.Second*10),
    )

    service.Init()

    hello.RegisterSayHandler(service.Server(), new(Say))

    if err := service.Run(); err != nil {
        log.Fatal(err)
    }
}
