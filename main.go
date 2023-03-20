package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type MyMux struct {
	Name string `json:"name"`
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	resp, err := http.Get("http://127.0.0.1:9090")
	if err != nil {
		fmt.Println(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	w.Write([]byte(body))
}

func main() {
	mux := &MyMux{Name: "liww"}
	http.ListenAndServe(":9091", mux)
}
