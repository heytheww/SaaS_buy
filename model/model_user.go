package model

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
	Id int `json:"id"`
	ReqAddUser
}

type RespPatchUser struct {
	Data   Data2  `json:"data"`
	Result Result `json:"result"`
}

type ReqGetUser struct {
	Data
}
