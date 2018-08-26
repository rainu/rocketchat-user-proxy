package messages

import (
	"strconv"
	"time"
)

type GeneralResponse struct {
	Id  string `json:"id"`
	Msg string `json:"msg"`
}

type MethodCall struct {
	Id     string        `json:"id"`
	Msg    string        `json:"msg"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

func NewMethodCall(method string, params []interface{}) *MethodCall {
	return &MethodCall{
		Id:     strconv.Itoa(time.Now().Nanosecond()),
		Msg:    "method",
		Method: method,
		Params: params,
	}
}

type Connect struct {
	Msg     string   `json:"msg"`
	Version string   `json:"version"`
	Support []string `json:"support"`
}

func NewConnect() *Connect {
	return &Connect{
		Msg:     "connect",
		Version: "1",
		Support: []string{"1"},
	}
}

type PingResponse struct {
	Msg string `json:"msg"`
}

func NewPingResponse() *PingResponse {
	return &PingResponse{
		Msg: "ping",
	}
}
