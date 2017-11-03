package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	proto "github.com/CunTianXing/go_app/go-micro/api_demo/default/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

type Example struct{}

type Foo struct{}

func (e *Example) Call(ctx context.Context, req *proto.Request, resp *proto.Response) error {
	log.Print("Received Example.Call request")
	//fmt.Println(req)
	fmt.Printf("req: %+v\n", req)
	name, ok := req.Get["name"]
	if !ok || len(name.Values) == 0 {
		return errors.BadRequest("go.micro.api.example", "no content")
	}
	resp.StatusCode = 200
	b, _ := json.Marshal(map[string]string{
		"message": "got your request " + strings.Join(name.Values, " "),
	})
	resp.Body = string(b)
	return nil
}

func (f *Foo) Bar(ctx context.Context, req *proto.Request, resp *proto.Response) error {
	log.Print("Received Foo.Bar request")
	if req.Method != "POST" {
		return errors.BadRequest("go.micro.api.example", "require post")
	}
	ct, ok := req.Header["Content-Type"]
	if !ok || len(ct.Values) == 0 {
		return errors.BadRequest("go.micro.api.example", "need content-type")
	}

	if ct.Values[0] != "application/json" {
		return errors.BadRequest("go.micro.api.example", "expect application/json")
	}
	var body map[string]interface{}
	json.Unmarshal([]byte(req.Body), &body)
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
