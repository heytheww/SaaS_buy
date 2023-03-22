package model

import "time"

type ReqBuy struct {
	Product_Id string `json:"product_id"`
	User_Id    string `json:"user_id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	Remarks    string `json:"remarks"`
}

type BuyData struct {
	ReqBuy
	Id          string    `json:"id"`
	Pay_Id      string    `json:"pay_id"`
	Status      int8      `json:"status"`
	Create_Time time.Time `json:"create_Time"`
}

type RespBuy struct {
	Data   BuyData `json:"data"`
	Result Result  `json:"result"`
}
