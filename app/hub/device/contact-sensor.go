package device

import (
	"fmt"
	"log"

	"github.com/am-x/homekit/app/shared/messages"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type ContactSensor struct {
	h        *Hub
	devices  []gobot.Device
	deviceID uint32
	isOpen   bool
}

func (c *ContactSensor) GetGobotDevices() []gobot.Device {
	return c.devices
}

func NewContactSensor(h *Hub, dc *messages.DeviceConfig) HubDevice {
	if cfg := dc.GetContactSensor(); cfg != nil {
		dr := h.GetDigitalReader()

		if dr == nil {
			log.Println("DigitalReader is nil")
			return nil
		}

		driver := gpio.NewButtonDriver(dr, fmt.Sprintf("%d", cfg.GetInPin()))

		dev := &ContactSensor{
			h:        h,
			devices:  []gobot.Device{driver},
			deviceID: dc.DeviceID,
		}

		_ = driver.On(gpio.ButtonPush, dev.onEvent)
		_ = driver.On(gpio.ButtonRelease, dev.onEvent)
		_ = driver.On(gpio.Error, dev.onError)

		return dev
	}

	return nil
}

func (c *ContactSensor) OnMessage(message *messages.ToHub) error {
	if msg := message.GetGetContactSensorState(); msg != nil && msg.DeviceID == c.deviceID {
		c.sendState()
	}

	return nil
}

func (c *ContactSensor) onError(s interface{}) {
	log.Println("Err", s)
}

func (c *ContactSensor) onEvent(s interface{}) {
	if isOpen, ok := s.(int); ok {
		c.isOpen = !(isOpen == 1)

		c.sendState()
		return
	}
}

func (c *ContactSensor) sendState() {
	c.h.MessageBus <- &messages.ToAccessory{
		FromHub: c.h.HubID,
		Message: &messages.ToAccessory_ContactSensorState{ContactSensorState: &messages.ContactSensorState{DeviceID: c.deviceID, Open: c.isOpen}},
	}

	log.Println("state is sent", c.isOpen)
}
