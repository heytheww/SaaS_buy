package util

import (
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		// log.Println("%s: %v", msg, err)
		log.Fatalln("%s: %v", msg, err)
	}
}
