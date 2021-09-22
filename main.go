package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var revision = "unknown"

func main() {
	var fromUrl, toUrl, fromSubject, toSubject, queue string
	var ncFrom, ncTo *nats.Conn
	var statPrintPeriod, processedCounter uint
	var pendingMsgLimit, pendingByteLimit int
	var debug bool
	var err error

	_, err = os.Stderr.WriteString(fmt.Sprintf("nats-resender version %s\n", revision))
	if err != nil {
		log.Fatalln(err)
	}

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

	ncFrom, err = nats.Connect(fromUrl)
	if err != nil {
		log.Fatalln(fmt.Errorf("source nats (%s) connect error: %s", fromUrl, err))
	}
	defer ncFrom.Close()
	if fromUrl == toUrl {
		ncTo = ncFrom
	} else {
		ncTo, err = nats.Connect(toUrl)
		if err != nil {
			log.Fatalln(fmt.Errorf("destination nats (%s) to connect error: %s", toUrl, err))
		}
		defer ncTo.Close()
	}

	subscription, err := ncFrom.QueueSubscribe(fromSubject, queue, func(m *nats.Msg) {
		err = ncTo.Publish(toSubject, m.Data)
		if err != nil {
			log.Fatalln(fmt.Errorf("error resending message: %s", err))
		}

		if debug {
			log.Printf("MSG '%s' -> '%s': '%s'\n", m.Subject, toSubject, m.Data)
		}
		processedCounter++
	})
	if err != nil {
		log.Fatalln(fmt.Errorf("subscribe to '%s' error: %s", fromSubject, err))
	}
	err = subscription.SetPendingLimits(pendingMsgLimit, pendingByteLimit)
	if err != nil {
		log.Fatalln(fmt.Errorf("subscribe set pending limits error: %s", err))
	}

	for {
		start := time.Now()
		time.Sleep(time.Duration(statPrintPeriod) * time.Second)
		elapsed := time.Since(start).Seconds()
		log.Printf("%d messages processed, elapsed %.2fs, %.2f msg/sec", processedCounter, elapsed, float64(processedCounter)/elapsed)
		processedCounter = 0
	}
}
