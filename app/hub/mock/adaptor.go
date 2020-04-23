package mock

import (
	"github.com/go-redis/redis/v7"
	"log"
	"sync"
)

type Adaptor struct {
	name        string
	pins        sync.Map
	redisClient *redis.Client
}

func NewAdaptor() *Adaptor {
	r := &Adaptor{
		pins: sync.Map{},
	}

	return r
}

func (a *Adaptor) Name() string {
	return a.name
}

func (a *Adaptor) SetName(n string) {
	a.name = n
}

func (a *Adaptor) Connect() error {
	a.redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6320",
		DB:   0,
	})

	//var interval = 5 * time.Second
	//
	//go func() {
	//	for {
	//		p := "13"
	//		v, err := a.DigitalRead(p)
	//
	//		if err == nil {
	//			var val byte = 0
	//
	//			if v == 0 {
	//				val = 1
	//			}
	//
	//			_ = a.DigitalWrite(p, val)
	//		}
	//
	//		select {
	//		case <-time.After(interval):
	//			//case <-b.halt:
	//			//	return
	//		}
	//	}
	//}()

	log.Println("Redis mock connected...")

	return nil
}

func (a *Adaptor) Finalize() error {
	return nil
}

func (a *Adaptor) DigitalRead(p string) (val int, err error) {
	val, err = a.redisClient.HGet("raspi", p).Int()

	if err == redis.Nil {
		err = nil
	}

	return
}

func (a *Adaptor) DigitalWrite(p string, v byte) (err error) {
	a.redisClient.HSet("raspi", p, int(v))
	return
}
