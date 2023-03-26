package service

import (
	"SaaS_buy/model"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func (s *Service) DelUserService(c *gin.Context) {
	db := s.DB.DBconn

	resp := model.RespDel{}
	resp.Result = model.Result{}

	req := model.Data{}
	err := c.ShouldBind(&req)
	if err != nil {
		resp.Result.Code = http.StatusBadRequest
		resp.Result.Message = "parameter error"
		c.JSON(http.StatusOK, resp)
		return
	}

	// 尝试从数据库中查询数据
	if db != nil {

		sqlStr := s.Sj.User.Delete
		err, s, r := s.DB.PrepareURDRows(sqlStr, req.Id)

		// 删除失败
		if err != nil {
			resp.Result.Code = http.StatusInternalServerError
			resp.Result.Message = "system error"
			c.JSON(http.StatusOK, resp)
			return
		}
		defer s.Close()

		var num int64
		num, err = r.RowsAffected()
		fmt.Println(num)
		// 获取删除记录的id失败
		if err != nil || num == 0 {
			resp.Result.Code = http.StatusInternalServerError
			resp.Result.Message = "delete user error"
			c.JSON(http.StatusOK, resp)
			return
		}

		resp.Result = model.Result{Code: http.StatusOK, Message: "success"}
		c.JSON(http.StatusOK, resp)
		return
	}

	// 数据库连接失败
	resp.Result.Code = http.StatusInternalServerError
	resp.Result.Message = "system error"
	c.JSON(http.StatusOK, resp)
}

func (s *Service) AddUserService(c *gin.Context) {
	db := s.DB.DBconn

	resp := model.RespAdd{}
	resp.Result = model.Result{}

	req := model.ReqAddUser{}
	err := c.ShouldBind(&req)
	if err != nil {
		resp.Result.Code = http.StatusBadRequest
		resp.Result.Message = "parameter error"
		c.JSON(http.StatusOK, resp)
		return
	}

	// 尝试从数据库中查询数据
	if db != nil {

		sqlStr := s.Sj.User.Insert
		create_time := time.Now().Format("2006-01-02 15:04:05")
		update_time := create_time
		err, s, r := s.DB.PrepareURDRows(sqlStr, req.Username, req.Password, req.Phone,
			req.Role, req.Grade, create_time, update_time)

		// 插入失败
		if err != nil {
			resp.Result.Code = http.StatusInternalServerError
			resp.Result.Message = "system error"
			c.JSON(http.StatusOK, resp)
			return
		}
		defer s.Close()

		var id int64
		id, err = r.LastInsertId()
		// 获取新插入的记录的id失败
		if err != nil {
			resp.Result.Code = http.StatusInternalServerError
			resp.Result.Message = "add user error"
			c.JSON(http.StatusOK, resp)
			return
		}

		resp.Result = model.Result{Code: http.StatusOK, Message: "success"}
		resp.Data = model.Data{Id: int(id)}
		c.JSON(http.StatusOK, resp)
		return
	}

	// 数据库连接失败
	resp.Result.Code = http.StatusInternalServerError
	resp.Result.Message = "system error"
	c.JSON(http.StatusOK, resp)
}

func (s *Service) GetserService(c *gin.Context) {

	db := s.DB.DBconn

	resp := model.RespGet{}
	resp.Result = model.Result{}

	tb := model.TableUser{}

	var userId int = 0
	id := c.Query("id")
	if id != "" {
		userId, _ = strconv.Atoi(id)
	} else {
		resp.Result.Code = http.StatusBadRequest
		resp.Result.Message = "parameter error"
		c.JSON(http.StatusOK, resp)
		return
	}

	// 尝试从数据库中查询数据
	if db != nil {
		sqlStr := s.Sj.User.Select
		err, s, r := s.DB.PrepareCRow(sqlStr, userId)

		if err != nil {
			resp.Result.Code = http.StatusInternalServerError
			resp.Result.Message = "system error"
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
