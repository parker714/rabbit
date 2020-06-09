package rabbit

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func (s *Session) Consume(queueName string) (<-chan amqp.Delivery, error) {
	s.RLock()
	defer s.RUnlock()

	return s.Channel.Consume(
		queueName,
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
}

// Consume used for consumption queue message.
// must handle message receipts(Ack or Nack) in callback func,
// otherwise the message will be returned to the queue after consumer disconnect.
func (s *Session) ReConsume(queueName string, callback func(amqp.Delivery)) {
	for {
		dd, err := s.Consume(queueName)
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
				goto RECONSUME
			case msg := <-dd:
				go callback(msg)
			}
		}
	RECONSUME:
	}
}
