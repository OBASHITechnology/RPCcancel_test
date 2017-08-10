package main

import (
	"flag"
	"fmt"
	"github.com/OBASHITechnology/RPCcancel_test/protobuf"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"time"
)

const SleepTime = 1

var ip string

func main() {

	parseArguments()
	client := connectToRPCServer()

	// Generate 10 requests and handle any fraudulent ones
	for i := 0; i < 10; i++ {

		r := &protoTypes.Request{Message: "Hello", Id: int32(i)}
		go SendRequest(r, client)

		time.Sleep(time.Second * SleepTime)

		if IsFraudulent(r) {
			HandleFraudulent(r)
		}
	}
}

// IsFraudulent will return true if a requests ID is even, otherwise false
func IsFraudulent(r *protoTypes.Request) bool {
	return r.Id%2 == 0
}

func HandleFraudulent(r *protoTypes.Request) {
	fmt.Println("Request", r.Id, "is FRAUDULENT")
}

func SendRequest(r *protoTypes.Request, client protoTypes.FraudtestClient) {
	fmt.Println("Sending", r, "to", ip)

	success, err := client.TransferMessage(context.Background(), r)
	if err != nil {
		panic(err)
	}

	if !success.Success {
		panic("Transfer Returned False for Request")
	}

}

func parseArguments() {
	ipFlag := flag.String("ip", "192.168.100.31:4455", "The IP address (and port) to forward messages to")

	flag.Parse()

	ip = *ipFlag
}

func connectToRPCServer() protoTypes.FraudtestClient {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	return protoTypes.NewFraudtestClient(conn)
}