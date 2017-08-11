package main

import (
	"errors"
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

var oddTimeout = 30

var evenTimeout = 60

var ip string

func main() {

	parseArguments()
	client := connectToRPCServer()

	// Generate 10 requests and handle any fraudulent ones
	for i := 0; i < 10; i++ {

		r := &protoTypes.Request{Message: "Hello", Id: int32(i)}
		SendRequest(r, client)

		time.Sleep(time.Second * time.Duration(sleepTime))
	}
}

/*****************
 * RPC Functions *
 *****************/

// IsEven will return true if a requests ID is even, otherwise false
func IsEven(r *protoTypes.Request) bool {
	return r.Id%2 == 0
}

func SendRequest(r *protoTypes.Request, client protoTypes.FraudtestClient) context.Context {
	fmt.Println("Sending", r, "to", ip)

	// Set Timeout
	ctx := context.Background()
	if IsEven(r) {
		fmt.Println("Request", r.Id, "timeout is", evenTimeout)
		ctx.Deadline()
	} else {
		fmt.Println("Request", r.Id, "timeout is", oddTimeout)
		ctx.Deadline()
	}

	success := new(protoTypes.SuccessIndicator)
	var err error
	rpcReturn := make(chan *protoTypes.SuccessIndicator)

	go func() {
		success, err = client.TransferMessage(ctx, r)
		if err != nil {
			panic(err)
		}
	}()

	if !success.Success {
		panic("Transfer Returned False for Request")
	}

	select {
	case success := <-rpcReturn:
		fmt.Println("Call was successful!", success)
	case <-ctx.Done():
		panic(errors.New("Timeout occurred!"))
	}

	return ctx

}

/*********************
 * Utility Functions *
 *********************/

func parseArguments() {
	ipFlag := flag.String("ip", "0.0.0.0:4454", "The IP address (and port) to forward messages to")
	sleepFlag := flag.Int("sleep", 10, "The number of seconds to sleep between sending a request and checking if it is fraudulent")
	evenFlag := flag.Int("even-timeout", 60, "The number of seconds given to even requests to complete")
	oddFlag := flag.Int("odd-timeout", 30, "The number of seconds given to odd requests to complete")

	flag.Parse()

	ip = *ipFlag
	sleepTime = *sleepFlag
	evenTimeout = *evenFlag
	oddTimeout = *oddFlag
}

func connectToRPCServer() protoTypes.FraudtestClient {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return protoTypes.NewFraudtestClient(conn)
}
