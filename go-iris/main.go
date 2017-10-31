package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())
	app.Controller("/",new(ExampleController))
	app.Run(iris.Addr(":8080"))
}

type ExampleController struct {
	mvc.C
}

func (c *ExampleController) Get() mvc.Result {
	return mvc.Response{
		ContentType:"text/html",
		Text:"<h1>Wellcome</h1>",
	}
}

func (c *ExampleController) GetPing() string {
	return "pong"
}

func (c *ExampleController) GetHello() interface{} {
	return map[string]string{"message":"hello Iris!"}
}
