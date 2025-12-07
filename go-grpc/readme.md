# Prerequisites

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

# Run gRPC Code Generation

Change working directory to `routeguide`

```shell
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       --plugin=protoc-gen-go=$HOME/go/bin/protoc-gen-go \
       --plugin=protoc-gen-go-grpc=$HOME/go/bin/protoc-gen-go-grpc \
       route_guide.proto
```

# Commands

## Refresh dependencies

```shell
go mod tidy
```

## Run Server / Client

```shell
cd server
go run .

cd client
go run .
```