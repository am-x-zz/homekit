package device

import (
	"log"

	"github.com/am-x/homekit/app/shared/messages"
	"github.com/pkg/errors"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type Hub struct {
	HubID      uint32
	Adaptor    gobot.Adaptor
	Devices    []HubDevice
	MessageBus chan *messages.ToAccessory
}

func NewHub(cfg *messages.HubConfig, a gobot.Adaptor, f []Factory) *Hub {
	h := &Hub{
		HubID:      cfg.GetHubID(),
		Adaptor:    a,
		Devices:    make([]HubDevice, 0),
		MessageBus: make(chan *messages.ToAccessory),
	}

	for _, cfg := range cfg.Devices {
		for _, factory := range f {
			if dev := factory(h, cfg); dev != nil {
				h.Devices = append(h.Devices, dev)
			}
		}
	}

	return h
}

func (h *Hub) GetDigitalReader() gpio.DigitalReader {
	if a, ok := h.Adaptor.(gpio.DigitalReader); ok {
		return a
	}

	return nil
}

func (h *Hub) GetDigitalWriter() gpio.DigitalWriter {
	if a, ok := h.Adaptor.(gpio.DigitalWriter); ok {
		return a
	}

	return nil
}

func (h *Hub) GetGobotDevices() (devices []gobot.Device) {
	for _, d := range h.Devices {
		devices = append(devices, d.GetGobotDevices()...)
	}
	return
}

func (h *Hub) OnMessage(message *messages.ToHub) {
	log.Println("hub on message", message)

	for _, d := range h.Devices {
		if err := d.OnMessage(message); err != nil {
			log.Println(errors.Wrap(err, "process device message"))
		}
	}
}
