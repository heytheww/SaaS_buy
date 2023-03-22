package main

import (
	"SaaS_buy/service"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 简单的路由组: v1
	v1 := router.Group("/general")
	{
		v1.POST("/buy", service.BuyService)
	}
}
