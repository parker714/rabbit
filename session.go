package rabbit

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

const (
	// When reconnecting to the server after connection failure
	reConnectionDelay = 3 * time.Second

	// When open channel after a channel exception
	reOpenChannelDelay = 3 * time.Second

	defaultHeartbeat = 10 * time.Second
	defaultLocale    = "en_US"
)

type Session interface {
	Close() error
	IsClosed() bool
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	RePublish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	ReConsume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table, callback func(amqp.Delivery))
}

type session struct {
	sync.RWMutex

	connection *amqp.Connection
	Channel    *amqp.Channel

	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation

	notifyClose chan struct{}
	close       chan struct{}
}

func New(url string) (Session, error) {
	config := amqp.Config{
		Heartbeat: defaultHeartbeat,
		Locale:    defaultLocale,
	}

	return NewConfig(url, config)
}

func NewConfig(url string, config amqp.Config) (Session, error) {
	var err error
	s := &session{
		notifyClose: make(chan struct{}),
		close:       make(chan struct{}),
	}

	if s.connection, err = amqp.DialConfig(url, config); err != nil {
		return nil, err
	}
	if err = s.openChannel(); err != nil {
		return nil, err
	}

	go s.handleConnect(url, config)
	return s, err
}

func (s *session) openChannel() (err error) {
	s.Lock()
	defer s.Unlock()

	s.Channel, err = s.connection.Channel()
	if err != nil {
		return
	}

	if err = s.Channel.Confirm(false); err != nil {
		return
	}

	s.notifyChanClose = make(chan *amqp.Error)
	s.Channel.NotifyClose(s.notifyChanClose)

	s.notifyConfirm = make(chan amqp.Confirmation, 1)
	s.Channel.NotifyPublish(s.notifyConfirm)

	return
}

func (s *session) handleConnect(url string, config amqp.Config) {
	var err error
	for {
		s.connection, err = amqp.DialConfig(url, config)
		if err != nil {
			log.Printf("rabbit: connection failed, err: %s\n", err)
			<-time.After(reConnectionDelay)
			continue
		}

		s.notifyConnClose = make(chan *amqp.Error)
		s.connection.NotifyClose(s.notifyConnClose)

		if done := s.handleChannel(); done {
			break
		}
	}
}

func (s *session) handleChannel() bool {
	for {
		if err := s.openChannel(); err != nil {
			log.Printf("rabbit: channel open failed, err: %s\n", err)
			<-time.After(reOpenChannelDelay)
			continue
		}

		select {
		case <-s.close:
			log.Println("rabbit: session was closed")
			return true
		case <-s.notifyConnClose:
			log.Println("rabbit: notifyConnClose received, re-connection")
			s.notifyClose <- struct{}{}
			return false
		case <-s.notifyChanClose:
			log.Println("rabbit: notifyChanClose received, re-openChannel")
			s.notifyClose <- struct{}{}
		}
	}
}

func (s *session) Close() error {
	s.RLock()
	defer s.RUnlock()

	close(s.close)
	if s.Channel != nil {
		if err := s.Channel.Close(); err != nil {
			return err
		}
	}
	if s.connection != nil {
		return s.connection.Close()
	}

	return nil
}

func (s *session) IsClosed() bool {
	s.RLock()
	defer s.RUnlock()

	if s.connection == nil || s.Channel == nil {
		return true
	}

	return s.connection.IsClosed()
}
