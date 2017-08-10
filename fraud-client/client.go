package main

import (
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

func main() {
	for i := 0; i < 10; i++ {
		// send request
		r := &Request{Data: "Hello", ID: i}
		// sleep 10 seconds
		time.Sleep(time.Second * 10)
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
