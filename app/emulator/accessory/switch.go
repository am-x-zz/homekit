package accessory

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type Switch struct {
	*accessory.Accessory
	ID     string
	Switch *service.Switch
}

// NewSwitch returns a switch which implements model.Switch.
func NewSwitch(id, name string) *Switch {
	info := accessory.Info{
		Name:         name,
		Manufacturer: Manufacturer,
	}

	acc := Switch{
		ID:        id,
		Switch:    service.NewSwitch(),
		Accessory: accessory.New(info, accessory.TypeSwitch),
	}

	acc.AddService(acc.Switch.Service)

	return &acc
}

func (acc Switch) GetID() string {
	return acc.ID
}

func (acc Switch) GetAccessory() *accessory.Accessory {
	return acc.Accessory
}
