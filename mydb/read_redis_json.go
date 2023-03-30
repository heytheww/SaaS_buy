package mydb

import (
	"encoding/json"
	"os"
)

type Config struct {
	Address  string `json:"address"`
	DB_Index int    `json:"db_index"`
	Password string `json:"password"`
}

func ReadRedisJson(path string) (*Config, error) {

	buf, err := os.ReadFile(path)
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
