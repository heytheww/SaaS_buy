package model

import "time"

type ReqGetOrder struct {
	Data
}

type RespGetOrder struct {
	TableOrder
}

type ReqDelOrder struct {
	Data
	Create_Time time.Time `json:"create_Time"`
	Update_Time time.Time `json:"update_Time"`
}

type RespDelOrder struct {
	Data   struct{} `json:"data"`
	Result Result   `json:"result"`
}
