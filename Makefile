.PHONY: install build

PROTO_FILES := $(shell find proto -name '*.proto')
GENERATED_FILES := $(shell find pb -name '*.go')


$(GENERATED_FILES): $(PROTO_FILES)
	protoc \
		--go_out=pb \
		--go-grpc_out=pb \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=./proto $(PROTO_FILES)

install:
	go get github.com/grpc/grpc-go/cmd/protoc-gen-go-grpc
	go get github.com/golang/protobuf/protoc-gen-go

build: $(GENERATED_FILES)
	CGO_ENABLED=0 GOOS=linux go build

build_client: $(GENERATED_FILES)
	cd cmd/client && go build
