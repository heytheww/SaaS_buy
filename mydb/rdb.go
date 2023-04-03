package mydb

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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
	MQName       string
	CusGroupName string
	firstId      string
}

func (db *RDB) InitRDB() (err error) {

	pwd, _ := os.Getwd() // 获取当前所在工作目录
	f_path := filepath.Join(pwd, "mydb", "redis.json")
	c, err := ReadRedisJson(f_path)
	if err != nil {
		return err
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: c.Address,
		// 默认没有密码
		Password: c.Password,
		// redis默认有16个数据库，命令行通过 select index 切换
		DB: c.DB_Index,
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
func (db *RDB) InitMQ(name string) (error, MQ) {
	mq := MQ{
		MQName: name,
	}
	// XAdd如果发现stream不存在则会创建
	a := db.AddMsg(context.Background(), &mq, "init", "init")
	if a.Err() != nil {
		return a.Err(), mq
	}
	mq.firstId = a.Val()

	return nil, mq
}

func (db *RDB) AddMsg(ctx context.Context, mq *MQ, values ...any) *redis.StringCmd {
	// *：消息Id自动生成
	arg := redis.XAddArgs{
		Stream:     mq.MQName,
		NoMkStream: false,
		ID:         "*",
		Values:     values,
	}
	a := db.RDBconn.XAdd(ctx, &arg)
	return a
}

// 创建消费组
// 创建成功返回“OK”
func (db *RDB) CreateGroup(ctx context.Context, mq *MQ, name string) error {
	// 该组从firstId之后的消息开始消费
	g := db.RDBconn.XGroupCreate(ctx, mq.MQName, name, mq.firstId)
	if g.Val() == "OK" {
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
		Group:    mq.CusGroupName,
		Consumer: cusName,
		Streams:  []string{mq.MQName, ">"}, // list of streams and ids, e.g. stream1 stream2 id1 id2
		Count:    1,
		Block:    0, // 堵塞读
		NoAck:    false,
	}
	xrg := db.RDBconn.XReadGroup(ctx, &arg)
	xs := xrg.Val()
	err := xrg.Err()
	return &xs, err
}

// 确认消费
func (db *RDB) ACK(ctx context.Context, mq *MQ, id string) *redis.IntCmd {
	ack := db.RDBconn.XAck(ctx, mq.MQName, mq.CusGroupName, id)
	return ack
}

// script：lua脚本
func (db *RDB) RunLua(ctx context.Context, script string, keys []string, values []interface{}) (int, error) {

	deductDock := redis.NewScript(script)

	num, err := deductDock.Run(ctx, db.RDBconn, keys, values...).Int()
	if err != nil {
		return 0, err
	}

	return num, err
}
