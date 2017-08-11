package main

import (
	"flag"
	"fmt"
	"github.com/OBASHITechnology/RPCcancel_test/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

// Command Line Arguments
var (
	sleepTime        int    = 10
	oddTimeout       int    = 30
	evenTimeout      int    = 60
	outboundLocation string = "0.0.0.0:4454"
)

func main() {

	parseArguments()
	client := connectToRPCServer()

	// Generate 10 requests and handle any fraudulent ones
	for i := 0; i < 10; i++ {

		r := &protoTypes.Request{Message: "Hello", Id: int32(i)}
		success := SendRequest(r, client)
		if success {
			fmt.Println("SUCCESS")
		} else {
			fmt.Println("unsuccessful")
		}
		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}

/*****************
 * RPC Functions *
 *****************/

// SendRequest will send a request to the specified client, and return whether it succeeded
func SendRequest(r *protoTypes.Request, client protoTypes.FraudtestClient) bool {
	fmt.Println("Sending", r, "to", outboundLocation)

	// Set Timeout
	ctx := CreateContextWithTimeout(r)

	rpcReturn := make(chan *protoTypes.SuccessIndicator)

	go func() {
		success, err := client.TransferMessage(ctx, r)
		if err != nil {
			fmt.Println("Error occurred when sending message")
		}
		rpcReturn <- success
	}()

	select {
	case success := <-rpcReturn:
		fmt.Println("Call was successful!", success)
		return success.Success
	case <-ctx.Done():
		fmt.Println("Timeout occurred!")
		return false
	}

}

// CreateContextWithTimeout will take a request, and return a context with a certain timeout, based on it's ID.
// - If the ID of the request is even, then the timeout is set to the evenTimeout, specified via Command Line Arguments
// - Likewise for an odd ID
func CreateContextWithTimeout(r *protoTypes.Request) context.Context {
	ctx := context.Background()

	if IsEven(r) {

		fmt.Println("Request", r.Id, "timeout is", evenTimeout)
		ctx, _ = context.WithTimeout(ctx, time.Duration(evenTimeout)*time.Second)

	} else {

		fmt.Println("Request", r.Id, "timeout is", oddTimeout)
		ctx, _ = context.WithTimeout(ctx, time.Duration(oddTimeout)*time.Second)

	}

	return ctx
}

/*********************
 * Utility Functions *
 *********************/

// IsEven will return true if a requests ID is even, otherwise false
func IsEven(r *protoTypes.Request) bool {
	return r.Id%2 == 0
}

func parseArguments() {
	flag.StringVar(&outboundLocation, "outboundLocation", outboundLocation, "The IP address (and port) to forward messages to")
	flag.IntVar(&sleepTime, "sleepTime", sleepTime, "The number of seconds to sleep between sending a request and checking if it is fraudulent")
	flag.IntVar(&evenTimeout, "evenTimeout", evenTimeout, "The number of seconds given to even requests to complete")
	flag.IntVar(&oddTimeout, "oddTimeout", oddTimeout, "The number of seconds given to odd requests to complete")

	flag.Parse()
}

func connectToRPCServer() protoTypes.FraudtestClient {
	conn, err := grpc.Dial(outboundLocation, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return protoTypes.NewFraudtestClient(conn)
}
