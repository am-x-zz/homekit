package hub

import (
	"github.com/am-x/homekit/app/shared/transport"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

type Hub struct {
	HubID     uint32
	Transport transport.MessageTransport
	Adaptor   gobot.Adaptor
}

func (h *Hub) GetDigitalReader() gpio.DigitalReader {
	if a, ok := h.Adaptor.(gpio.DigitalReader); ok {
		return a
	}

	return nil
}
