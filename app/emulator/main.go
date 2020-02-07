package main

import (
	"errors"
	"log"

	"github.com/am-x/homekit/app/emulator/accessory"
	"github.com/am-x/homekit/app/emulator/device"
	"github.com/am-x/homekit/app/shared/messages"
	"github.com/am-x/homekit/app/shared/transport/mqtt"
	"github.com/brutella/hc"
	"gobot.io/x/gobot"
	"golang.org/x/sync/errgroup"

	hclog "github.com/brutella/hc/log"
	mqtta "gobot.io/x/gobot/platforms/mqtt"
)

var processors = make(map[uint32]accessory.MessageProcessor)

var MqttClientID = "homekit-central"

func main() {
	var (
		connections = make([]gobot.Connection, 0)
	)

	mqttAdaptor := mqtta.NewAdaptor("tcp://10.0.1.2:1883", MqttClientID)
	connections = append(connections, mqttAdaptor)

	messageTransport := mqtt.NewTransport(mqttAdaptor)

	hclog.Debug.Enable()

	pool := device.NewDevicePool()

	accessories := getAccessories()

	for _, a := range accessories {
		if at, ok := a.(accessory.WithHomeKitAccessory); ok {
			pool.Add(at)
		}
	}

	hc.OnTermination(func() {
		<-pool.Stop()
	})

	var g errgroup.Group

	g.Go(func() error {
		pool.Start()
		return errors.New("pool exited")
	})

	g.Go(func() error {
		robot := gobot.NewRobot("homekit", connections, func() {
			err := messageTransport.OnAccessoryMessage(func(message *messages.ToAccessory) error {
				log.Println("message", message)

				for _, a := range accessories {
					if at, ok := a.(accessory.WithAccessoryMessageProcessing); ok {
						if err := at.ProcessMessage(message); err != nil {
							log.Println(err)
						}
					}
				}

				return nil
			})

			if err != nil {
				log.Fatalln("can't subscribe")
			}
		})

		if err := robot.Start(); err != nil {
			return errors.New("robot exited")
		}

		return nil
	})

	if err := g.Wait(); err == nil {
		log.Fatalln("All exited.")
	}
}

func getAccessories() []accessory.Accessory {
	//acc1 := func() accessory.Accessory {
	//	a := accessory.NewSwitch("lamp1", "Lamp in bedroom")
	//
	//	// Log to console when client (e.g. iOS app) changes the value of the on characteristic
	//	a.Switch.On.OnValueRemoteUpdate(func(on bool) {
	//		if on == true {
	//			log.Debug.Println("Client changed switch to on")
	//		} else {
	//			log.Debug.Println("Client changed switch to off")
	//		}
	//	})
	//
	//	return a
	//}()

	acc2 := func() accessory.Accessory {
		a := accessory.NewLightbulb("bulb1", "Lightbulb")

		return a
	}()

	acc3 := func() accessory.Accessory {
		a := accessory.NewTemperatureSensor("temp1", "Temp sensor in bedroom", 1, -20, 40, 0.1)

		return a
	}()

	acc4 := func() accessory.Accessory {
		a := accessory.NewContactSensor("wnd1", "Window in bedroom")

		return a
	}()

	acc5 := func() accessory.Accessory {
		a := accessory.NewThermostat("thermo", "Thermostat in bedroom")

		return a
	}()

	return []accessory.Accessory{
		//acc1,
		acc2,
		acc3,
		acc4,
		acc5,
	}
}
