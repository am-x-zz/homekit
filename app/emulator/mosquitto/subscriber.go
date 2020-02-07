package mosquitto

import (
	"github.com/awesome/homekit/app/config"
	"github.com/eclipse/paho.mqtt.golang"
)

func InitMqttClient(cfg *config.MqttConfig) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.BrokerURL)
	//opts.SetClientID(*id)
	//opts.SetUsername(*user)
	//opts.SetPassword(*password)
	//opts.SetCleanSession(*cleansess)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}
