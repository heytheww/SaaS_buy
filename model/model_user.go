package model

import "time"

type ReqAddUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     int8   `json:"role"`
	Grade    int    `json:"grade"`
}

type RespAddUser struct {
	Data   Data   `json:"data"`
	Result Result `json:"result"`
}

type ReqDelUser struct {
	Data
}

type RespqDelUser struct {
	Data   struct{} `json:"data"`
	Result Result   `json:"result"`
}

type ReqPatchUser struct {
	Id string `json:"id"`
	ReqAddUser
}

type RespPatchUser struct {
	Data   Data2  `json:"data"`
	Result Result `json:"result"`
}

type ReqGetUser struct {
	Data
}

type DataGetUser struct {
	Id string `json:"id"`
	ReqAddUser
	Create_Time time.Time `json:"create_time"`
	Update_Time time.Time `json:"update_time"`
}

type RespGetUser struct {
	Data   DataGetUser `json:"data"`
	Result Result      `json:"result"`
}
