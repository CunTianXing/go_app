package main

import (
	"fmt"
	"os"

	proto "github.com/CunTianXing/go_app/go-micro/service/proto"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"golang.org/x/net/context"
)

type Greeter struct{}

func (g *Greeter) Hello(ctx context.Context, req *proto.HelloRequest, rsp *proto.HelloResponse) error {
	rsp.Message = "hello " + req.Name
	return nil
}

func runClient(service micro.Service) {
	greeter := proto.NewGreeterClient("greeter", service.Client())
	rsp, err := greeter.Hello(context.TODO(), &proto.HelloRequest{Name: "gary"})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rsp.Message)
}

func main() {
	service := micro.NewService(
		micro.Name("greeter"),
		micro.Version("latest"),
		micro.Metadata(map[string]string{
			"type": "helloworld",
		}),
		micro.Flags(cli.BoolFlag{
			Name:  "run_client",
			Usage: "Launch the client",
		}),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			if c.Bool("run_client") {
				runClient(service)
				os.Exit(0)
			}
		}),
	)
	proto.RegisterGreeterHandler(service.Server(), new(Greeter))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
