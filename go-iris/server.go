package main

import (
    "github.com/kataras/iris"
    "github.com/kataras/iris/middleware/logger"
    "github.com/kataras/iris/middleware/recover"
)

func main() {
    app := iris.New()
    app.Use(recover.New())
    app.Use(logger.New())

    app.Handle("GET","/",func(ctx iris.Context) {
        ctx.HTML("<h1>Wellcome</h1>")
    })

    app.Get("/ping",func(ctx iris.Context) {
        ctx.WriteString("pong")
    })

    app.Get("/hello", func(ctx iris.Context) {
        ctx.JSON(iris.Map{"message": "Hello Iris!"})
    })

    app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
