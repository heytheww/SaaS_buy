package model

import "time"

type ReqAddUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Role     int8   `json:"role" binding:"required,oneof=1 3 7"`
	Grade    int    `json:"grade" binding:"required"`
}

type ReqUpdateUser struct {
	Id int `json:"id" binding:"required"`
	ReqAddUser
	Update_Time time.Time `json:"update_time"`
}
