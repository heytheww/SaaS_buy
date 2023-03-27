package model

import "time"

type Result struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
}

type Data struct {
	Id int `json:"id" binding:"required"`
}

type Data2 struct {
	Update_Time time.Time `json:"update_time"`
}

type RespAdd struct {
	Data   any    `json:"data"`
	Result Result `json:"result"`
}

type RespDel struct {
	Data   struct{} `json:"data"`
	Result Result   `json:"result"`
}

type RespUpdate struct {
	Data   struct{} `json:"data"`
	Result Result   `json:"result"`
}

type RespGet struct {
	Data   []any  `json:"data"`
	Result Result `json:"result"`
}
