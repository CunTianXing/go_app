package main

import (
	"log"
	"net"

	"github.com/CunTianXing/go_app/docker-micro/proto/users"
	"github.com/CunTianXing/go_app/docker-micro/shared"
	"github.com/CunTianXing/go_app/docker-micro/users/server"
	"google.golang.org/grpc"
)

func main() {
	listener, _ := net.Listen("tcp", ":8082")
	log.Print("[main] service started")

	shared.Init()

	srv := grpc.NewServer()
	users.RegisterUsersServer(srv, &server.Server{})
	srv.Serve(listener)
}
