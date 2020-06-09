package rabbit

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
	"time"
)

const (
	// When re-publish times after the publish fails
	rePublishTimes = 3
)

var (
	ErrOverPublishTimes = errors.New("rabbit: re-publish failed, more than rePublishTimes")
	ErrServerClosed     = errors.New("rabbit: server was closed")
)

func (s *Session) Publish(exchangeName, routingKey string, msg []byte) error {
	s.RLock()
	defer s.RUnlock()

	return s.Channel.Publish(
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			// Marking messages as persistent doesn't fully guarantee that a message won't be lost.
			// if you need a stronger guarantee then you can use publish confirms.
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         msg,
		},
	)
}

// RePublish is a rabbit server confirmation publish mode.
// when the publish fails, it will be published back multiple times by default. (eg: rePublishTimes)
func (s *Session) RePublish(exchangeName, routingKey string, msg []byte) error {
	var publishTimes int
	for {
		if publishTimes > rePublishTimes {
			return ErrOverPublishTimes
		}
		publishTimes++

		if err := s.Publish(exchangeName, routingKey, msg); err != nil {
			log.Printf("rabbit: RePublish failed, err: %s\n", err)
			continue
		}

		// Waiting for the rabbit server to notify the release
		for {
			select {
			case <-s.close:
				return ErrServerClosed
			case <-s.notifyClose:
				log.Println("rabbit: re-publish failed, connection or channel was closed, waiting to reconnect and re-publish")
				<-time.After(reConnectionDelay)
				goto PUBLISH
			case confirm := <-s.notifyConfirm:
				if !confirm.Ack {
					log.Println("rabbit: re-publish didn't confirm, retry")
					goto PUBLISH
				}
			}
			return nil
		}
	PUBLISH:
	}
}
