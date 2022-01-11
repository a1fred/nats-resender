package options

import (
	"flag"
	"fmt"

	"github.com/nats-io/nats.go"
)

func NewOptions() *Options {
	options := &Options{
		FromNats: NewNatsOptions("from", "nats source"),
		ToNats:   NewNatsOptions("to", "nats source"),
	}
	// Subscription params flags
	flag.StringVar(&options.Queue, "queue", "", "queue, disabled by default")
	flag.IntVar(&options.PendingMsgLimit, "pendingMsgLimit", nats.DefaultSubPendingMsgsLimit, "pending subscription message limit")
	flag.IntVar(&options.PendingByteLimit, "pendingByteLimit", nats.DefaultSubPendingBytesLimit, "pending subscription byte limit")

	// App params
	flag.BoolVar(&options.Debug, "debug", false, "Debug")
	flag.UintVar(&options.StatPrintPeriod, "print-period", 10, "Print period seconds")
	flag.Parse()

	fmt.Printf(
		"Resending %s/%s -> %s/%s\n",
		options.FromNats.Url, options.FromNats.Subject,
		options.ToNats.Url, options.ToNats.Subject,
	)
	if options.Queue != "" {
		fmt.Printf("Queue: %s\n", options.Queue)
	}

	return options
}

type Options struct {
	FromNats *NatsOptions
	ToNats   *NatsOptions

	Queue            string
	PendingMsgLimit  int
	PendingByteLimit int

	Debug           bool
	StatPrintPeriod uint ``
}
