package internal

import (
	"log"

	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	Connection *nats.Conn
}

func NewNatsPublisher(Connection *nats.Conn) *NatsPublisher {
	return &NatsPublisher{
		Connection: Connection,
	}
}

func (nc *NatsPublisher) Publish(subject string, payload []byte) error {
	err := nc.Connection.Publish(subject, payload)
	if err != nil {
		return err
	}
	nc.Connection.Flush()
	return nil
}

type NatsSubscriber struct {
	Connection   *nats.Conn
	Subject      string
	ch           chan (*nats.Msg)
	subscription *nats.Subscription
}

func NewNatsSubscriber(connection *nats.Conn, subject string) *NatsSubscriber {
	return &NatsSubscriber{
		Connection:   connection,
		Subject:      subject,
		ch:           nil,
		subscription: nil,
	}
}

func (nc *NatsSubscriber) Subscribe() (chan (*nats.Msg), error) {
	ch := make(chan *nats.Msg, 64)
	sub, err := nc.Connection.ChanSubscribe(nc.Subject, ch)
	sub.SetPendingLimits(65536, 64*1024*1024)

	if err != nil {
		return nil, err
	}

	nc.ch = ch
	log.Printf("Started to subscribe: [%s]", nc.Subject)
	return nc.ch, nil
}

func (nc *NatsSubscriber) Unsubscribe() error {
	if nc.subscription == nil {
		return nil
	}
	return nc.subscription.Unsubscribe()
}
