package model

import (
	simulator "github.com/brocaar/chirpstack-simulator/simulator"
	"github.com/brocaar/lorawan"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
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

type RespData struct {
	Channel int `json:"channel"`
	Status int `json:"status"`
	Devid string `json:"devid"`
}
// Setup configures the NS MQTT gateway backend.
func Client() mqtt.Client {
	var mqttClient mqtt.Client
	mq_addr := os.Getenv("mq")
	if mq_addr == "" {
		mq_addr = "mq.nlecloud.com:1883"
	}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(mq_addr)
	opts.SetUsername("")
	opts.SetPassword("")
	opts.SetCleanSession(true)
	opts.SetAutoReconnect(true)
	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil
	}

	return mqttClient
}