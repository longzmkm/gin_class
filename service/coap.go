package service

import (
	"gin_class/utils"
)

//@function: EmailTest
//@description: 发送邮件测试
//@return: err error

func CoapTest(host string, path string, datas map[string]interface{}) (err error) {
	err = utils.CoapTest(host, path, datas)
	return err
}
