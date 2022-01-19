package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "grpc-intro/usermgmt"

	"google.golang.org/grpc"
)

const address = "localhost:50521"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	newUsers := map[string]int32{
		"Alice": 43,
		"Bob":   30,
	}

	for name, age := range newUsers {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatal("an error occurred ", err)
		}

		fmt.Printf("\nUser Details:\nName: %s\nAge: %d\nID: %d\n", r.GetName(), r.GetAge(), r.GetId())
	}
}
