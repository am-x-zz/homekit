package main

import (
	"github.com/nats-io/go-nats"
	"log"
	"time"
)

func main() {
	nc, err := nats.Connect("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := nc.Publish("blink", []byte("All is Well")); err != nil {
				log.Fatal(err)
			}

			_ = nc.Flush()
		}
	}
}
