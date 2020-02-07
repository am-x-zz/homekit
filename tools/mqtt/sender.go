package main

import (
	"log"
	"time"

	"github.com/awesome/homekit/app/config"
	"github.com/awesome/homekit/app/messages"
	"github.com/awesome/homekit/app/mosquitto"
	"github.com/golang/protobuf/proto"
)

func main() {
	cfg := &config.Config{MqttConfig: &config.MqttConfig{
		BrokerURL: "tcp://127.0.0.1:1883",
	}}

	client := mosquitto.InitMqttClient(cfg.MqttConfig)

	t := time.NewTicker(5 * time.Second)
	defer t.Stop()

	state := false

	for {
		select {
		case <-t.C:
			state = !state

			msg := &messages.MessageWrapper{
				Message: &messages.MessageWrapper_ToDevice{
					ToDevice: &messages.ToDevice{
						DeviceID: 1,
						Message: &messages.ToDevice_SetSwitchState{
							SetSwitchState: &messages.SetSwitchState{State: state},
						},
					},
				},
			}

			b, err := proto.Marshal(msg)

			if err != nil {
				log.Panic(err)
			}

			log.Println(b)

			if token := client.Publish("to/device/3652907", 0, false, b); token.Wait() && token.Error() != nil {
				log.Panic(token.Error())
			}
		}
	}
}
