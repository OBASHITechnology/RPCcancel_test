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
		// sleep 10 seconds
		time.Sleep(time.Second * 10)
		// if ID is even then FRAUD
		if i%2 == 0 {
			fmt.Println("Request", i, "is FRAUDULENT")
		}
	}
}

// CheckFraud will return true if a requests ID is even, otherwise false
func CheckFraud(r Request) bool {
	return r.ID%2 == 0
}
