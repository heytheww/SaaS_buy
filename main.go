package main

import (
	"SaaS_buy/model"
	"SaaS_buy/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	sv := service.Service{
		Limit:  1 * time.Second,
		Bursts: 10,
	}
	// 限流器响应
	r := model.Result{
		Code:    http.StatusBadGateway,
		Message: "人太多啦，请稍后重试！",
	}
	resp := model.RespBuy{
		Data:   model.NilData{},
		Result: r,
	}

	// 服务初始化
	sv.InitService()

	// 启动订单生成模块
	go sv.AddOrder()

	// 简单的路由组: v1
	v1 := router.Group("/general")
	{
		v1.POST("/buy", sv.Limiter(5*time.Second, resp), sv.BuyService)
	}

	v2 := router.Group("/manage")
	{
		v2.GET("/getUser", sv.GetserService)
		v2.POST("/addUser", sv.AddUserService)
		v2.DELETE("/delUser", sv.DelUserService)
		v2.PUT("/updateUser", sv.UpdateUserService)
	}

	router.Run(":8080")
}
