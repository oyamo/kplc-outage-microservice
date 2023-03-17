# Set the directory containing your protobuf files
PROTO_DIR := ./

# Find all .proto files in the directory and its subdirectories
PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

# Compile all the protobuf files
all: $(patsubst %.proto,%.pb.go,$(PROTO_FILES)) $(patsubst %.proto,%.pb.py,$(PROTO_FILES))

# Compile a single protobuf file for Go
%.pb.go: %.proto
	protoc -I=$(PROTO_DIR) --go_out=$(PROTO_DIR) --go_opt=paths=source_relative --go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative $<

# Compile a single protobuf file for Python
%.pb.py: %.proto
	python -m grpc_tools.protoc -I $(PROTO_DIR) --python_out=$(PROTO_DIR) --grpc_python_out=$(PROTO_DIR) $<

# Clean up generated files
clean:
	find $(PROTO_DIR) -name '*.pb.go' -delete
	find $(PROTO_DIR) -name '*_pb2.py' -delete
	find $(PROTO_DIR) -name '*_pb2_grpc.py' -delete
