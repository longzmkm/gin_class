package request

type Coap struct {
	Host     string                 `json:"host"`
	Path     string                 `json:"path"`
	Port     string                 `json:"port"'`
	CoapData map[string]interface{} `json:"datas"`
}
