package mydb

import (
	"encoding/json"
	"os"
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

func ReadSqlJson(path string) (any, error) {

	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	s := SqlJSON{}
	err = json.Unmarshal(buf, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}
