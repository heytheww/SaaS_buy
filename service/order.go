package service

import (
	"SaaS_buy/model"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (s Service) AddOrder() error {
	fmt.Println("订单生成模块:启动")

	var waitErr chan error

	// 堵塞读异步消息队列
	// 消息队列长度1000
	q := s.Queue
	ch := s.MQCh
	db := s.DB.DBconn

	sqlStr := s.Sj.Order.Insert
	err, stmt := s.DB.PrepareURDRows(sqlStr)
	failOnError(err, "prepare failed, err:")
	defer stmt.Close()

	// 无限循环
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer // 消费者id自动生成
		false,  // auto-ack // 	取消自动确认，尽量少用
		false,  // exclusive
		false,  // no-local // 未支持
		false,  // no-wait  // 不等待服务确认请求，立即开始传送，如果无法消费channel就会报错
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func(waitErr chan error) {
		for d := range msgs {

			// log.Printf("Received a message: %s", d.Body)

			// now := time.Now().Format("2006-01-02 15:04:05")
			order := model.TableOrder{
				Status: 1,
			}
			// 获取消息的标记
			// tag:=d.DeliveryTag
			ids := strings.Split(string(d.Body), ":")
			i, _ := strconv.Atoi(ids[0])
			i2, _ := strconv.Atoi(ids[1])
			order.User_Id = i
			order.Product_Id = i2
			order.Remarks = ids[2]

			// 尝试从数据库中插入数据
			if db != nil {

				r, err := stmt.Exec(order.User_Id, order.Product_Id, order.Status, order.Remarks)

				if err != nil {
					waitErr <- err
				}

				// 插入失败
				if err != nil {
					fmt.Println(err)
				}

				var id int64
				id, err = r.LastInsertId()
				// 获取新插入的记录的id失败
				if err != nil {
					waitErr <- errors.New("get id error")
					// fmt.Println(errors.New("get id error").Error())
				}
				fmt.Println(id)

				// 单条应答，处理一条，应答一条
				err2 := d.Ack(false)
				if err2 != nil {
					waitErr <- errors.New("ack error")
				}

			} else {
				waitErr <- errors.New("nil mysql connection")
				// fmt.Println(errors.New("nil mysql connection").Error())
			}
		}
	}(waitErr)

	_err := <-waitErr

	return _err
}
