package server

import (
	"context"
	pb "grpc-intro/usermgmt"
	"log"
	"net"

	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

type UserManagementServer struct {
	conn *pgx.Conn
	pb.UnimplementedUserManagementServer
}

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

func (u *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("\nRecieved %s\n\n", in.GetName())
	// createSQL := `
	// 	CREATE TABLE IF NOT EXISTS users (
	// 		id SERIAL PRIMARY KEY,
	// 		name text,
	// 		age int
	// 	);
	// `
	// _, err := u.conn.Exec(ctx, createSQL)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Table creation failed: %w\n", err)
	// 	os.Exit(1)
	// }

	createdUser := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
	}
	tx, err := u.conn.Begin(ctx)
	if err != nil {
		log.Fatal("connection begin failed: %w\n", err)
	}
	// defer tx.Conn().Close(ctx)
	_, err = tx.Exec(ctx, "INSERT INTO users (name, age) VALUES($1,$2)", createdUser.Name, createdUser.Age)
	if err != nil {
		log.Fatal("Cannot insert new user: %w", err)
	}
	tx.Commit(ctx)
	return createdUser, nil
}

func (u *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var usersList *pb.UserList = &pb.UserList{}
	rows, err := u.conn.Query(ctx, "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}

		usersList.Users = append(usersList.Users, &user)
	}
	return usersList, nil
}

// StartServer ...
func StartServer(port string) {
	dbaseURL := "postgres://postgres:password@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), dbaseURL)
	if err != nil {
		log.Fatal("Unable to connect to db")
	}
	// defer conn.Close(context.Background())
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("failed to listen: %w", err)
	}
	usr := NewUserManagementServer()
	usr.conn = conn

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, usr)
	log.Printf("Server listening at: %s", lis.Addr().String())
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %w", err)
	}
}
