package model

import (
	simulator "github.com/brocaar/chirpstack-simulator/simulator"
	"github.com/brocaar/lorawan"
	"time"
)

type Config struct {
	GatewayID lorawan.EUI64
	DevEUI    lorawan.EUI64
	AppKey    lorawan.AES128Key
	//Data      map[string]interface{}
	Data      string
	Topic     string
	Frequency string
}

type DeviceState struct {
	Date time.Time
	*simulator.Device
}