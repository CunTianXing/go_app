package main

import (
	"log"

	proto "github.com/CunTianXing/go_app/go-micro/api_demo/rpc/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Example struct{}

type Foo struct{}

func (e *Example) Call(ctx context.Context, req *proto.CallRequest, resp *proto.CallResponse) error {
	log.Print("Received Example.Call request")
	if len(req.Name) == 0 {
		return errors.BadRequest("go.micro.api.example", "no content")
	}
	resp.Message = "your request " + req.Name
	return nil
}

func (f *Foo) Bar(ctx context.Context, req *proto.EmptyRequest, resp *proto.EmptyResponse) error {
	log.Print("Received Foo.Bar request")
	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.example"),
	)
	service.Init()
	proto.RegisterExampleHandler(service.Server(), new(Example))
	proto.RegisterFooHandler(service.Server(), new(Foo))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
