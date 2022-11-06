gen:
	protoc --go-grpc_out=require_unimplemented_servers=false:./ --go_out=./ ./proto/chat.proto

clean:
	rm chatserver/*.go

run:
	go run server.go