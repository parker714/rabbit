package main

import (
	"github.com/parker714/rabbit"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func main() {
	s, err := rabbit.New("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalln(err)
	}
	go monitor(s)

	// test
	publish(s)
}

func monitor(s rabbit.Session) {
	for {
		log.Printf("rabbit: session status(%t)\n", !s.IsClosed())
		time.Sleep(1 * time.Second)
	}
}

// publish example
func publish(s rabbit.Session) {
	msg := amqp.Publishing{
		// Marking messages as persistent doesn't fully guarantee that a message won't be lost.
		// if you need a stronger guarantee then you can use publish confirms.
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         []byte("pb"),
	}

	if err := s.RePublish("test", "#", false, false, msg); err != nil {
		log.Fatalln(err)
	}
}

// consume example
func consume(s rabbit.Session) {
	table := amqp.Table{
		"uid": 121,
	}
	s.ReConsume("test-consume", "pb", true, false, false, false, table, handle)
}

// consume example, when you have finished processing the message, be sure to Ack or Nack !!!
func handle(d amqp.Delivery) {
	//defer d.Ack(false)
	//defer _ = d.Nack(false, false)
	log.Println(string(d.Body))
}
