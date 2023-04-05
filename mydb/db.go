package mydb

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

type DB struct {
	DBconn *sql.DB  // 数据库连接
	Sj     *SqlJSON // 数据库配置
}

// 定义一个初始化数据库的函数
func (db *DB) InitDB(maxConn int) (err error) {
	// DSN:Data Source Name
	pwd, _ := os.Getwd() // 获取当前所在工作目录
	f_path := filepath.Join(pwd, "mydb", "sql.json")
	j, err := ReadSqlJson(f_path)
	if err != nil {
		return err
	}
	db.Sj = j
	dsn := db.Sj.DSN
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
	db.DBconn.SetMaxIdleConns(maxConn)
	db.DBconn.SetMaxOpenConns(maxConn)
	return nil
}

// CURD：Insert Update Select Delete

func (db *DB) PrepareURDRows(sqlStr string) (error, *sql.Stmt) {
	stmt, err := db.DBconn.Prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("prepare failed, err:%v\n", err), nil
	}

	return nil, stmt
}

func (db *DB) PrepareURDRowsAndExec(sqlStr string, query ...any) (error, *sql.Stmt, sql.Result) {
	stmt, err := db.DBconn.Prepare(sqlStr)
	if err != nil {
		return fmt.Errorf("prepare failed, err:%v\n", err), nil, nil
	}

	res, err := stmt.Exec(query...)
	if err != nil {
		return fmt.Errorf("query failed, err:%v\n", err), nil, nil
	}

	return nil, stmt, res
}

func (db *DB) PrepareCRow(sqlStr string, query ...any) (error, *sql.Stmt, *sql.Rows) {

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
