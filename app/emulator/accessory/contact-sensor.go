package accessory

import (
	"github.com/awesome/homekit/app/shared/messages"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type ContactSensor struct {
	id  string
	hid uint32

	hks *service.ContactSensor
	hka *accessory.Accessory
}

func NewContactSensor(id, name string) *ContactSensor {
	info := accessory.Info{
		Name:         name,
		Manufacturer: Manufacturer,
	}

	acc := &ContactSensor{
		id:  id,
		hka: accessory.New(info, accessory.TypeSensor),
	}

	acc.hks = service.NewContactSensor()
	acc.hka.AddService(acc.hks.Service)

	return acc
}

func (acc *ContactSensor) GetID() string {
	return acc.id
}

func (acc *ContactSensor) GetAccessory() *accessory.Accessory {
	return acc.hka
}

func (acc *ContactSensor) GetHardwareID() uint32 {
	return 1
}

func (acc *ContactSensor) ProcessMessage(msg *messages.ToAccessory) error {
	if m := msg.GetContactSensorState(); m != nil && m.GetDeviceID() == acc.GetHardwareID() {
		s := 0

		if m.GetOpen() {
			s = 1
		}

		acc.hks.ContactSensorState.SetValue(s)
	}

	return nil
}
