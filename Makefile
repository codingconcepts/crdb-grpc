.PHONY: protos

protos:
	protoc --proto_path=protos protos/*.proto --go_out=. --go-grpc_out=.

server:
	go run main.go

rpc_list:
	grpcurl -plaintext localhost:8080 list

rpc_hello:
	grpcurl \
		-plaintext \
		-d '{"name": "Rob"}' \
			localhost:8080 hello.HelloService/SayHello