package serializer

// 响应数据格式
type Response struct {
	Status uint        `json:"status" example:"1"`
	Msg    string      `json:"msg" example:"msg"`
	Data   interface{} `json:"data" example:"数据"`
	Error  string      `json:"error" example:"错误"`
}


