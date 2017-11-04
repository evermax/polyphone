.PHONY: build

build: protos/auth.pb.go
	go build .

protos/auth.pb.go:
	protoc --go_out="plugins=grpc:." protos/auth.proto

test:
	go test ./...