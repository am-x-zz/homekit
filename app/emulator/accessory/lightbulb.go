package accessory

import (
	"github.com/brutella/hc/accessory"
	"github.com/brutella/hc/service"
)

type Lightbulb struct {
	*accessory.Accessory
	ID        string
	Lightbulb *service.Lightbulb
}

// NewLightbulb returns an light bulb accessory which one light bulb service.
func NewLightbulb(id, name string) *Lightbulb {
	info := accessory.Info{
		Name:         name,
		Manufacturer: Manufacturer,
	}

	acc := Lightbulb{
		ID:        id,
		Lightbulb: service.NewLightbulb(),
		Accessory: accessory.New(info, accessory.TypeLightbulb),
	}

	acc.Lightbulb.Brightness.SetValue(100)

	acc.AddService(acc.Lightbulb.Service)

	return &acc
}

func (acc Lightbulb) GetID() string {
	return acc.ID
}

func (acc Lightbulb) GetAccessory() *accessory.Accessory {
	return acc.Accessory
}
