package server

import (
	"log"

	"github.com/CunTianXing/go_app/docker-micro/proto/list"
	"github.com/CunTianXing/go_app/docker-micro/proto/users"
	"github.com/CunTianXing/go_app/docker-micro/shared"
	"golang.org/x/net/context"
)

type Server struct{}

func (s *Server) CreateUser(ctx context.Context, in *users.CreateUserRequest) (*users.CreateUserResponse, error) {
	userID, err := shared.CreateUser(in.Email)

	response := new(users.CreateUserResponse)
	if err == nil {
		log.Printf("[user.Create] New user ID: %s", userID)

		createInitialItem(userID)

		// TODO: send email to user when it's created.

		response.Message = "User created successfully"
		response.Id = userID
		response.Code = 200
	} else {
		response.Message = err.Error()
		response.Code = 500
	}

	return response, err
}

// Create initial item in todo list
func createInitialItem(userID string) {
	_, err := shared.ListClient.CreateItem(context.Background(), &list.CreateItemRequest{
		Message: "Welcome to Workshop!",
		UserId:  userID,
	})
	if err != nil {
		log.Printf("[user.Create] Cannot create item: %v", err)
	}
}
