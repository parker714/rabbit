package main

import (
	"github.com/parker714/rabbit"
	"github.com/streadway/amqp"
	"log"
	"time"
)

// go run consume.go
func main() {
	s, err := rabbit.New("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalln(err)
	}
	go monitor(s)

	// test
	publish(s)
}

func monitor(s *rabbit.Session) {
	for {
		log.Printf("rabbit: session status(%t)\n", !s.IsClosed())
		time.Sleep(1 * time.Second)
	}
}

// publish example
func publish(s *rabbit.Session) {
	if err := s.RePublish("test", "#", []byte("re-hello")); err != nil {
		log.Fatalln(err)
	}
}

// consume example
func consume(s *rabbit.Session) {
	s.ReConsume("test-consume", handle)
}

// consume example, when you have finished processing the message, be sure to Ack or Nack !!!
func handle(d amqp.Delivery) {
	//defer d.Ack(false)
	//defer _ = d.Nack(false, false)
	log.Println(string(d.Body))
}
