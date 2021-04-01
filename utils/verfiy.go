package utils

var (
	LoginVerify = Rules{"CaptchaId": {NotEmpty()}, "Captcha": {NotEmpty()}, "Username": {NotEmpty()}, "Password": {NotEmpty()}}
	CoapVerify  = Rules{"host": {NotEmpty()}, "path": {NotEmpty()}, "datas": {NotEmpty()}}
)
