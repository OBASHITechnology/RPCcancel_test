package main

import (
	"flag"
	protoTypes "github.com/OBASHITechnology/RPCcancel_test/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
)

/*

	Setup

*/

// Variables we need to be global
var (
	acceptedMessages map[int32]string = make(map[int32]string)
)

// Variables we'll be configuring at argument parsing
var (
	inboundLocation    string = "0.0.0.0:4455"
	inspectionLocation string = "0.0.0.0:8181"
)

// What we'll be serving on
type MessageAcceptor bool

/*

	Functions and methods

*/
func (MessageAcceptor) TransferMessage(ctx context.Context, request *protoTypes.Request) (*protoTypes.SuccessIndicator, error) {
	acceptedMessages[request.Id] = request.Message
	return &protoTypes.SuccessIndicator{Success: true}, nil
}

func serveMessages(res http.ResponseWriter, _ *http.Request) {
	for id, message := range acceptedMessages {
		res.Write([]byte(strconv.FormatInt(int64(id), 10)))
		res.Write([]byte(": "))
		res.Write([]byte(message))
		res.Write([]byte("\n"))
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

	wg.Add(2)

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

	// Now that we've got a goroutine running the gRPC server, let's also serve HTTP.
	go func() {
		http.HandleFunc("/", serveMessages)
		if err := http.ListenAndServe(inspectionLocation, nil); err != nil {
			log.Fatalf(err.Error())
			wg.Add(-1)
		}
	}()

	wg.Wait()
}

func parseArguments() {
	flag.StringVar(&inboundLocation, "inLocation", inboundLocation, "The location where messages will be listened for on. Format IP:PORT")
	flag.StringVar(&inspectionLocation, "readLocation", inspectionLocation, "The location where messages will be served out again over HTTP, so we can see them on the web")
	flag.Parse()
}
