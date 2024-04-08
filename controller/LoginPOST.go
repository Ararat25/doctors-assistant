package controller

import (
	"errors"
	"gorm.io/gorm"
	"net/http"
	"webApp/initializers"
	"webApp/model"
)

func LoginPOST(res http.ResponseWriter, req *http.Request) {
	user := model.User{
		Email:    req.PostFormValue("email"),
		Password: req.PostFormValue("password"),
	}

	userFound := model.User{}

	result := initializers.DB.Where("email = ?", user.Email).First(&userFound)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(res, result.Error.Error(), http.StatusNoContent)
		return
	}

	result = initializers.DB.Where("password = ?", user.Password).First(&userFound)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(res, result.Error.Error(), http.StatusUnauthorized)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}
