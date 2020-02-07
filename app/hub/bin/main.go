package main

import (
	"log"

	"github.com/am-x/homekit/app/hub"
	"github.com/am-x/homekit/app/hub/device"
	"github.com/am-x/homekit/app/shared/transport/mqtt"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/raspi"

	mqtta "gobot.io/x/gobot/platforms/mqtt"
)

func main() {
	var (
		h = &hub.Hub{}

		connections = make([]gobot.Connection, 0)
		devices     = make([]gobot.Device, 0)

		initializers = []device.Initializer{device.NewContactSensor}
	)

	hubConfig := hub.LoadHubConfig()

	mqttAdaptor := mqtta.NewAdaptor("tcp://10.0.1.2:1883", "homekit-hub")
	connections = append(connections, mqttAdaptor)

	h.HubID = hubConfig.HubID
	h.Adaptor = raspi.NewAdaptor()
	h.Transport = mqtt.NewTransport(mqttAdaptor)

	connections = append(connections, h.Adaptor)

	for _, cfg := range hubConfig.Devices {
		for _, init := range initializers {
			if dev := init(h, cfg); dev != nil {
				devices = append(devices, dev.GetDevices()...)
			}
		}
	}

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

	robot := gobot.NewRobot("homekit", connections, devices)

	if err := robot.Start(); err != nil {
		log.Fatal(err)
	}
}
