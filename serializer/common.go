package serializer

import (
	"Fire/pkg/e"
	"net/http"
)

type Response struct {
	Status int         `json:"code"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

type TokenData struct {
	User  interface{} `json:"user"`
	Token string      `json:"token"`
}

type DataList struct {
	Item  interface{} `json:"item"`
	Total int         `json:"total"`
}

func BuildListResponse(items interface{}, total int) Response {
	return Response{
		Status: http.StatusOK,
		Data: DataList{
			Item:  items,
			Total: total,
		},
		Msg: e.GetMsg(http.StatusOK),
	}
}

// HandleError 通用错误处理方法
func HandleError(code int) Response {
	return Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   nil,
	}
}
