package accessory

import (
	"github.com/am-x/homekit/app/shared/messages"
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type Thermostat struct {
	id  string
	hid uint32

	hks *service.Thermostat
	hka *accessory.Accessory
}

func NewThermostat(id, name string) *Thermostat {
	info := accessory.Info{
		Name:         name,
		Manufacturer: Manufacturer,
	}

	acc := &Thermostat{
		id:  id,
		hka: accessory.New(info, accessory.TypeSensor),
	}

	acc.hks = service.NewThermostat()

	//acc.hks.CurrentTemperature.SetValue(temp)
	//acc.hks.CurrentTemperature.SetMinValue(min)
	//acc.hks.CurrentTemperature.SetMaxValue(max)
	//acc.hks.CurrentTemperature.SetStepValue(steps)

	acc.hka.AddService(acc.hks.Service)

	return acc
}

func (acc *Thermostat) GetID() string {
	return acc.id
}

func (acc *Thermostat) GetAccessory() *accessory.Accessory {
	return acc.hka
}

func (acc *Thermostat) GetHardwareID() uint32 {
	return 4
}

func (acc *Thermostat) ProcessMessage(msg *messages.ToAccessory) error {
	//switch true {
	//case msg.GetTemperatureValue() != nil:
	//	acc.hks.CurrentTemperature.SetValue(float64(msg.GetTemperatureValue().GetTemperatureValue()))
	//}

	return nil
}
