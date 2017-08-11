package main

import (
	"flag"
	protoTypes "github.com/OBASHITechnology/RPCcancel_test/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
	"time"
	"errors"
)

/*

	Setup

*/

// Variables we'll be configuring at argument parsing
var (
	inboundLocation  string = "0.0.0.0:4455"
	outboundLocation string = "0.0.0.0:4455"
	sleepLength int = 1 // Time to sleep in seconds
)

// What we'll be serving on
type MessageAcceptor bool

/*

	Functions and methods

*/
func (MessageAcceptor) TransferMessage(ctx context.Context, request *protoTypes.Request) (*protoTypes.SuccessIndicator, error) {
	// Sleep a configurable length of time before processing the recieved message
	time.Sleep(time.Duration(sleepLength) * time.Second)

	conn, err := grpc.Dial(outboundLocation, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	client := protoTypes.NewFraudtestClient(conn)

	// Make an asynchronous RPC call
	rpcReturn := make(chan *protoTypes.SuccessIndicator)
	go func() {
		success, err := client.TransferMessage(ctx, request)
		if err != nil {
			log.Fatal(err.Error())
		}
		rpcReturn <- success
	}()

	// Check to see whether our async call returns before our context times out
	select {
	case success := <- rpcReturn:
		return success, nil
	case <- ctx.Done():
		return &protoTypes.SuccessIndicator{false}, errors.New("Timeout occurred!")
	}
}

func main() {
	// Configure from the commandline
	parseArguments()

	// Variables we'll use
	var (
		lis        net.Listener
		err        error
		server     *grpc.Server = grpc.NewServer()
		messageAcc MessageAcceptor
		wg         sync.WaitGroup
	)

	wg.Add(1)

	// Get a Listener on the required location, set up by parsing the arguments earlier
	if lis, err = net.Listen("tcp", inboundLocation); err != nil {
		log.Fatalf(err.Error())
	}

	// Register the server with the gRPC server, and for reflection
	protoTypes.RegisterFraudtestServer(server, protoTypes.FraudtestServer(messageAcc))
	reflection.Register(server)

	// Perform the serving!
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf(err.Error())
			wg.Add(-1)
		}
	}()

	wg.Wait()
}

func parseArguments() {
	flag.StringVar(&inboundLocation, "inLocation", inboundLocation, "The location where messages will be listened for on. Format IP:PORT")
	flag.StringVar(&outboundLocation, "outLocation", outboundLocation, "The location which messages will be forwarded to.")
	flag.IntVar(&sleepLength, "sleep", sleepLength, "Amount of time in seconds that the process will sleep for before processing a recieved message")
	flag.Parse()
}
