package service

import (
	"SaaS_buy/mydb"
	"time"

	"golang.org/x/time/rate"
)

type Service struct {
	DB *mydb.DB
	Sj mydb.SqlJSON
	l  *rate.Limiter
}

func (s *Service) InitService() {
	// 创建mysql连接
	db := mydb.DB{}
	// 初始化数据库连接和配置
	db.InitDB()
	// 传给service使用
	s.DB = &db
	s.Sj = db.Sj

	// 创建限流器
	// 每1秒投放一个令牌，桶大小10个，初始大小10个
	l := rate.NewLimiter(rate.Every(5*time.Second), 5)
	s.l = l
}
