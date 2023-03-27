package main

import (
	"SaaS_buy/service"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	sv := service.Service{}
	// 服务初始化
	sv.InitService()
	// 简单的路由组: v1
	v1 := router.Group("/general")
	{
		v1.POST("/buy", sv.BuyService)
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
