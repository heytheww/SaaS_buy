package model

type ReqAddProduct struct {
	Name         string `json:"name"`
	Sub_Name     string `json:"sub_name"`
	Main_Image   string `json:"main_image"`
	Detail_Image string `json:"detail_image"`
	Price        int    `json:"price"`
}
