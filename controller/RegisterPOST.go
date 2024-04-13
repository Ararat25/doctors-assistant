package controller

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"webApp/initializers"
	"webApp/model"
)

func RegisterPOST(res http.ResponseWriter, req *http.Request) {
	user := model.User{
		Email:      req.PostFormValue("email"),
		Password:   req.PostFormValue("password"),
		LastName:   req.PostFormValue("lastName"),
		FirstName:  req.PostFormValue("firstName"),
		MiddleName: req.PostFormValue("middleName"),
		Specialty:  req.PostFormValue("specialty"),
	}

	hashPassword := sha256.Sum256([]byte(user.Password))
	hashString := hex.EncodeToString(hashPassword[:])
	user.Password = hashString

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		http.Error(res, result.Error.Error(), http.StatusConflict)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
}
