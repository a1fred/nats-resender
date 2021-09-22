package resender

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type Resender struct {
	ncFrom *nats.Conn
	ncTo   *nats.Conn
	sub    *nats.Subscription

	processedCounter uint
}

func NewResender(fromUrl, toUrl string) *Resender {
	var ncFrom, ncTo *nats.Conn

	ncFrom, err := nats.Connect(fromUrl)
	if err != nil {
		log.Fatalln(fmt.Errorf("source nats (%s) connect error: %s", fromUrl, err))
	}

	if fromUrl == toUrl {
		ncTo = ncFrom
	} else {
		ncTo, err = nats.Connect(toUrl)
		if err != nil {
			log.Fatalln(fmt.Errorf("destination nats (%s) to connect error: %s", toUrl, err))
		}
		defer ncTo.Close()
	}

	return &Resender{
		ncFrom: ncFrom,
		ncTo:   ncTo,
	}
}

func (r *Resender) Close() error {
	if r.sub != nil {
		err := r.sub.Drain()
		if err != nil {
			return err
		}
		r.sub = nil
	}

	err := r.ncFrom.Drain()
	if err != nil {
		return err
	}
	r.ncFrom.Close()

	if r.ncFrom != r.ncTo {
		r.ncTo.Close()
	}

	return nil
}

func (r *Resender) FlushCounter() (uint, int, int) {
	counter := r.processedCounter
	r.processedCounter = 0

	q_msgs, q_bytes, err := r.sub.Pending()
	if err != nil {
		q_msgs, q_bytes = 0, 0
	}

	return counter, q_msgs, q_bytes
}

func (r *Resender) Subscribe(
	fromSubject, queue, toSubject string,
	pendingMsgLimit, pendingByteLimit int,
	debug bool,
) (*nats.Subscription, error) {
	subscription, err := r.ncFrom.QueueSubscribe(fromSubject, queue, func(m *nats.Msg) {
		err := r.ncTo.Publish(toSubject, m.Data)
		if err != nil {
			log.Fatalln(fmt.Errorf("error resending message: %s", err))
		}

		if debug {
			log.Printf("MSG '%s' -> '%s': '%s'\n", m.Subject, toSubject, m.Data)
		}
		r.processedCounter++
	})

	if err != nil {
		return nil, fmt.Errorf("subscribe to '%s' error: %s", fromSubject, err)
	}
	err = subscription.SetPendingLimits(pendingMsgLimit, pendingByteLimit)
	if err != nil {
		return nil, fmt.Errorf("subscribe set pending limits error: %s", err)
	}

	r.sub = subscription
	return subscription, nil
}
