package main

import (
	"grpc-intro/server"
)

func main() {
	const port = ":50521"
	server.StartServer(port)
}
