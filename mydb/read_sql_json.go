package mydb

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type CURD struct {
	Insert string `json:"insert"`
	Delete string `json:"delete"`
	Update string `json:"update"`
	Select string `json:"select"`
}

type SqlJSON struct {
	DSN     string `json:"dsn"`
	User    CURD   `json:"user"`
	Act     CURD   `json:"act"`
	Order   CURD   `json:"order"`
	Product CURD   `json:"product"`
}

func ReadSqlJson() (*SqlJSON, error) {

	pwd, _ := os.Getwd() // 获取当前所在工作目录
	f_path := filepath.Join(pwd, "conf", "sql.json")
	buf, err := os.ReadFile(f_path)
	if err != nil {
		return nil, err
	}

	s := SqlJSON{}
	err = json.Unmarshal(buf, &s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
