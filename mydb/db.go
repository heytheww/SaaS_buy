package mydb

import (
	"SaaS_buy/util"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

type DB struct {
	DBconn *sql.DB // 数据库连接
}

// 定义一个初始化数据库的函数
func (db *DB) InitDB() (err error) {
	// DSN:Data Source Name
	pwd, _ := os.Getwd()
	f_path := filepath.Join(pwd, "mydb", "sql.json")
	j, err := util.ReadSqlJson(f_path)
	if err != nil {
		return err
	}
	dsn := j.(util.SqlJSON).Mysql
	// 不会校验账号密码是否正确
	db.DBconn, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	// 尝试与数据库建立连接（校验dsn是否正确）
	err = db.DBconn.Ping()
	if err != nil {
		return err
	}
	db.DBconn.SetConnMaxLifetime(0)
	db.DBconn.SetMaxIdleConns(50)
	db.DBconn.SetMaxOpenConns(50)
	return nil
}

func (db *DB) PrepareQueryRow(sqlStr string, query ...any) (error, *sql.Stmt, *sql.Rows) {
	stmt, err := db.DBconn.Prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("prepare failed, err:%v\n", err), nil, nil

	}

	rows, err := stmt.Query(query...)
	if err != nil {
		return fmt.Errorf("query failed, err:%v\n", err), nil, nil
	}

	return nil, stmt, rows
}
