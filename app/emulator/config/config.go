package config

type Config struct {
	MqttConfig *MqttConfig `json:"mqtt_config"`
}

type MqttConfig struct {
	BrokerURL string `json:"broker_url"`
}
