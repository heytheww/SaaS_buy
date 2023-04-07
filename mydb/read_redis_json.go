package mydb

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Address  string `json:"address"`
	DB_Index int    `json:"db_index"`
	Password string `json:"password"`
}

func ReadRedisJson() (*Config, error) {

	pwd, _ := os.Getwd() // 获取当前所在工作目录
	f_path := filepath.Join(pwd, "mydb", "redis.json")

	buf, err := os.ReadFile(f_path)
	if err != nil {
		return nil, err
	}

	c := Config{}
	err = json.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
