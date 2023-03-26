package service

import (
	"SaaS_buy/mydb"
)

type Service struct {
	DB *mydb.DB
	Sj mydb.SqlJSON
}

func (s *Service) InitService() {
	// 创建mysql连接
	db := mydb.DB{}
	// 初始化数据库连接和配置
	db.InitDB()
	// 传给service使用
	s.DB = &db
	s.Sj = db.Sj
}
