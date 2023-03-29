package mydb

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RDB struct {
	RDBconn  *redis.Client // redis数据库连接
	Addr     string
	Password string
	DBIndex  int
}

type MQ struct {
	MQName        string
	MaxLen        int64
	ConsumerGroup string
}

func (db *RDB) InitRDB() (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: db.Addr,
		// 默认没有密码
		Password: db.Password,
		// redis默认有16个数据库，命令行通过 select index 切换
		DB: db.DBIndex,
	})
	if rdb == nil {
		return errors.New("nil connection")
	}
	db.RDBconn = rdb
	return nil
}

func (db *RDB) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	set := db.RDBconn.Set(ctx, key, value, expiration)
	if set.Val() == "OK" {
		return nil
	}
	return set.Err()
}

func (db *RDB) Get(ctx context.Context, key string) (string, error) {
	val, err := db.RDBconn.Get(ctx, key).Result()
	var err2 error
	switch {
	case err == redis.Nil:
		err2 = errors.New("key不存在")
	case err != nil:
		err2 = errors.New("错误：" + err.Error())
	case val == "":
		err2 = errors.New("值是空字符串")
	}
	if err2 != nil {
		return "", err2
	}
	return val, nil
}

// 使用redis.Stream实现message queue
// Stream不需要事先创建，redis会在xadd时自动创建
func (db *RDB) InitMQ(name string, maxLen int64) MQ {
	mq := MQ{
		MQName: name,
		MaxLen: maxLen,
	}
	return mq
}

func (db *RDB) AddMsg(ctx context.Context, mq *MQ, key string, value string) {
	// *：消息Id自动生成
	val := []string{key, value}
	arg := redis.XAddArgs{
		Stream:     mq.MQName,
		NoMkStream: false,
		MaxLen:     mq.MaxLen,
		// Approx causes MaxLen and MinID to use "~" matcher (instead of "=").
		Approx: true,
		ID:     "*",
		Values: val,
	}
	db.RDBconn.XAdd(ctx, &arg)
}

// 创建消费组
// 创建成功返回“OK”
func (db *RDB) GetGroup(ctx context.Context, mq *MQ, name string) error {
	g := db.RDBconn.XGroupCreate(ctx, mq.MQName, name, "0-0")
	if g.Val() == "OK" {
		mq.ConsumerGroup = name
		return nil
	}

	return g.Err()
}

// 销毁消费组
// 删除成功返回 1
func (db *RDB) DestroyGroup(ctx context.Context, mq *MQ, name string) error {
	ic := db.RDBconn.XGroupDestroy(ctx, mq.MQName, name)
	if ic.Val() == 1 {
		return nil
	}
	return ic.Err()
}

// 消费组消费
func (db *RDB) GetMsgByGroup(ctx context.Context, mq *MQ, cusName string) (*[]redis.XStream, error) {
	arg := redis.XReadGroupArgs{
		Group:    mq.ConsumerGroup,
		Consumer: cusName,
		Streams:  []string{mq.MQName, ">"}, // list of streams and ids, e.g. stream1 stream2 id1 id2
		Count:    1,
		Block:    time.Millisecond * 5000,
		NoAck:    false,
	}
	xrg := db.RDBconn.XReadGroup(ctx, &arg)
	xs := xrg.Val()
	err := xrg.Err()
	return &xs, err
}
