package main

import (
	"bufio"
	"context"
	"fmt"
	"grpcChatServer/chatserver"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Enter Server IP:Port ::: ")
	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Failed to read from console :: %v", err)
	}

	serverID = strings.Trim(serverID, "\r\n")
	log.Println("Connecting : " + serverID)

	//connect to grpc server --> dial grpc server to initiate connection using grpc.Dial() method
	conn, err := grpc.Dial(serverID, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to cconnect to gRPC server :: %v", err)
	}

	defer conn.Close()

	//call chatService to create a stream
	client := chatserver.NewServicesClient(conn)
	stream, err := client.ChatService(context.Background())
	//ChatService RPC method to create a stream passing a context with no deadline
	if err != nil {
		log.Fatalf("Failed to call ChatService :: %v", err)
	}

	//running 2 go routines at client side for sending and receiving messages

	//implement communication with gRPC server
	ch := clienthandle{stream: stream}
	ch.clientConfig()

	go ch.sendMessage()
	go ch.receiveMessage()

	//blocker
	bl := make(chan bool)
	<-bl
}

// clientHandle : stream and client's name
type clienthandle struct {
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func (ch *clienthandle) clientConfig() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Your Name : ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read from console :: %v", err)
	}
	ch.clientName = strings.Trim(name, "\r\n")
}

// send message
func (ch *clienthandle) sendMessage() {

	for {
		reader := bufio.NewReader(os.Stdin)

		clientMessage, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read from console :: %v", err)
		}
		clientMessage = strings.Trim(clientMessage, "\r\n")

		clientMessageBox := &chatserver.FromClient{
			Name: ch.clientName,
			Body: clientMessage,
		}

		err = ch.stream.Send(clientMessageBox)
		if err != nil {
			log.Printf("Error while sending message to the server :: %v", err)
		}
	}
}

// receive message
func (ch *clienthandle) receiveMessage() {
	for {
		mssg, err := ch.stream.Recv()
		if err != nil {
			log.Printf("Error in receiving the message :: %v", err)
		}
		//print the message to the console
		fmt.Printf("\n %s : %s \n", mssg.Name, mssg.Body)
	}
}
