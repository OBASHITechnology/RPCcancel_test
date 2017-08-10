package main

import (
	protoTypes "github.com/OBASHITechnology/RPCcancel_test/protobuf"
	"flag"
	"net"
	"log"
	"google.golang.org/grpc"
	"context"
	"google.golang.org/grpc/reflection"
)

/*

	Setup

 */

// Variables we need to be global
var (
	acceptedMessages map[int32]string
)

// Variables we'll be configuring at argument parsing
var (
	inboundLocation string = "0.0.0.0:4455"
	inspectionLocation string = "0.0.0.0:8181"
)

// What we'll be serving on
type MessageAcceptor bool


/*

	Functions and methods

 */
func (MessageAcceptor) TransferMessage(ctx context.Context, request *protoTypes.Request) (*protoTypes.SuccessIndicator, error) {
	acceptedMessages[request.Id] = request.Message
	return &protoTypes.SuccessIndicator{Success:true}, nil
}

func main() {
	// Configure from the commandline
	parseArguments()

	// Variables we'll use
	var (
		lis net.Listener
		err error
		server *grpc.Server = grpc.NewServer()
		messageAcc MessageAcceptor
	)

	// Get a Listener on the required location, set up by parsing the arguments earlier
	if lis, err = net.Listen("tcp", inboundLocation); err != nil {
		log.Fatalf(err.Error())
	}

	// Register the server with the gRPC server, and for reflection
	protoTypes.RegisterFraudtestServer(server, protoTypes.FraudtestServer(messageAcc))
	reflection.Register(server)

	// Perform the serving!
	if err := server.Serve(lis); err != nil {
		log.Fatalf(err.Error())
	}
}

func parseArguments() {
	flag.StringVar(&inboundLocation, "inLocation", inboundLocation, "The location where messages will be listened for on. Format IP:PORT")
	flag.StringVar(&inspectionLocation, "readLocation", inspectionLocation, "The location where messages will be served out again over HTTP, so we can see them on the web")
	flag.Parse()
}
