package main

import (
	"grpcChatServer/chatserver"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

func main() {
	//assign port
	Port := os.Getenv("PORT")
	if Port == "" {
		Port = "8080" //default port
	}

	//init listener
	listener, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen @ %v :: %v", Port, err)
	}
	log.Println("Listening @ :" + Port)

	//gRPC server instance
	grpc_server := grpc.NewServer()

	//register chatService
	cs := chatserver.ChatServer{}
	chatserver.RegisterServicesServer(grpc_server, &cs)

	//grpc listen and server
	err = grpc_server.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to start gRPC server :: %v", err)
	}
}
