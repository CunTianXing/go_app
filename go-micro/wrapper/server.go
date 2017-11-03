package main

import (
	"fmt"
	"log"

	proto "github.com/CunTianXing/go_app/go-micro/service/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Message = "Hello " + req.Name
	return nil
}

// server middleware
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Printf("[wrapper] server request: %v", req.Method())
		
		err := fn(ctx, req, rsp)
		return err
	}
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
		micro.WrapHandler(logWrapper),
	)
	service.Init()
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
