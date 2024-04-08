package controller

import (
	"fmt"
	"net/http"
)

func HandleMain(res http.ResponseWriter, req *http.Request) {
	s := fmt.Sprintf("Method: %s\nHost: %s\nPath: %s",
		req.Method, req.Host, req.URL.Path)
	_, err := res.Write([]byte(s))
	if err != nil {
		fmt.Printf("Ошибка записи ответа: %s\n", err.Error())
		return
	}
}
