package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"

	//_ "github.com/denisenkom/go-mssqldb"
	"net/http"
	"webApp/controller"
	"webApp/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	r := chi.NewRouter()
	r.Get("/main", controller.HandleMain)
	r.Get("/login", controller.HandleLogin)
	r.Post("/login/user", controller.HandlePOST)

	fileServer := http.FileServer(http.Dir("./view/staticFiles"))
	r.Handle("/static/", http.StripPrefix("/static", fileServer))

	err := http.ListenAndServe(":5050", r)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %s\n", err.Error())
		return
	}
}
