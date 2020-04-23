package accessory

import (
	"github.com/am-x/homekit/app/shared/messages"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type Lightbulb struct {
	id  string
	hid uint32

	hks *service.Lightbulb
	hka *accessory.Accessory
}

// NewLightbulb returns an light bulb accessory which one light bulb service.
func NewLightbulb(id, name string, hardwareID uint32) *Lightbulb {
	info := accessory.Info{
		Name:         name,
		Manufacturer: Manufacturer,
	}

	acc := Lightbulb{
		id:  id,
		hid: hardwareID,
		hks: service.NewLightbulb(),
		hka: accessory.New(info, accessory.TypeLightbulb),
	}

	acc.hka.AddService(acc.hks.Service)

	return &acc
}

func (acc Lightbulb) GetID() string {
	return acc.id
}

func (acc Lightbulb) GetAccessory() *accessory.Accessory {
	return acc.hka
}

func (acc *Lightbulb) GetHardwareID() uint32 {
	return acc.hid
}

func (acc *Lightbulb) GetService() *service.Lightbulb {
	return acc.hks
}

func (acc *Lightbulb) ProcessMessage(msg *messages.ToAccessory) error {
	if m := msg.GetSwitchState(); m != nil && m.GetDeviceID() == acc.GetHardwareID() {
		acc.hks.On.SetValue(m.GetState())
	}

	return nil
}
