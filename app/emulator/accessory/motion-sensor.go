package accessory

import (
	"github.com/am-x/homekit/app/shared/messages"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type MotionSensor struct {
	id  string
	hid uint32

	hks *service.MotionSensor
	hka *accessory.Accessory
}

func NewMotionSensor(id, name string, hardwareID uint32) *MotionSensor {
	info := accessory.Info{
		Name:         name,
		Manufacturer: Manufacturer,
	}

	acc := &MotionSensor{
		id:  id,
		hid: hardwareID,
		hka: accessory.New(info, accessory.TypeSensor),
	}

	acc.hks = service.NewMotionSensor()
	acc.hka.AddService(acc.hks.Service)

	return acc
}

func (acc *MotionSensor) GetID() string {
	return acc.id
}

func (acc *MotionSensor) GetAccessory() *accessory.Accessory {
	return acc.hka
}

func (acc *MotionSensor) GetHardwareID() uint32 {
	return acc.hid
}

func (acc *MotionSensor) ProcessMessage(msg *messages.ToAccessory) error {
	if m := msg.GetContactSensorState(); m != nil && m.GetDeviceID() == acc.GetHardwareID() {
		acc.hks.MotionDetected.SetValue(m.GetOpen())
	}

	return nil
}
