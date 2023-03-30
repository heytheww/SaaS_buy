package service

import (
	"SaaS_buy/model"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
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

	// 1.执行布隆过滤器，判断商品是否存在
	// 2.执行redis lua，扣减库存

	// 读取lua脚本
	pwd, _ := os.Getwd() // 获取当前所在工作目录
	f_path := filepath.Join(pwd, "mydb", "stock.lua")
	buf, err2 := os.ReadFile(f_path)
	if err2 != nil {
		log.Fatalln(err2)
	}
	// 准备执行脚本的参数
	pId := strconv.Itoa(req.Product_Id)
	keys := []string{"stock", pId}
	values := []interface{}{}
	num, err3 := s.RDB.RunLua(c.Request.Context(), string(buf), keys, values)
	if err3 != nil {
		log.Fatalln(err3)
	}

	switch num {
	case -1: // 库存不限
		resp.Result.Code = http.StatusOK
		resp.Result.Message = "恭喜您，抢购成功！"
	case -2: // 商品不存在
		resp.Result.Code = http.StatusNotFound
		resp.Result.Message = "商品不存在"
	case 0: // 库存不足
		resp.Result.Code = http.StatusForbidden
		resp.Result.Message = "库存不足"
	default: // 扣前库存还剩：
		// fmt.Println("扣前库存还剩：", num)
		resp.Result.Code = http.StatusOK
		resp.Result.Message = "恭喜您，抢购成功！"
	}

	c.JSON(http.StatusOK, resp)

	// 3.向异步消息队列推送 订单生成源信息

}
