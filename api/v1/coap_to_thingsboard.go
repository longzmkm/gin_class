package v1

import (
	"fmt"
	"gin_class/model/request"
	"gin_class/model/response"
	"gin_class/service"
	"gin_class/utils"
	"github.com/gin-gonic/gin"
)

// @Tags Coap
// @Summary 发送Coap
// @Security application/json
// @Produce  application/json
// @Param data body request.Coap true "IP, path, 数据"
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /coap/send [post]
func Send(c *gin.Context) {
	var coap_data request.Coap
	_ = c.ShouldBindJSON(&coap_data)

	if err := utils.Verify(coap_data, utils.CoapVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := service.CoapTest(coap_data.Host, coap_data.Path, coap_data.CoapData); err != nil {
		fmt.Println(err)
		response.FailWithMessage("发送失败", c)
	} else {
		response.OkWithData("发送成功", c)
	}
}
