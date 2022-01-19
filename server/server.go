package server

import (
	"context"
	pb "grpc-intro/usermgmt"
	"log"
	"math/rand"
	"net"

	"google.golang.org/grpc"
)

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (u *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("\nRecieved %s\n\n", in.GetName())
	userID := rand.Intn(1000)

	return &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   int32(userID),
	}, nil
}

// StartServer ...
func StartServer(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, &UserManagementServer{})
	log.Printf("Server listening at: %s", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}
