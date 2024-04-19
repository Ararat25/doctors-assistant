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
	"webApp/controller/refreshToken"
	"webApp/controller/register"
	"webApp/initializers"
	"webApp/model"
	myMiddleware "webApp/model/middleware"
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
	authMiddleware := myMiddleware.NewAuthMiddleware(authService)

	loginHandler := login.NewHandler(authService)
	registerHandler := register.NewHandler(authService)
	refreshTokenHandler := refreshToken.NewHandler(authService)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.CheckToken)
		r.Get("/account", account.Page)
	})

	r.Method(http.MethodPost, "/login/user", loginHandler)
	r.Method(http.MethodPost, "/register/user", registerHandler)
	r.Method(http.MethodGet, "/refresh-token", refreshTokenHandler)

	r.Get("/main", basic.Page)
	r.Get("/login", login.Page)
	r.Get("/register", register.Page)

	fileServer := http.FileServer(http.Dir("./view/staticFiles"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	err := http.ListenAndServe(":5500", r)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %s\n", err.Error())
		return
	}
}
