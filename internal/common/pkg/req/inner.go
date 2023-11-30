package req

type Response struct {
	Code  int         `json:"code"`  // 状态码,这个状态码是与前端和APP约定的状态码,非HTTP状态码
	Data  interface{} `json:"data"`  // 返回数据
	Msg   string      `json:"msg"`   // 自定义返回的消息内容
	Field string      `json:"field"` // 自定义返回的消息内容
}

