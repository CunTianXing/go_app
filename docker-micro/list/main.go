package main

import (
	"log"
	"net"

	"github.com/CunTianXing/go_app/docker-micro/list/server"
	"github.com/CunTianXing/go_app/docker-micro/proto/list"
	"github.com/CunTianXing/go_app/docker-micro/shared"
	"google.golang.org/grpc"
)

func main() {
	listener, _ := net.Listen("tcp", ":8081")
	log.Print("[main] service started")

	shared.Init()

	srv := grpc.NewServer()
	list.RegisterListServer(srv, &server.Server{})
	srv.Serve(listener)
}
