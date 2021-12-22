package services

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"github.com/abdullohsattorov/api-gateway/config"
	pb "github.com/abdullohsattorov/api-gateway/genproto"
)

type IServiceManager interface {
	TodoService() pb.TodoServiceClient
}

type serviceManager struct {
	todoService pb.TodoServiceClient
}

func (s *serviceManager) TodoService() pb.TodoServiceClient {
	return s.todoService
}

func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connTodo, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.TodoServiceHost, conf.TodoServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	serviceManager := &serviceManager{
		todoService: pb.NewTodoServiceClient(connTodo),
	}

	return serviceManager, nil
}
