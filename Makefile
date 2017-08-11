all:
	make dependencies

dependencies:
	go get github.com/golang/protobuf/proto
	go get golang.org/x/net/context
	go get google.golang.org/grpc
