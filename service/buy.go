package service

import (
	"SaaS_buy/model"
	"SaaS_buy/util"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (s Service) BuyService(c *gin.Context) {

	resp := model.RespBuy{}
	resp.Result = model.Result{}

	req := model.ReqBuy{}
	err := c.ShouldBind(&req)
	if err != nil {
		resp.Result.Code = http.StatusBadRequest
		resp.Result.Message = "parameter error"
		c.JSON(http.StatusOK, resp)
		return
	}

	rdb := s.RDBClient

	// 1.执行布隆过滤器，判断商品是否存在
	// TODO

	// 2.执行redis lua，扣减库存
	// 读取lua脚本
	pwd, _ := os.Getwd() // 获取当前所在工作目录
	f_path := filepath.Join(pwd, "conf", "stock.lua")
	buf, err2 := os.ReadFile(f_path)
	if err2 != nil {
		log.Fatalln(err2)
	}
	// 准备执行脚本的参数
	pId := strconv.Itoa(req.Product_Id)
	keys := []string{"stock", pId}
	values := []interface{}{}
	num, err3 := rdb.RunLua(context.Background(), string(buf), keys, values)
	if err3 != nil {
		fmt.Println(err3)
		resp.Result.Code = http.StatusBadGateway
		resp.Result.Message = "请求失败，请重试"
		c.JSON(http.StatusBadGateway, resp)
		return
	}

	switch num {
	case -1: // 库存不限
		resp.Result.Code = http.StatusOK
		resp.Result.Message = "恭喜您，抢购成功！"
	case -2: // 商品不存在
		resp.Result.Code = http.StatusNotFound
		resp.Result.Message = "商品不存在"
		c.JSON(http.StatusOK, resp)
		return
	case 0: // 库存不足
		resp.Result.Code = http.StatusForbidden
		resp.Result.Message = "库存不足"
		c.JSON(http.StatusOK, resp)
		return
	default: // 扣前库存还剩：
		// fmt.Println("扣前库存还剩：", num)
		resp.Result.Code = http.StatusOK
		resp.Result.Message = "恭喜您，抢购成功！"
	}
	c.JSON(http.StatusOK, resp)

	// 3.向异步消息队列推送 订单生成源信息
	// 按照用户id:商品id组装成消息
	body := strconv.Itoa(req.User_Id) + ":" + strconv.Itoa(req.Product_Id) + ":" + req.Remarks

	err = s.MQCh.PublishWithContext(context.Background(),
		"",
		s.Queue.Name,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	util.FailOnError(err, "Failed to publish a message")
	log.Printf("msg have sent：%s\n", body)
}
