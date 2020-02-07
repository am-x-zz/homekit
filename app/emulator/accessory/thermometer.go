package accessory

import (
	"github.com/am-x/homekit/app/shared/messages"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type Thermometer struct {
	id  string
	hid uint32

	hks *service.TemperatureSensor
	hka *accessory.Accessory
}

// NewTemperatureSensor returns a Thermometer which implements model.Thermometer.
func NewTemperatureSensor(id, name string, temp, min, max, steps float64) *Thermometer {
	info := accessory.Info{
		Name:         name,
		Manufacturer: Manufacturer,
	}

	acc := &Thermometer{
		id:  id,
		hka: accessory.New(info, accessory.TypeSensor),
	}

	acc.hks = service.NewTemperatureSensor()

	acc.hks.CurrentTemperature.SetValue(temp)
	acc.hks.CurrentTemperature.SetMinValue(min)
	acc.hks.CurrentTemperature.SetMaxValue(max)
	acc.hks.CurrentTemperature.SetStepValue(steps)

	acc.hka.AddService(acc.hks.Service)

	return acc
}

func (acc *Thermometer) GetID() string {
	return acc.id
}

func (acc *Thermometer) GetAccessory() *accessory.Accessory {
	return acc.hka
}

func (acc *Thermometer) SetValue(v float64) {
	acc.hks.CurrentTemperature.SetValue(v)
}

func (acc *Thermometer) GetHardwareID() uint32 {
	return 2
}

func (acc *Thermometer) ProcessMessage(msg *messages.ToAccessory) error {
	//switch true {
	//case msg.GetTemperatureValue() != nil:
	//	acc.SetValue(float64(msg.GetTemperatureValue().GetTemperatureValue()))
	//}

	return nil
}
