package cli

import (
	"flag"
	"log"
	"time"

	"github.com/a1fred/nats-resender/src/resender"
	"github.com/nats-io/nats.go"
)

func Main() {
	var fromUrl, toUrl, fromSubject, toSubject, queue string

	var statPrintPeriod uint
	var pendingMsgLimit, pendingByteLimit int
	var debug bool
	var err error

	// From nats
	flag.StringVar(&fromUrl, "from-url", nats.DefaultURL, "nats source url")
	flag.StringVar(&fromSubject, "from-subj", "*", "nats source subject")

	// To nats
	flag.StringVar(&toUrl, "to-url", nats.DefaultURL, "nats destination url")
	flag.StringVar(&toSubject, "to-subj", "nats-resender", "nats destination subject")

	// Subscription params flags
	flag.StringVar(&queue, "queue", "", "queue, disabled by default")
	flag.IntVar(&pendingMsgLimit, "pendingMsgLimit", nats.DefaultSubPendingMsgsLimit, "pending subscription message limit")
	flag.IntVar(&pendingByteLimit, "pendingByteLimit", nats.DefaultSubPendingBytesLimit, "pending subscription byte limit")

	// App params
	flag.BoolVar(&debug, "debug", false, "Debug")
	flag.UintVar(&statPrintPeriod, "print-period", 10, "Print period seconds")

	flag.Parse()

	nn := resender.NewResender(fromUrl, toUrl)
	defer nn.Close()

	_, err = nn.Subscribe(fromSubject, queue, toSubject, pendingMsgLimit, pendingByteLimit, debug)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		start := time.Now()
		time.Sleep(time.Duration(statPrintPeriod) * time.Second)
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
