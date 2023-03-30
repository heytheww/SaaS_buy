package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s Service) BuyService(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}
