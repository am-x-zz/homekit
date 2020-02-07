package accessory

import (
	"github.com/awesome/homekit/app/shared/messages"
	"github.com/brutella/hc/accessory"
)

const (
	Manufacturer = "Apple Inc."
)

type Accessory interface {
	GetID() string
}

type WithHomeKitAccessory interface {
	Accessory
	GetAccessory() *accessory.Accessory
}

type WithAccessoryMessageProcessing interface {
	ProcessMessage(msg *messages.ToAccessory) error
}

type MessageProcessor func(msg *messages.ToAccessory) error
