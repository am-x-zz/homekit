package main

import (
	"fmt"
	"github.com/am-x/homekit/app/hub/mock"
	"github.com/am-x/homekit/app/shared/messages"
	"github.com/am-x/homekit/app/shared/transport"
	"github.com/pkg/errors"
	"log"

	"github.com/am-x/homekit/app/hub"
	"github.com/am-x/homekit/app/hub/device"
	"github.com/am-x/homekit/app/shared/transport/mqtt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"

	mqtta "gobot.io/x/gobot/platforms/mqtt"
)

const (
	AdaptorRaspi = "raspi"
	AdaptorMock  = "mock"
)

func main() {
	var (
		h                *device.Hub
		adaptor          gobot.Adaptor
		messageTransport transport.MessageTransport

		adaptorType = "mock"

		connections = make([]gobot.Connection, 0)
	)

	switch adaptorType {
	case AdaptorRaspi:
		adaptor = raspi.NewAdaptor()
	case AdaptorMock:
		adaptor = mock.NewAdaptor()
	default:
		log.Fatal(fmt.Sprintf("Unsupported adaptor type: %s", adaptorType))
	}

	mqttAdaptor := mqtta.NewAdaptor("tcp://10.0.1.2:1883", "homekit-hub")
	connections = append(connections, mqttAdaptor)
	connections = append(connections, adaptor)

	messageTransport = mqtt.NewTransport(mqttAdaptor)

	h = device.NewHub(
		hub.LoadHubConfig(),
		adaptor,
		[]device.Factory{
			device.NewContactSensor,
			device.NewSwitch,
		},
	)

	//cfg := i2c.NewMCP23017Driver(piAdaptor, i2c.WithBus(1), i2c.WithAddress(0x20))

	//led := gpio.NewLedDriver(piAdaptor, "12")
	//button := gpio.NewButtonDriver(piAdaptor, "13")

	//v := 0

	//work := func() {
	//gobot.Every(2*time.Second, func() {
	//	if v == 0 {
	//		v = 1
	//	} else {
	//		v = 0
	//	}

	//if err := cfg.WriteGPIO(uint8(2), uint8(v), "A"); err != nil {
	//	log.Println(err)
	//}
	//
	//log.Println(fmt.Sprintf("state: %v", v))
	//})

	//	_ = button.On(gpio.ButtonPush, func(s interface{}) {
	//		log.Println(fmt.Sprintf("button push!"))
	//	})
	//}

	//devices = append(devices, cfg)

	robot := gobot.NewRobot(
		"homekit",
		connections,
		h.GetGobotDevices(),
		func() {
			err := messageTransport.OnHubMessage(h.HubID, func(message *messages.ToHub) error {
				h.OnMessage(message)
				return nil
			})

			if err != nil {
				log.Println(errors.Wrap(err, "subscribe on hub message"))
			}

			go func() {
				for {
					select {
					case msg := <-h.MessageBus:
						if err := messageTransport.ToAccessory(msg); err != nil {
							log.Println(errors.Wrap(err, "send message to accessory"))
						}
					}
				}
			}()
		},
	)

	if err := robot.Start(); err != nil {
		log.Fatal(err)
	}
}
