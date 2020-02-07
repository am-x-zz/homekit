package main

import (
	"github.com/nats-io/go-nats"
	"log"
)

func main() {
	nc, err := nats.Connect("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe
	if _, err := nc.Subscribe("blink", func(m *nats.Msg) {
		log.Println(m)
	}); err != nil {
		log.Fatal(err)
	}

	select {}
}
