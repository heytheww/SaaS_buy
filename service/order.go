package service

import (
	"SaaS_buy/model"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

func (s Service) AddOrder() error {
	fmt.Println("订单生成模块:启动")

	// 堵塞读异步消息队列
	// 消息队列长度1000
	rdb := s.RDB
	ctx := context.Background()
	mq := s.MQ

	// 消息id
	msgId := ""

	// 无限循环
	for {
		xs, err := rdb.GetMsgByGroup(ctx, &mq, "c1")
		if err != nil {
			return err
		}
		now := time.Now().Format("2006-01-02 15:04:05")
		order := model.TableOrder{
			Status:      1,
			Create_Time: now,
			Update_Time: now,
		}

		// 目前只用到一条队列
		// 组装数据表参数
		messages := (*xs)[0].Messages
		for _, v := range messages {
			msgId = v.ID
			for key, value := range v.Values {
				value, ok := value.(string)
				if ok {
					switch key {
					case "user_id":
						id, err := strconv.Atoi(value)
						if err != nil {
							return errors.New("user_id Atoi error")
							// log.Fatalln(errors.New("user_id Atoi error"))
						}
						order.User_Id = id
					case "product_id":
						id, err := strconv.Atoi(value)
						if err != nil {
							return errors.New("product_id Atoi error")
							// log.Fatalln(errors.New("product_id Atoi error"))
						}
						order.Product_Id = id
					case "name":
						order.Name = value
					case "phone":
						order.Phone = value
					case "address":
						order.Address = value
					case "remarks":
						order.Remarks = value
					}
				} else {
					return errors.New("type assertion not ok")
					// log.Fatalln(errors.New("type assertion not ok"))
				}
			}
		}

		db := s.DB.DBconn

		// 尝试从数据库中插入数据
		if db != nil {
			sqlStr := s.Sj.Order.Insert
			// time.Sleep(2 * time.Second) // 测试用
			err, s, r := s.DB.PrepareURDRows(sqlStr, order.User_Id, order.Product_Id, order.Status,
				order.Name, order.Phone, order.Address, order.Remarks, order.Create_Time, order.Update_Time)

			// 插入失败
			if err != nil {
				log.Fatalln(err)
			}
			defer s.Close()

			var id int64
			id, err = r.LastInsertId()
			// 获取新插入的记录的id失败
			if err != nil {
				return errors.New("get id error")
				// fmt.Println(errors.New("get id error").Error())
			}

			// 消息ack这块其实不是最重要的，可以简单处理
			// 因为读消息时，读取的是未读过的
			ack := rdb.ACK(ctx, &mq, msgId)
			if ack.Err() == nil {
				fmt.Println(id, ack.Err())
			}

		} else {
			return errors.New("nil mysql connection")
			// fmt.Println(errors.New("nil mysql connection").Error())
		}
	}
}
