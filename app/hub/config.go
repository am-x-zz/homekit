package hub

import "github.com/awesome/homekit/app/shared/messages"

func LoadHubConfig() *messages.HubConfig {
	return &messages.HubConfig{
		HubID: 1,
		Devices: []*messages.DeviceConfig{
			{
				DeviceID: 1,
				Config: &messages.DeviceConfig_ContactSensor{
					ContactSensor: &messages.ContactSensorConfig{
						InPin: 13,
					},
				},
			},
		},
	}
}
