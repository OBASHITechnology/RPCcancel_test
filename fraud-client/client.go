package main

import (
	"flag"
	"fmt"
	"github.com/OBASHITechnology/RPCcancel_test/protobuf"
	"time"
)

const SleepTime = 1

var ip string

func main() {

	parseArguments()

	for i := 0; i < 10; i++ {
		// Generate request
		r := &protoTypes.Request{Message: "Hello", Id: int32(i)}

		// send request
		go SendRequest(r)

		// sleep for a while
		time.Sleep(time.Second * SleepTime)

		// if ID is even then FRAUD
		if IsFraudulent(r) {
			fmt.Println("Request", r.Id, "is FRAUDULENT")
		}
	}
}

// IsFraudulent will return true if a requests ID is even, otherwise false
func IsFraudulent(r *protoTypes.Request) bool {
	return r.Id%2 == 0
}

func SendRequest(r *protoTypes.Request) {
	fmt.Println("Sending", r, "to", ip)
}

func parseArguments() {
	ipFlag := flag.String("ip", "192.168.100.31:4455", "The IP address (and port) to forward messages to")

	flag.Parse()

	ip = *ipFlag
}
