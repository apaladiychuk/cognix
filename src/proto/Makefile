
gen-proto-go:
	rm -f ../golang/core/proto/*.pb.go
	protoc -I=. -I=vendor -I=${GOPATH}/src *.proto --go_out=.. --go-grpc_out=..

# gen-proto-python:
# do not use thi commad it is not correct for python
# 	rm -f ./python/embedder/*pb2*.py*
#	python3 -m grpc_tools.protoc -I . --python_out=. --pyi_out=. --grpc_python_out=. embedd_service.proto