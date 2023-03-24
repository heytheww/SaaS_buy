package model

import "time"

type ReqAddProduct struct {
	Name         string `json:"name"`
	Sub_Name     string `json:"sub_name"`
	Main_Image   string `json:"main_image"`
	Detail_Image string `json:"detail_image"`
	Price        int    `json:"price"`
}

type RespAddProduct struct {
	Data   Data   `json:"data"`
	Result Result `json:"result"`
}

type ReqDelProduct struct {
	Data
}

type RespqDelProduct struct {
	Data   struct{} `json:"data"`
	Result Result   `json:"result"`
}

type ReqPatchProduct struct {
	Id int `json:"id"`
	ReqAddProduct
}

type RespPatchProduct struct {
	Data   Data2  `json:"data"`
	Result Result `json:"result"`
}

type ReqGetProduct struct {
	Data
}

type DataGetProduct struct {
	TableProduct
	Create_Time time.Time `json:"create_Time"`
	Update_Time time.Time `json:"update_Time"`
}

type RespGetProduct struct {
	Data   DataGetProduct `json:"data"`
	Result Result         `json:"result"`
}
