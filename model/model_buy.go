package model

type ReqBuy struct {
	User_Id    int    `json:"user_id" binding:"required"`
	Product_Id int    `json:"product_id" binding:"required"`
	Remarks    string `json:"remarks"`
}
