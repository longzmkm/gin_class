package request

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Captcha  string `json:"captcha"`
	CaptchaId string `json:"captcha_id"`
}
