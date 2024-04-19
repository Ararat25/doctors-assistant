package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"webApp/controller/account"
	"webApp/controller/basic"
	"webApp/controller/login"
	"webApp/controller/register"
	"webApp/initializers"
	"webApp/model"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authService := model.NewAuthService([]byte(os.Getenv("AUTH_SALT")), []byte(os.Getenv("TOKEN_SALT")))

	loginHandler := login.NewHandler(authService)

	r.Get("/main", basic.MainPage)
	r.Get("/login", login.LoginPage)
	r.Get("/register", register.RegisterPage)
	r.Post("/account", account.AccountPage)

	//r.Post("/login/user", controller.LoginPOST)
	r.Method(http.MethodPost, "/login/user", loginHandler)

	r.Post("/register/user", register.RegisterPOST)

	fileServer := http.FileServer(http.Dir("./view/staticFiles"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	err := http.ListenAndServe(":5500", r)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %s\n", err.Error())
		return
	}
}
