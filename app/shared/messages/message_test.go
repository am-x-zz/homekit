package messages

import (
	"log"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	msg := &MessageWrapper{
		Message: &MessageWrapper_ToDevice{
			ToDevice: &ToDevice{
				DeviceID: 1,
				Message: &ToDevice_SetSwitchState{
					SetSwitchState: &SetSwitchState{State: true},
				},
			},
		},
	}

	b, err := proto.Marshal(msg)
	assert.Nil(t, err)

	log.Println(b)
}

func Test2(t *testing.T) {
	msg := &FromDevice{
		DeviceID: 1,
		Message: &FromDevice_SwitchState{
			SwitchState: &SwitchState{
				State: false,
			},
		},
	}

	b, err := proto.Marshal(msg)
	assert.Nil(t, err)

	log.Println(b)
}
