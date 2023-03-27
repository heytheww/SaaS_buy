package model

import "time"

type ReqAddAct struct {
	Product_Id int       `json:"product_id"`
	Burst      int       `json:"burst"`
	Limt       int       `json:"limt"`
	Stock      int       `json:"stock"`
	Name       string    `json:"name"`
	Sub_Name   string    `json:"sub_name"`
	Start_Time time.Time `json:"start_time"`
	Ground     int8      `json:"ground"`
}
