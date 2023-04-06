package main

import (
	"SaaS_buy/model"
	"SaaS_buy/service"
	"SaaS_buy/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

func main() {

	conf, _err := util.ReadConfigJson()
	util.FailOnError(_err, "main:")

	router := gin.Default()
	sv := service.Service{
		Limit:    1 * time.Millisecond,
		Bursts:   conf.Burst,
		AMQP_URL: conf.AMQP_URL,
		MaxConn:  conf.MaxConn,
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

	server01 := &http.Server{
		Addr:         ":1234",
		Handler:      router.Handler(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		// 启动订单生成模块
		return sv.AddOrder()
	})

	g.Go(func() error {
		return server01.ListenAndServe()
	})

	if err := g.Wait(); err != nil {
		util.FailOnError(err, "main:")
	}

	// router.Run(":8080")
}
