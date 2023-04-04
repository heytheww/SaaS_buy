package util

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	Burst     int `json:"burst"`
	MsgMaxLen int `json:"msg_max_len"`
}

func ReadConfigJson() (*Config, error) {
	pwd, _ := os.Getwd() // 获取当前所在工作目录
	f_path := filepath.Join(pwd, "util", "config.json")

	buf, err := os.ReadFile(f_path)
	if err != nil {
		return nil, err
	}

	c := Config{}
	err2 := json.Unmarshal(buf, &c)
	if err2 != nil {
		return nil, err2
	}

	return &c, nil
}
