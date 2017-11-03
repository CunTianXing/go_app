package main

import (
	"fmt"

	proto "github.com/CunTianXing/go_app/go-micro/service/proto"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
)

type logWrapper struct {
	client.Client
}

// client middleware
func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	//rsp = "test client middleware"
	fmt.Println(req.ContentType())
	fmt.Printf("[wrapper] client request service: %s method: %s\n", req.Service(), req.Method())
	return l.Client.Call(ctx, req, rsp)
}

func logWrap(c client.Client) client.Client {
	return &logWrapper{c}
}

func main() {
	service := micro.NewService(
		micro.Name("greeter.client"),
		micro.WrapClient(logWrap),
	)
	service.Init()
	greeter := proto.NewGreeterClient("greeter", service.Client())
	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "gary"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp.Message)
}
