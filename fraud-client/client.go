package main

import (
	"flag"
	"fmt"
	"time"
)

type Request struct {
	Data string
	ID   int
}

type SuccessIndicator struct {
	Ok bool
}

const SleepTime = 1

var ip string

func main() {

	parseArguments()

	for i := 0; i < 10; i++ {
		// Generate request
		r := &Request{Data: "Hello", ID: i}

		// send request
		go SendRequest(r)

		// sleep 10 seconds
		time.Sleep(time.Second * SleepTime)
		// if ID is even then FRAUD
		if IsFraudulent(r) {
			fmt.Println("Request", r.ID, "is FRAUDULENT")
		}
	}
}

// IsFraudulent will return true if a requests ID is even, otherwise false
func IsFraudulent(r *Request) bool {
	return r.ID%2 == 0
}

func SendRequest(r *Request) {
	fmt.Println("Sending", r, "to", ip)
}

func parseArguments() {
	ipFlag := flag.String("ip", "192.168.100.31:4455", "The IP address (and port) to forward messages to")

	flag.Parse()

	ip = *ipFlag
}
