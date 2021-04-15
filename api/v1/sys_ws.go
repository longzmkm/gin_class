package v1

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gin_class/model"
	"github.com/brocaar/chirpstack-api/go/v3/common"
	"github.com/brocaar/chirpstack-api/go/v3/gw"
	"github.com/brocaar/chirpstack-simulator/simulator"
	"github.com/brocaar/lorawan"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// @Tags Websockket
// @Summary chirp stack转发
// @Security application/json
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /ws/stack [GET]
func Stack(c *gin.Context) {
	ws, _ := upgrader.Upgrade(c.Writer, c.Request, nil)
	defer ws.Close()
	//// 先读取一个数据配置设备
	_, message, _ := ws.ReadMessage()
	var cf model.Config
	//var sgw *simulator.Gateway
	err := json.Unmarshal(message, &cf)
	if err != nil {
		fmt.Println(err)
		fmt.Println("解析数据error")
	}
	fmt.Println(cf)
	fmt.Println(cf.Topic)

	u, _ := strconv.ParseUint(cf.Frequency, 10, 32)
	//
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//
	deviceMap := make(map[lorawan.EUI64]*model.DeviceState)
	// 先创建一个网关 用来收发数据
	client := model.Client()
	defer client.Disconnect(0)
	sgw, err := createGateway(&cf,client)

	if err != nil {
		panic(err)
	}
	// 创建一个设备 关联网关
	sdv, err := createDevice(ctx, &cf, sgw, ws, u)
	if err != nil {
		panic(err)
	}
	deviceState := model.DeviceState{time.Now(), sdv}
	deviceMap[cf.DevEUI] = &deviceState

	for {
		//读取ws中的数据
		_, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println("ws读取数据错误")
			break
		}
		var cf model.Config
		er := json.Unmarshal(message, &cf)
		if er != nil {
			fmt.Println("ws解析数据错误")
			continue
		}
		u, _ := strconv.ParseUint(cf.Frequency, 10, 32)
		_, found := deviceMap[cf.DevEUI]
		if !found {
			sdv, err := createDevice(ctx, &cf, sgw, ws, u)
			if err != nil {
				fmt.Println("循环中创建设备失败")
			}
			deviceState := model.DeviceState{time.Now(), sdv}
			deviceMap[cf.DevEUI] = &deviceState
		} else {
			deviceMap[cf.DevEUI].Date = time.Now()
		}

		_ = simulator.WithUplinkPayload(true, 10, []byte(cf.Data))(deviceMap[cf.DevEUI].Device)
	}
	//cancel()
}

func createGateway(cf *model.Config,client mqtt.Client) (*simulator.Gateway, error) {
	sgw, err := simulator.NewGateway(
		simulator.WithMQTTClient(client),
		//simulator.WithMQTTCredentials("mq.nlecloud.com:1883", "", ""),
		simulator.WithGatewayID(cf.GatewayID),
		simulator.WithEventTopicTemplate(cf.Topic+"/gateway/{{ .GatewayID }}/event/{{ .Event }}"),
		simulator.WithCommandTopicTemplate(cf.Topic+"/gateway/{{ .GatewayID }}/command/{{ .Command }}"),
	)
	return sgw, err
}

func createDevice(ctx context.Context, cf *model.Config, sgw *simulator.Gateway, ws *websocket.Conn, u uint64) (*simulator.Device, error) {
	var wg sync.WaitGroup
	sdv, err := simulator.NewDevice(ctx, &wg,
		simulator.WithDevEUI(cf.DevEUI),
		simulator.WithAppKey(cf.AppKey),
		simulator.WithRandomDevNonce(),
		simulator.WithUplinkInterval(time.Second*1),
		simulator.WithUplinkCount(0),
		simulator.WithUplinkPayload(true, 10, []byte(cf.Data)),
		simulator.WithUplinkTXInfo(gw.UplinkTXInfo{
			Frequency:  uint32(u),
			Modulation: common.Modulation_LORA,
			ModulationInfo: &gw.UplinkTXInfo_LoraModulationInfo{
				LoraModulationInfo: &gw.LoRaModulationInfo{
					Bandwidth:       125,
					SpreadingFactor: 7,
					CodeRate:        "3/4",
				},
			},
		}),
		simulator.WithGateways([]*simulator.Gateway{sgw}),
		simulator.WithDownlinkHandlerFunc(func(conf, ack bool, fCntDown uint32, fPort uint8, data []byte) error {

			if len(data) > 0 {
				var respdata model.RespData
				err := json.Unmarshal(data, &respdata)
				if err != nil {
					fmt.Println(err)
					fmt.Println("解析数据error")
				}
				respdata.Devid = cf.DevEUI.String()

				datas, _ := json.Marshal(respdata)

				// TODO 这边是向下发送数据
				ws.WriteMessage(1, datas)
				fmt.Printf("Recive Server data: %s\n ", datas)
			}
			return nil
		}),
	)
	return sdv, err
}