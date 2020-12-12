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
	// ErrOverPublishTimes ...
	ErrOverPublishTimes = errors.New("rabbit: re-publish failed, more than rePublishTimes")
	// ErrServerClosed ...
	ErrServerClosed = errors.New("rabbit: server was closed")
)

// Publish used send msg.
func (s *session) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	s.RLock()
	defer s.RUnlock()

	return s.Channel.Publish(
		exchange,
		key,
		mandatory,
		immediate,
		msg,
	)
}

// RePublish is a rabbit server confirmation publish mode.
// when the publish fails, it will be published back multiple times by default. (eg: rePublishTimes)
func (s *session) RePublish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	var publishTimes int
	for {
		if publishTimes > rePublishTimes {
			return ErrOverPublishTimes
		}
		publishTimes++

		if err := s.Publish(exchange, key, mandatory, immediate, msg); err != nil {
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
