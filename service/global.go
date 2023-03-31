package service

import (
	"SaaS_buy/mydb"
	"log"
	"time"

	"golang.org/x/time/rate"
)

type Service struct {
	DB     *mydb.DB
	RDB    *mydb.RDB
	Sj     *mydb.SqlJSON // 执行查询sql
	l      *rate.Limiter
	Limit  time.Duration // 每 Limit 时间生成一个令牌
	Bursts int           // 桶初始大小、突发申请令牌数
	MQ     mydb.MQ
}

func (s *Service) InitService() {
	// 创建mysql连接
	db := mydb.DB{}
	// 初始化数据库连接和配置
	db.InitDB()
	// 传给service使用
	s.DB = &db
	s.Sj = db.Sj

	// 创建redis连接
	rdb := mydb.RDB{}
	// 初始化redis数据库连接和配置
	err := rdb.InitRDB()
	if err != nil {
		log.Fatal(err)
	}
	// 传给service使用
	s.RDB = &rdb

	// 创建一个异步消息队列
	mq := rdb.InitMQ("mq")
	s.MQ = mq

	// 创建限流器
	// 每1秒投放一个令牌，桶大小10个，初始大小10个
	l := rate.NewLimiter(rate.Every(s.Limit), s.Bursts)

	s.l = l
}
