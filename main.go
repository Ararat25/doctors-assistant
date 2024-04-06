package main

import (
	"fmt"
	"net/http"
	"webApp/controller"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/main", controller.HandleMain)
	mux.HandleFunc("/login", controller.HandleLogin)
	mux.HandleFunc("/login/user", controller.HandlePOST)

	fileServer := http.FileServer(http.Dir("./view/staticFiles"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	err := http.ListenAndServe(":5050", mux)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %s\n", err.Error())
		return
	}
}
