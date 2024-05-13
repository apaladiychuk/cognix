
gen-proto-go:
	rm -f ../backend/core/proto/*.pb.go
	protoc -I=. -I=vendor -I=${GOPATH}/src *.proto --go_out=.. --go-grpc_out=..

gen-proto-python:
	rm -f ai/core/proto/*.pb.py
	protoc -I=. *.proto --python_out=. --python-grpc_out=.