package model

import "time"

type ReqBuy struct {
	Product_Id int    `json:"product_id"`
	User_Id    int    `json:"user_id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	Remarks    string `json:"remarks"`
}

type BuyData struct {
	ReqBuy
	Id          int       `json:"id"`
	Pay_Id      string    `json:"pay_id"`
	Status      int8      `json:"status"`
	Create_Time time.Time `json:"create_Time"`
}
