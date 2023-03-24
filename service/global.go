package service

import (
	"SaaS_buy/mydb"
)

type Service struct {
	DB *mydb.DB
}

func (s *Service) InitService() {
	// 创建mysql连接
	db := mydb.DB{}
	// 初始化
	db.InitDB()
	// 传给service使用
	s.DB = &db
}
