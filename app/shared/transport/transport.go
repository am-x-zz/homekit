package transport

import "github.com/am-x/homekit/app/shared/messages"

type AccessoryMessageHandler func(message *messages.ToAccessory) error
type HubMessageHandler func(message *messages.ToHub) error

type Options struct {
	HubDeviceID uint32
}

type Option func(*Options)

func DefaultOptions() *Options {
	return &Options{
		HubDeviceID: 0,
	}
}

func WithHubDeviceID(id uint32) Option {
	return func(o *Options) {
		o.HubDeviceID = id
	}
}

type MessageTransport interface {
	ToAccessory(message *messages.ToAccessory) error
	OnAccessoryMessage(h AccessoryMessageHandler) error
	ToHub(hubID uint32, message *messages.ToHub, opts ...Option) error
	OnHubMessage(hubID uint32, h HubMessageHandler, opts ...Option) error
}
