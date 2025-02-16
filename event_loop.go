package main

import (
	"fmt"
	"math/rand"
	"time"
)

// event handler function signature
type EventHandler func(int)

type Event struct {
	id      int
	handler EventHandler
}

type EventLoop struct {
	eventQueue chan Event
	done       chan bool // channel to signal event loop termination
}

func (loop *EventLoop) AddEvent(event Event) {
	loop.eventQueue <- event
}

func (loop *EventLoop) Start() {
	fmt.Println("Event loop started...")

	// start the loop in a separate thread
	go func() {
		for {
			select {
			case event := <-loop.eventQueue:
				go event.handler(event.id)
			case <-loop.done:
				return
			}
		}
	}()
}

func (loop *EventLoop) Stop() {
	fmt.Println("Stopping event loop...")
	loop.done <- true
}

func main() {

	// create and start the event loop
	loop := EventLoop{
		eventQueue: make(chan Event, 10), // buffered channel to hold events
		done:       make(chan bool),
	}

	loop.Start()

	// create and add events
	handlerFunc := func(id int) {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Microsecond)
		fmt.Printf("Event of type %d has been processed\n", id)
	}

	event1 := Event{id: 1, handler: handlerFunc}
	event2 := Event{id: 2, handler: handlerFunc}
	event3 := Event{id: 3, handler: handlerFunc}

	loop.AddEvent(event1)
	loop.AddEvent(event2)
	loop.AddEvent(event3)

	// give some time for events to finish
	time.Sleep(1 * time.Second)

	loop.Stop()
}
