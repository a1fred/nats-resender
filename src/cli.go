package cli

import (
	"flag"
	"log"
	"time"

	"github.com/a1fred/nats-resender/src/options"
	"github.com/a1fred/nats-resender/src/resender"
)

func Main() {
	options := options.NewOptions()
	var err error

	flag.Parse()

	nn := resender.NewResender(options.FromNats, options.ToNats)
	defer nn.Close()

	_, err = nn.Subscribe(
		options.FromNats.Subject,
		options.Queue,
		options.ToNats.Subject,

		options.PendingMsgLimit, options.PendingByteLimit, options.Debug)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		start := time.Now()
		time.Sleep(time.Duration(options.StatPrintPeriod) * time.Second)
		elapsed := time.Since(start).Seconds()
		processedCounter, q_msgs, q_bytes := nn.FlushCounter()

		log.Printf(
			"%d processed, %.2fs elapsed, %.2f msg/sec, buffer:%dmsgs/%dbytes",
			processedCounter,
			elapsed,
			float64(processedCounter)/elapsed,
			q_msgs, q_bytes,
		)
	}
}
