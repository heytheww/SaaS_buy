package service

import (
	"SaaS_buy/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func (s *Service) GetserService(c *gin.Context) {

	db := s.DB.DBconn

	resp := model.Resp{}
	resp.Result = model.Result{}

	tb := model.TableUser{}

	var userId int = 0
	id := c.Query("id")
	if id != "" {
		userId, _ = strconv.Atoi(id)
	} else {
		resp.Result.Code = http.StatusNotFound
		resp.Result.Message = "missing parameter"
		c.JSON(http.StatusOK, resp)
		return
	}

	// 尝试从数据库中查询数据
	if db != nil {
		sqlStr := `SELECT id,username,password,phone,role,grade,create_time,update_time FROM user WHERE id>=?`

		err, s, r := s.DB.PrepareQueryRow(sqlStr, userId)

		if err != nil {
			resp.Result.Code = http.StatusInternalServerError
			c.JSON(http.StatusOK, resp)
			return
		}
		defer s.Close()
		defer r.Close()

		slices := make([]any, 0)
		for r.Next() {
			err = r.Scan(&tb.Id, &tb.Username, &tb.Password, &tb.Phone, &tb.Role, &tb.Grade, &tb.Create_Time, &tb.Update_Time)
			if err != nil {
				break
			}
			slices = append(slices, tb)
		}

		resp.Result = model.Result{Code: http.StatusOK, Message: "success"}
		resp.Data = slices
		c.JSON(http.StatusOK, resp)
		return
	}

	resp.Result.Code = http.StatusInternalServerError
	resp.Result.Message = "system error"
	c.JSON(http.StatusOK, resp)
}
