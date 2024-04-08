package controller

import (
	"fmt"
	"net/http"
	"webApp/initializers"
	"webApp/model"
)

func HandlePOST(res http.ResponseWriter, req *http.Request) {
	fmt.Printf("Email: %s\nPassword: %s\n", req.PostFormValue("email"), req.PostFormValue("password"))

	var users = []model.User{}

	_ = initializers.DB.Select("email", "password").Find(&users)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}
