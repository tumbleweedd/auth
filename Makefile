.PHONY: migrate

migrate:
	go run cmd/migrator/main.go -storage-path "postgres:postgres@localhost:5433/auth_service?sslmode=disable" -migrations-path migrations

# Path configuration
PROTO_DIR := ./proto
GO_OUT_DIR := ./gen

# Protobuf files
TOKEN_PROTO_FILES := $(wildcard $(PROTO_DIR)/token/*.proto)
USER_PROTO_FILES := $(wildcard $(PROTO_DIR)/user/*.proto)

# Protoc commands
PROTOC_GEN_GO := --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative
PROTOC_GEN_GO_GRPC := --go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative

# Ensure that all dependencies are installed and paths are set
.PHONY: deps
deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	export PATH="$PATH:$(go env GOPATH)/bin"

# Generate Go files from proto files
.PHONY: gen-proto
gen-proto:
	mkdir -p $(GO_OUT_DIR)
	protoc -I $(PROTO_DIR) $(TOKEN_PROTO_FILES) $(USER_PROTO_FILES) $(PROTOC_GEN_GO) $(PROTOC_GEN_GO_GRPC)

# Clean generated files
.PHONY: clean
clean:
	rm -rf $(GO_OUT_DIR)

# Default target
all: gen-proto
