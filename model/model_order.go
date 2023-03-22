package model

import "time"

type ReqGetOrder struct {
	Data
}

type RespGetOrder struct {
	Id          string    `json:"id"`
	User_Id     string    `json:"user_id"`
	Product_Id  string    `json:"product_id"`
	Pay_Id      string    `json:"pay_id"`
	Status      int8      `json:"status"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Address     string    `json:"address"`
	Remarks     int       `json:"remarks"`
	Create_Time time.Time `json:"create_time"`
	Update_Time time.Time `json:"update_time"`
}

type ReqDelOrder struct {
	Data
}

type RespDelOrder struct {
	Data   struct{} `json:"data"`
	Result Result   `json:"result"`
}
