package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"webApp/controller/account"
	"webApp/controller/assistant"
	"webApp/controller/basic"
	"webApp/controller/login"
	"webApp/controller/logout"
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

	authService := model.NewAuthService([]byte(os.Getenv("AUTH_SALT")), []byte(os.Getenv("TOKEN_SALT")), initializers.DB)
	authMiddleware := myMiddleware.NewAuthMiddleware(authService)

	loginHandler := login.NewHandler(authService)
	logoutHandler := logout.NewHandler(authService)
	registerHandler := register.NewHandler(authService)
	refreshTokenHandler := refreshToken.NewHandler(authService)
	assistantHandler := assistant.NewHandler(authService)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.CheckToken)
		r.Get("/account/user", account.Page)
		r.Post("/account/assistant/message", assistantHandler.ServeHTTP)
	})

	r.Post("/login/user", loginHandler.ServeHTTP)
	r.Post("/register/user", registerHandler.ServeHTTP)
	r.Get("/refresh-token", refreshTokenHandler.ServeHTTP)
	r.Get("/logout", logoutHandler.ServeHTTP)
	r.Get("/main", basic.Page)
	r.Get("/login", login.Page)
	r.Get("/register", register.Page)

	fileServer := http.FileServer(http.Dir("./view/staticFiles"))
	r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")), r)
	if err != nil {
		fmt.Printf("Ошибка запуска сервера: %s\n", err.Error())
		return
	}
}
