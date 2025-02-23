package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Go's select construct can be leveraged to send the same request to multiple servers
// and process the response from whichever server responds first.
// This approach is very useful in distributed systems with multiple replicated servers,
// as it allows us to serve the user the fastest response available.
// Additionally, we can use a timeout pattern to control the total execution time of the request.

func dbQuery(dbID int, out chan<- string) {
	// use random delay for simulating processing requests
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	out <- fmt.Sprintf("DB server %d: Query result", dbID)
}

func main() {

	response := make(chan string) // response channel
	defer close(response)

	start := time.Now()

	// fire request to 3 DB servers
	for i := 0; i < 3; i++ {
		go dbQuery(i+1, response)
	}

	// process first response that arrives or timeout within 70ms
	timeout := time.After(70 * time.Millisecond)
	select {
	case result := <-response:
		fmt.Println(result)
	case <-timeout:
		fmt.Println("All requests timed out!")
	}

	fmt.Println("Time elapsed:", time.Since(start))
}
