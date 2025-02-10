package main

import (
	"fmt"
	"math/rand"
	"time"
)

func worker(id int, requests <-chan int, response chan<- string) {

	// keep accepting requests as long as the request channel is open
	for in := range requests {
		// simulate the request work by random delay
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		response <- fmt.Sprintf("Worker %d has finished a request %d", id, in)
	}
}

func loadBalancer(numWorkers int, requests <-chan int) <-chan string {
	out := make(chan string)
	for i := 0; i < numWorkers; i++ {
		go worker(i, requests, out)
	}
	return out
}

func main() {

	const numWorkers int = 4
	const numRequests int = 10

	requests := make(chan int, 100) // buffered requests channel
	responseCh := loadBalancer(numWorkers, requests)

	// fire the requests (in the background)
	go func(ch chan<- int) {
		for i := 0; i < numRequests; i++ {
			requests <- i
		}
		close(requests) // close the request channel once done firing them all
	}(requests)

	// read the responses
	for i := 0; i < numRequests; i++ {
		fmt.Println(<-responseCh)
	}
}
