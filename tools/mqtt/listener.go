package main

import (
	"github.com/awesome/homekit/app/config"
	"github.com/awesome/homekit/app/messages"
	"github.com/awesome/homekit/app/mosquitto"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/golang/protobuf/proto"
	"log"
)

func main() {
	cfg := &config.Config{MqttConfig: &config.MqttConfig{
		BrokerURL: "tcp://127.0.0.1:1883",
	}}

	client := mosquitto.InitMqttClient(cfg.MqttConfig)

	topic := "from/device/3652907"

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}

	callback2 := func(c mqtt.Client, msg mqtt.Message) {
		log.Printf("%#v", msg)

		var m messages.FromDevice

		if err := proto.Unmarshal(msg.Payload(), &m); err != nil {
			log.Println("err:", err)
		} else {
			switch true {
			case m.GetSwitchState() != nil:
				log.Printf("State: %t\r\n", m.GetSwitchState().GetState())

			case m.GetTemperatureValue() != nil:
				log.Printf("Temp: %f, hum: %f\r\n", m.GetTemperatureValue().GetTemperatureValue(), m.GetTemperatureValue().GetHumidityValue())
			}
		}
	}

	client.AddRoute(topic, callback2)

	select {}
}
