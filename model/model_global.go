package model

import "time"

type Result struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
}

type Data struct {
	Id string `json:"id"`
}

type Data2 struct {
	Update_Time time.Time `json:"update_time"`
}
