package model

import "time"

// 表的model

type TableUser struct {
	Id          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Phone       string    `json:"phone"`
	Role        int8      `json:"role"`
	Grade       int       `json:"grade"`
	Del_Flag    int8      `json:"del_flag"`
	Create_Time time.Time `json:"create_time"`
	Update_Time time.Time `json:"update_time"`
}

type TableProduct struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Sub_Name     string `json:"sub_name"`
	Main_Image   string `json:"main_image"`
	Detail_Image string `json:"detail_image"`
	Price        int    `json:"price"`
	Del_Flag     int8   `json:"del_flag"`
	Create_Time  string `json:"create_time"`
	Update_Time  string `json:"update_time"`
}

type TableOrder struct {
	Id          int    `json:"id"`
	User_Id     int    `json:"user_id"`
	Product_Id  int    `json:"product_id"`
	Pay_Id      string `json:"pay_id"`
	Status      int8   `json:"status"`
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Remarks     int    `json:"remarks"`
	Del_Flag    int8   `json:"del_flag"`
	Create_Time string `json:"create_time"`
	Update_Time string `json:"update_time"`
}

type TableActivities struct {
	Id          int    `json:"id"`
	Product_Id  int    `json:"product_id"`
	Burst       int    `json:"burst"`
	Limt        int    `json:"limt"`
	Stock       int    `json:"stock"`
	Name        string `json:"name"`
	Sub_Name    string `json:"sub_name"`
	Start_Time  string `json:"start_time"`
	Ground      int8   `json:"ground"`
	Del_Flag    int8   `json:"del_flag"`
	Create_Time string `json:"create_time"`
	Update_Time string `json:"update_time"`
}
