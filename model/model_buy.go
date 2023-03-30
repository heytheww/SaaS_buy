package model

import "time"

type ReqBuy struct {
	Product_Id int    `json:"product_id" binding:"required"`
	User_Id    int    `json:"user_id" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Address    string `json:"address" binding:"required"`
	Remarks    string `json:"remarks"`
}

type BuyData struct {
	ReqBuy
	Id          int       `json:"id"`
	Pay_Id      string    `json:"pay_id"`
	Status      int8      `json:"status"`
	Create_Time time.Time `json:"create_Time"`
}
