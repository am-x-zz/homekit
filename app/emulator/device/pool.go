package device

import (
	"github.com/awesome/homekit/app/emulator/accessory"
	"github.com/brutella/hc"
	"github.com/brutella/hc/log"
	"strconv"
	"sync"
)

const (
	StartingPort = 12340
	DefaultPin   = "12344321"
	StoragePath  = "./tmp/storage/"
)

type device struct {
	ID        string
	Accessory accessory.Accessory
	transport hc.Transport
}

type devicePool struct {
	deviceCount int
	wg          *sync.WaitGroup
	devices     map[string]*device
	stopped     chan struct{}
}

func NewDevicePool() *devicePool {
	return &devicePool{
		wg:      &sync.WaitGroup{},
		devices: make(map[string]*device),
		stopped: make(chan struct{}),
	}
}

func (rcv *devicePool) GetDevice(id string) *device {
	if d, ok := rcv.devices[id]; ok && d != nil {
		return d
	}

	return nil
}

func (rcv *devicePool) Add(hks accessory.WithHomeKitAccessory) {
	id := hks.GetID()
	config := hc.Config{Pin: DefaultPin, Port: strconv.Itoa(StartingPort + rcv.deviceCount), StoragePath: StoragePath + "/" + id}
	transport, err := hc.NewIPTransport(config, hks.GetAccessory())

	if err != nil {
		log.Info.Panic(err) // TODO: remove panic
	}

	device := &device{
		ID:        id,
		Accessory: hks,
		transport: transport,
	}

	rcv.deviceCount++
	rcv.devices[device.ID] = device
}

func (rcv *devicePool) Start() {
	for _, d := range rcv.devices {
		go d.transport.Start()
		rcv.wg.Add(1)
	}

	rcv.wg.Wait()
	rcv.stopped <- struct{}{}
}

func (rcv *devicePool) Stop() <-chan struct{} {
	for _, d := range rcv.devices {
		go func(d *device, wg *sync.WaitGroup) {
			<-d.transport.Stop()
			wg.Done()
		}(d, rcv.wg)
	}

	return rcv.stopped
}
