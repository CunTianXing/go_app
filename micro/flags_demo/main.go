package main

import (
	"fmt"
	"os"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
)

func main() {
	service := micro.NewService(
		micro.Flags(
			cli.StringFlag{
				Name:  "string_flag",
				Usage: "this is a string flag",
			},
			cli.IntFlag{
				Name:  "int_flag",
				Usage: "this is an int flag",
			},
			cli.BoolFlag{
				Name:  "bool_flag",
				Usage: "this is a bool flag",
			},
		),
	)

	service.Init(
		micro.Action(func(c *cli.Context) {
			fmt.Printf("the string flag is: %s\n", c.String("string_flag"))
			fmt.Printf("the int flag is: %d\n", c.Int("int_flag"))
			fmt.Printf("the bool flag is: %t\n", c.Bool("bool_flag"))
			os.Exit(0)
		}),
	)
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}

//go run main.go --string_flag="a string" --int_flag=10 --bool_flag=true
