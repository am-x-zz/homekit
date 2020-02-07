package main

import (
	"log"
	"time"

	"github.com/awesome/homekit/app/shared/messages"
	"github.com/awesome/homekit/app/shared/transport/mqtt"
	"gobot.io/x/gobot"

	mqtta "gobot.io/x/gobot/platforms/mqtt"
)

func main() {
	var (
		connections = make([]gobot.Connection, 0)
	)

	mqttAdaptor := mqtta.NewAdaptor("tcp://10.0.1.2:1883", "homekit")
	connections = append(connections, mqttAdaptor)

	messageTransport := mqtt.NewTransport(mqttAdaptor)

	state := false

	robot := gobot.NewRobot("homekit", connections, func() {
		gobot.Every(5*time.Second, func() {
			state = !state

			msg := &messages.ToAccessory{
				FromHub: 1,
				Message: &messages.ToAccessory_ContactSensorState{
					ContactSensorState: &messages.ContactSensorState{
						DeviceID: 1,
						Open:     state,
					},
				},
			}

			log.Println("sent", msg)

			if err := messageTransport.ToAccessory(msg); err != nil {
				log.Panic(err)
			}
		})
	})

	if err := robot.Start(); err != nil {
		log.Fatal("robot exited")
	}
}
