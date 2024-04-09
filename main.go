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
	r.Get("/main", controller.MainPage)
	r.Get("/login", controller.LoginPage)
	r.Get("/register", controller.RegisterPage)
	r.Get("/account", controller.AccountPage)
	r.Post("/login/user", controller.LoginPOST)
	r.Post("/register/user", controller.RegisterPOST)

	fileServer := http.FileServer(http.Dir("./view/staticFiles"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	err := http.ListenAndServe(":5500", r)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %s\n", err.Error())
		return
	}
}
