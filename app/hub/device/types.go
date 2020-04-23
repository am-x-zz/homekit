package device

import (
	"github.com/am-x/homekit/app/shared/messages"
	"gobot.io/x/gobot"
)

type HubDevice interface {
	GetGobotDevices() []gobot.Device
	OnMessage(message *messages.ToHub) error
}

type Factory func(h *Hub, cfg *messages.DeviceConfig) HubDevice
