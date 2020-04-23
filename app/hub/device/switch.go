package device

import (
	"fmt"
	"log"

	"github.com/am-x/homekit/app/shared/messages"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type Switch struct {
	h         *Hub
	deviceID  uint32
	devices   []gobot.Device
	outPin    *gpio.DirectPinDriver
	isEnabled bool
}

func (c *Switch) GetGobotDevices() []gobot.Device {
	return c.devices
}

func NewSwitch(h *Hub, dc *messages.DeviceConfig) HubDevice {
	if cfg := dc.GetSwitch(); cfg != nil {
		dr := h.GetDigitalReader()

		if dr == nil {
			log.Println("DigitalReader is nil")
			return nil
		}

		outPin := gpio.NewDirectPinDriver(h.Adaptor, fmt.Sprintf("%d", cfg.GetOutputPin()))

		dev := &Switch{
			h:        h,
			deviceID: dc.DeviceID,
			devices:  []gobot.Device{outPin},
			outPin:   outPin,
		}

		for _, p := range cfg.GetInputs() {
			driver := gpio.NewButtonDriver(dr, fmt.Sprintf("%d", p.GetPin()))

			switch p.GetInputType() {
			case messages.SwitchInputType_SwitchInputTypeOnOff:
				_ = driver.On(gpio.ButtonPush, dev.onEventOnOff)
				_ = driver.On(gpio.ButtonRelease, dev.onEventOnOff)
			case messages.SwitchInputType_SwitchInputTypePulse:
				_ = driver.On(gpio.ButtonPush, dev.onEventPulse)
				//_ = driver.On(gpio.ButtonRelease, dev.onEventPulse)
			case messages.SwitchInputType_SwitchInputTypeChange:
				_ = driver.On(gpio.ButtonPush, dev.onEventChange)
				_ = driver.On(gpio.ButtonRelease, dev.onEventChange)
			default:
				log.Println("not supported input type")
				return nil
			}

			_ = driver.On(gpio.Error, dev.onError)

			dev.devices = append(dev.devices, driver)
		}

		return dev
	}

	return nil
}

func (c *Switch) OnMessage(message *messages.ToHub) error {
	if msg := message.GetGetSwitchState(); msg != nil && msg.DeviceID == c.deviceID {
		c.sendState()
	}

	if msg := message.GetSetSwitchState(); msg != nil && msg.DeviceID == c.deviceID {
		c.isEnabled = msg.GetState()
		c.updateOutput()
	}

	return nil
}

func (c *Switch) onError(s interface{}) {
	log.Println("Err", s)
}

func (c *Switch) onEventOnOff(s interface{}) {
	if isOn, ok := s.(int); ok {
		c.isEnabled = isOn == 1

		c.updateOutput()
		c.sendState()
	}
}

func (c *Switch) onEventPulse(_ interface{}) {
	c.toggle()
	c.updateOutput()
	c.sendState()
}

func (c *Switch) onEventChange(_ interface{}) {
	c.toggle()
	c.updateOutput()
	c.sendState()
}

func (c *Switch) toggle() {
	c.isEnabled = !c.isEnabled
}

func (c *Switch) updateOutput() {
	if c.isEnabled {
		_ = c.outPin.On()
	} else {
		_ = c.outPin.Off()
	}

	log.Println("state is", c.isEnabled)
}

func (c *Switch) sendState() {
	c.h.MessageBus <- &messages.ToAccessory{
		FromHub: c.h.HubID,
		Message: &messages.ToAccessory_SwitchState{SwitchState: &messages.SwitchState{DeviceID: c.deviceID, State: c.isEnabled}},
	}

	log.Println("state is sent", c.isEnabled)
}
