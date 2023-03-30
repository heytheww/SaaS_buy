package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 限流
func (s Service) Limiter(timeout time.Duration, result any) gin.HandlerFunc {
	return func(c *gin.Context) {
		bgCtx := context.Background()
		// 超过30s就取消等待，立即返回
		ctx, _ := context.WithTimeout(bgCtx, timeout)

		err := s.l.Wait(ctx)
		if err != nil {
			log.Println("limiter wait error: ", err)
			c.JSON(http.StatusBadGateway, result)
			// 阻止执行余下的handler
			c.Abort()
			// 当前handler立即返回
			return
		}

		// 继续执行余下的handler
		c.Next()
	}
}
