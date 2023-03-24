package main

import (
	"SaaS_buy/service"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	sv := service.Service{}
	sv.InitService()

	// 简单的路由组: v1
	v1 := router.Group("/general")
	{
		v1.POST("/buy", sv.BuyService)
	}

	v2 := router.Group("/manage")
	{
		// v2.POST("/addUser", service.AddUserService)
		v2.GET("/getUser", sv.GetserService)
	}

	router.Run(":8080")
}
