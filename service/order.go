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

func (s Service) AddOrder() {
	fmt.Println("订单生成模块:启动")

	// 堵塞读异步消息队列
	// 消息队列长度1000
	rdb := s.RDB
	ctx := context.Background()
	mq := s.MQ
	rdb.GetGroup(ctx, &mq, "cg1")

	for {
		xs, err := rdb.GetMsgByGroup(ctx, &mq, "c1")
		if err != nil {
			log.Fatalln(err)
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

			value, ok := v.Values[v.ID].(string)
			if ok {
				switch v.ID {
				case "user_id":
					id, err := strconv.Atoi(value)
					if err != nil {
						log.Fatalln(errors.New("user_id Atoi error"))
					}
					order.User_Id = id
				case "product_id":
					id, err := strconv.Atoi(value)
					if err != nil {
						log.Fatalln(errors.New("user_id Atoi error"))
					}
					order.Product_Id = id
				case "name":
					order.Name = value
				case "address":
					order.Address = value
				case "remarks":
					order.Remarks = value
				}
			} else {
				log.Fatalln(errors.New("type assertion not ok"))
			}
		}

		db := s.DB.DBconn

		// 尝试从数据库中插入数据
		if db != nil {
			sqlStr := s.Sj.Order.Insert
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
				fmt.Println(errors.New("get id error").Error())
			}
			fmt.Println(id)
		} else {
			fmt.Println(errors.New("nil mysql connection").Error())
		}
	}

}
