
gen-proto-go:
	rm -f ../backend/core/proto/*.pb.go
	protoc -I=. -I=vendor -I=${GOPATH}/src *.proto --go_out=.. --go-grpc_out=..


#PROTOC_GEN_GRPC_PYTHON=$(shell which grpc_tools.protoc)
#.PHONY: gen-proto-python
#gen-proto-python:
#	rm -f ./python/embedder/*pb2*.py*
#	python -I=. --proto_path=. --python_out=./python --pyi_out=./python --grpc_python_out=./python semantic_data.proto
