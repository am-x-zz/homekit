package hub

import "github.com/am-x/homekit/app/shared/messages"

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
			{
				DeviceID: 2,
				Config: &messages.DeviceConfig_Switch{
					Switch: &messages.SwitchConfig{
						Inputs: []*messages.SwitchInput{
							{Pin: 1, InputType: messages.SwitchInputType_SwitchInputTypeChange},
							{Pin: 2, InputType: messages.SwitchInputType_SwitchInputTypeChange},
							{Pin: 3, InputType: messages.SwitchInputType_SwitchInputTypeChange},
						},
						OutputPin: 5,
					},
				},
			},
		},
	}
}
