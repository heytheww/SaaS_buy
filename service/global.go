package service

import (
	"SaaS_buy/mydb"
	"SaaS_buy/util"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/time/rate"
)

type Service struct {
	MaxConn   int // mysql连接池 最大打开连接数、最大空闲连接数
	DB        *mydb.DB
	Sj        *mydb.SqlJSON // 执行查询sql
	l         *rate.Limiter
	Limit     time.Duration   // 每 Limit 时间生成一个令牌
	Bursts    int             // 桶初始大小、突发申请令牌数
	AMQP_URL  string          // 消息队列 url
	amqpConn  amqp.Connection // 消息队列的连接
	MQCh      *amqp.Channel   // 与消息队列通信的通道
	Queue     *amqp.Queue     // 队列实体
	RDBClient *mydb.RDB
}

func (s *Service) InitService() {
	// 创建mysql连接
	db := mydb.DB{}
	// 初始化数据库连接和配置
	err := db.InitDB(s.MaxConn)
	util.FailOnError(err, "mysql init error")
	// 传给service使用
	s.DB = &db
	s.Sj = db.Sj

	// 初始化redis client
	rdb := mydb.RDB{}
	err = rdb.InitRDB()
	util.FailOnError(err, "redis init error")
	// 传给service使用
	s.RDBClient = &rdb

	// 创建一个异步消息队列
	conn, err := amqp.Dial(s.AMQP_URL)
	util.FailOnError(err, "Failed to connect to RabbitMQ")
	s.amqpConn = *conn
	// defer conn.Close()
	ch, err := conn.Channel()
	util.FailOnError(err, "Failed to open a channel")
	// defer ch.Close()
	q, err2 := ch.QueueDeclare(
		"order",
		false,
		false,
		false,
		false,
		nil,
	)
	util.FailOnError(err2, "Failed to declare a queue")
	// 传给service使用
	s.MQCh = ch
	s.Queue = &q

	// 创建限流器
	// 每1秒投放一个令牌，桶大小10个，初始大小10个
	l := rate.NewLimiter(rate.Every(s.Limit), s.Bursts)
	s.l = l
}
