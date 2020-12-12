package rabbit

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func (s *session) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	s.RLock()
	defer s.RUnlock()

	return s.Channel.Consume(
		queue,
		consumer,
		autoAck,
		exclusive,
		noLocal,
		noWait,
		args,
	)
}

// Consume used for consumption queue message.
// must handle message receipts(Ack or Nack) in callback func,
// otherwise the message will be returned to the queue after consumer disconnect.
func (s *session) ReConsume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table, callback func(amqp.Delivery)) {
	for {
		dd, err := s.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
		if err != nil {
			log.Printf("rabbit: reconsume failed, err: %s\n", err)
			<-time.After(reConnectionDelay)
			continue
		}

		for {
			select {
			case <-s.close:
				log.Println("rabbit: re-consume failed, session was closed, close consume")
				return
			case <-s.notifyClose:
				log.Println("rabbit: re-consume failed, connection or channel was closed, waiting to reconnect and re-consume")
				<-time.After(reConnectionDelay)
				goto ReConsume
			case msg := <-dd:
				go callback(msg)
			}
		}
	ReConsume:
	}
}
