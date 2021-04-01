package utils

import (
	"encoding/json"
	"fmt"
	"github.com/dustin/go-coap"
)

//@function: Coap
//@description: Coap发送方法
//@param: subject string, body string
//@return: error

func CoapTest(host string, path string, datas map[string]interface{}) error {
	return send(host, path, datas)
}

//@function: send
//@description: coap发送方法
//@param: subject string, body string
//@return: error

func send(host string, path string, datas map[string]interface{}) error {
	b, err := json.Marshal(datas)
	if err != nil {
		fmt.Println("Error dialing: %v", err)
	}

	req := coap.Message{
		Type: coap.Confirmable,
		Code: coap.POST,
		//MessageID: 12345,
		Payload: []byte(b),
	}
	//req.SetOption(coap.ETag, "weetag")
	//req.SetOption(coap.MaxAge, 3)
	// 设置路径
	req.SetPathString(path)
	c, err := coap.Dial("udp", host+":5683")
	if err != nil {
		fmt.Println("Error dialing: %v", err)
	}

	c.Send(req)
	//if err != nil {
	//	return err
	//}

	return nil
}
