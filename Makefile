# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /bin/bash

GO_TOOLS = 	google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
			google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
			github.com/gogo/protobuf/protoc-gen-gofast@v1.3.1 \
			github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.14.7 \
			github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.14.7 \
			github.com/envoyproxy/protoc-gen-validate@v0.6.1 \
			github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.3.2 \
			github.com/bold-commerce/protoc-gen-struct-transformer@v1.0.7 \
			github.com/google/wire/cmd/wire@v0.4.0 \
			golang.org/x/tools/cmd/goimports@v0.1.11 \
			github.com/bufbuild/buf/cmd/buf@v1.4.0 \

install-go-tools:
	@echo $(GO_TOOLS) | xargs -r -n1 go install

update:
	go mod tidy
	go mod vendor

go-doc:
	godoc -http=:6060