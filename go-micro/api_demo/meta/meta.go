package main

import (
	"log"

	proto "github.com/CunTianXing/go_app/go-micro/api_demo/meta/proto"
	"github.com/micro/go-api"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Example struct{}

type Foo struct{}

func (e *Example) Call(ctx context.Context, req *proto.CallRequest, rsp *proto.CallResponse) error {
	log.Print("Received Example.Call request")
	if len(req.Name) == 0 {
		return errors.BadRequest("go.micro.api.example", "no content")
	}
	rsp.Message = "got your request " + req.Name
	return nil
}

func (f *Foo) Bar(ctx context.Context, req *proto.EmptyRequest, rsp *proto.EmptyResponse) error {
	log.Print("Received Foo.Bar request")

	return nil
}

func main() {
	service := micro.NewService(
		micro.Name("go.micro.api.example"),
	)
	proto.RegisterExampleHandler(service.Server(), new(Example), api.WithEndpoint(&api.Endpoint{
		Name:    "Example.Call",
		Path:    []string{"/example"},
		Method:  []string{"GET", "POST"},
		Handler: api.Rpc,
	}))

	proto.RegisterFooHandler(service.Server(), new(Foo), api.WithEndpoint(&api.Endpoint{
		Name:    "Foo.Bar",
		Path:    []string{"/foo/bar"},
		Method:  []string{"POST"},
		Handler: api.Rpc,
	}))
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
