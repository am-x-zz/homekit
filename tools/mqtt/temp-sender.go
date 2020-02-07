package main

import (
	"log"
	"time"

	"github.com/am-x/homekit/app/config"
	"github.com/am-x/homekit/app/messages"
	"github.com/am-x/homekit/app/mosquitto"
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

			msg := &messages.FromDevice{
				DeviceID: 2,
				Message: &messages.FromDevice_TemperatureValue{
					TemperatureValue: &messages.TemperatureValue{
						TemperatureValue: 24.0,
					},
				},
			}

			b, err := proto.Marshal(msg)

			if err != nil {
				log.Panic(err)
			}

			log.Println(b)

			if token := client.Publish("from/device/3652907", 0, false, b); token.Wait() && token.Error() != nil {
				log.Panic(token.Error())
			}
		}
	}
}
