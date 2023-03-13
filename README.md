# crdb-grpc
A barebones example of accessing CockroachDB via gRPC

### Dependencies

* [Go](https://go.dev)
* [Protobuf](https://protobuf.dev)
* [Protobuf Gen Go](https://grpc.io/docs/languages/go/quickstart)
* [grpcurl](https://github.com/fullstorydev/grpcurl)

### Running Locally

Start a cluster
```
$ cockroach start-single-node \
    --listen-addr=localhost:26257 \
    --http-addr=localhost:8080 \
    --insecure
```

Create a table
```
$ cockroach sql --insecure < create.sql
```

Generate the proto service
```
$ protoc --proto_path=protos protos/*.proto --go_out=. --go-grpc_out=.
```

Run the server
```
$ go run main.go
```

Make a request to the gRPC reflection endpoint
```
$ grpcurl -plaintext localhost:8080 list TodoService
```

### Requests

Create a todo
```
$ grpcurl \
	-plaintext \
	-d '{"title": "todo d"}' \
		localhost:8080 TodoService.CreateTodo
```

Fetch all todos
```
$ grpcurl \
    -plaintext \
    localhost:8080 TodoService.GetTodos
```

Fetch one todo
```
$ grpcurl \
	-plaintext \
	-d '{"id": "15c494cd-1c87-45ae-b725-cbf436eda43c"}' \
		localhost:8080 TodoService.GetTodo
```

Delete a todo
```
$ grpcurl \
	-plaintext \
	-d '{"id": "15c494cd-1c87-45ae-b725-cbf436eda43c"}' \
		localhost:8080 TodoService.DeleteTodo
```