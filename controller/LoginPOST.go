package controller

import (
	"crypto/sha256"
	"encoding/hex"
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

	result := initializers.DB.Where(&model.User{Email: user.Email}).First(&userFound)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(res, result.Error.Error(), http.StatusNoContent)
		return
	}

	hashPassword := sha256.Sum256([]byte(user.Password))
	hashString := hex.EncodeToString(hashPassword[:])

	result = initializers.DB.Where(&model.User{Email: user.Email, Password: hashString}).First(&userFound)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		http.Error(res, result.Error.Error(), http.StatusUnauthorized)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}
