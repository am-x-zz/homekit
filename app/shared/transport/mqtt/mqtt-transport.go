package mqtt

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/am-x/homekit/app/shared/messages"
	"github.com/am-x/homekit/app/shared/transport"
	"github.com/golang/protobuf/jsonpb"
	"gobot.io/x/gobot/platforms/mqtt"
)

type Transport struct {
	adaptor *mqtt.Adaptor
}

func NewTransport(adaptor *mqtt.Adaptor) *Transport {
	adaptor.SetAutoReconnect(true)

	return &Transport{
		adaptor: adaptor,
	}
}

func (t *Transport) ToAccessory(message *messages.ToAccessory) error {
	m := jsonpb.Marshaler{}
	b, err := m.MarshalToString(message)

	if err != nil {
		return err
	}

	if ok := t.adaptor.Publish("to/accessory", []byte(b)); !ok {
		return errors.New("message not published")
	}

	return nil
}

func (t *Transport) OnAccessoryMessage(h transport.AccessoryMessageHandler) error {
	token, _ := t.adaptor.OnWithQOS("to/accessory", 0, func(msg mqtt.Message) {
		var message = new(messages.ToAccessory)

		if err := jsonpb.Unmarshal(bytes.NewReader(msg.Payload()), message); err != nil {
			log.Println("unmarshal message", err)
		}

		_ = h(message)
	})

	if ok := token.Wait(); !ok && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (t *Transport) ToHub(hubID uint32, message *messages.ToHub, opts ...transport.Option) error {
	m := jsonpb.Marshaler{}
	b, err := m.MarshalToString(message)

	if err != nil {
		return err
	}

	topic := fmt.Sprintf("to/hub/%d", hubID)

	if ok := t.adaptor.Publish(topic, []byte(b)); !ok {
		return errors.New("message not published")
	}

	return nil
}

func (t *Transport) OnHubMessage(hubID uint32, h transport.HubMessageHandler, opts ...transport.Option) error {
	topic := fmt.Sprintf("to/hub/%d", hubID)

	token, err := t.adaptor.OnWithQOS(topic, 0, func(msg mqtt.Message) {
		message := new(messages.ToHub)

		if err := jsonpb.Unmarshal(bytes.NewReader(msg.Payload()), message); err != nil {
			log.Println("unmarshal message", err)
		}

		_ = h(message)
	})

	if err != nil {
		return err
	}

	if ok := token.Wait(); !ok && token.Error() != nil {
		return token.Error()
	}

	return nil
}
