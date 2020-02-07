package main

import (
	"github.com/awesome/homekit/app/shared/messages"
	"github.com/awesome/homekit/app/shared/transport/mqtt"
	"gobot.io/x/gobot"
	"log"

	mqtta "gobot.io/x/gobot/platforms/mqtt"
)

func main() {
	var (
		connections = make([]gobot.Connection, 0)
	)

	mqttAdaptor := mqtta.NewAdaptor("tcp://10.0.1.2:1883", "homekit2")
	connections = append(connections, mqttAdaptor)

	messageTransport := mqtt.NewTransport(mqttAdaptor)

	robot := gobot.NewRobot("homekit", connections, func() {
		_ = messageTransport.OnAccessoryMessage(func(message *messages.ToAccessory) error {
			log.Println(message)
			return nil
		})
	})

	if err := robot.Start(); err != nil {
		log.Fatal("robot exited")
	}
}
