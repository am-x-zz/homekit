package device

import (
	"fmt"
	"log"

	"github.com/am-x/homekit/app/hub"
	"github.com/am-x/homekit/app/shared/messages"
	"github.com/am-x/homekit/app/shared/transport"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type ContactSensor struct {
	h        *hub.Hub
	devices  []gobot.Device
	deviceID uint32
	isOpen   bool
}

func (c *ContactSensor) GetDevices() []gobot.Device {
	return c.devices
}

func NewContactSensor(h *hub.Hub, dc *messages.DeviceConfig) HubDevice {
	if cfg := dc.GetContactSensor(); cfg != nil {
		dr := h.GetDigitalReader()

		if dr == nil {
			log.Println("DigitalReader is nil")
			return nil
		}

		driver := gpio.NewButtonDriver(dr, fmt.Sprintf("%d", cfg.InPin))

		dev := &ContactSensor{
			h:        h,
			devices:  []gobot.Device{driver},
			deviceID: dc.DeviceID,
		}

		_ = driver.On(gpio.ButtonPush, dev.onEvent)
		_ = driver.On(gpio.ButtonRelease, dev.onEvent)
		_ = driver.On(gpio.Error, dev.onError)
		_ = h.Transport.OnHubMessage(h.HubID, dev.onMessage, transport.WithHubDeviceID(dc.DeviceID))

		return dev
	}

	return nil
}

func (c *ContactSensor) onMessage(message *messages.ToHub) error {
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

	//TODO: log err
}

func (c *ContactSensor) sendState() {
	m := &messages.ToAccessory{
		FromHub: c.h.HubID,
		Message: &messages.ToAccessory_ContactSensorState{ContactSensorState: &messages.ContactSensorState{DeviceID: c.deviceID, Open: c.isOpen}},
	}

	if err := c.h.Transport.ToAccessory(m); err != nil {
		log.Println("send to accessory", err)
	}

	log.Println("state is sent", c.isOpen)
}
