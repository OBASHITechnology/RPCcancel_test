package main

import (
	"flag"
	"fmt"
	"github.com/OBASHITechnology/RPCcancel_test/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

/**************************
 * Command Line Arguments *
 **************************/

var sleepTime = 1

var ip string

func main() {

	parseArguments()
	client := connectToRPCServer()

	// Generate 10 requests and handle any fraudulent ones
	for i := 0; i < 10; i++ {

		r := &protoTypes.Request{Message: "Hello", Id: int32(i)}
		ctx := SendRequest(r, client)

		time.Sleep(time.Second * time.Duration(sleepTime))

		if IsFraudulent(r) {
			HandleFraudulent(r, ctx)
		}
	}
}

/*****************
 * RPC Functions *
 *****************/

// IsFraudulent will return true if a requests ID is even, otherwise false
func IsFraudulent(r *protoTypes.Request) bool {
	return r.Id%2 == 0
}

func HandleFraudulent(r *protoTypes.Request, ctx context.Context) {
	fmt.Println("Request", r.Id, "is FRAUDULENT")
	ctx.Done()
}

func SendRequest(r *protoTypes.Request, client protoTypes.FraudtestClient) context.Context {
	fmt.Println("Sending", r, "to", ip)

	ctx := context.Background()
	success, err := client.TransferMessage(ctx, r)
	if err != nil {
		panic(err)
	}

	if !success.Success {
		panic("Transfer Returned False for Request")
	}
	return ctx

}

/*********************
 * Utility Functions *
 *********************/

func parseArguments() {
	ipFlag := flag.String("ip", "0.0.0.0:4455", "The IP address (and port) to forward messages to")
	sleepFlag := flag.Int("sleep", 10, "The number of seconds to sleep between sending a request and checking if it is fraudulent")

	flag.Parse()

	ip = *ipFlag
	sleepTime = *sleepFlag
}

func connectToRPCServer() protoTypes.FraudtestClient {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return protoTypes.NewFraudtestClient(conn)
}
