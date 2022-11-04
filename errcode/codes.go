package errcode

import "fmt"

// CodeOk 错误码统一在这里注册
const (
	CodeOk = 0
)

type E struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 错误码对应msg模版
var ET = map[int]string{
	CodeOk: "ok",
}

func New(code int, params ...interface{}) *E {
	return &E{
		Code:    code,
		Message: fmt.Sprintf(ET[code], params),
	}
}

func (e E) Error() string {
	return fmt.Sprintf("%d:%s", e.Code, e.Message)
}
