package util

import (
	"encoding/json"
	"os"
)

type SqlJSON struct {
	Mysql string `json:"mysql"`
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
