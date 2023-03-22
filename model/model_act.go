package model

import "time"

type ReqAddAct struct {
	Product_Id string    `json:"product_id"`
	Burst      int       `json:"burst"`
	Limt       int       `json:"limt"`
	Stock      int       `json:"stock"`
	Name       string    `json:"name"`
	Sub_Name   string    `json:"sub_name"`
	Start_Time time.Time `json:"start_time"`
	Ground     int8      `json:"ground"`
}

type RespAddAct struct {
	Data   Data   `json:"data"`
	Result Result `json:"result"`
}

type ReqDelAct struct {
	Data
}

type RespDelAct struct {
	Data   struct{} `json:"data"`
	Result Result   `json:"result"`
}

type ReqPatchAct struct {
	Id string `json:"id"`
	ReqAddAct
}

type RespPatchAct struct {
	Data   Data2  `json:"data"`
	Result Result `json:"result"`
}

type ReqGetAct struct {
	Data
}

type DataGetAct struct {
	Id string `json:"id"`
	ReqAddAct
	Create_Time time.Time `json:"create_time"`
	Update_Time time.Time `json:"update_time"`
}

type RespGetAct struct {
	Data   DataGetAct `json:"data"`
	Result Result     `json:"result"`
}
