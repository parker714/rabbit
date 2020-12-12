Rabbit 
------------
This project is based on [ampq](https://github.com/streadway/amqp) and has been added auto reconnect and re-synchronization of client and more publish.

[![GitHub license](https://img.shields.io/github/license/parker714/rabbit)](https://github.com/parker714/rabbit/blob/master/LICENSE)
[![Release](https://img.shields.io/github/release/parker714/rabbit.svg)](https://github.com/parker714/rabbit/releases)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/parker714/rabbit)

Overview
------------
- AMQP 0.9.1 client with RabbitMQ extensions
- Auto reconnect
- Message replay

Supported
------------

rabbit is available as a [Go module](https://github.com/golang/go/wiki/Modules).

- go 1.12+ 

Installation
------------
```sh
go get github.com/parker714/rabbit
```

Usage
------------
```
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

	publish(s)
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

```

License
------------
This project is under the MIT License. See the [LICENSE](https://github.com/parker714/rabbit/blob/master/LICENSE) file for the full license text.
