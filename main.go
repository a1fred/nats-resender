package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/a1fred/nats-resender/src"
)

var revision = "unknown"

func main() {
	_, err := os.Stderr.WriteString(fmt.Sprintf("nats-resender version %s\n", revision))
	if err != nil {
		log.Fatalln(err)
	}

	cli.Main()
}
