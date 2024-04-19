package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type JSONBody struct {
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}

func VerifyUser(token string, email string) bool {
	claims := jwt.MapClaims{}

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SALT")), nil
	})
	if err != nil {
		fmt.Printf("Failed to parse token: %s\n", err)
		return false
	}

	if !jwtToken.Valid {
		return false
	}

	emailRow, ok := claims["sub"]
	if !ok {
		return false
	}

	expRow, ok := claims["exp"]
	if !ok {
		return false
	}

	emailFromReq, ok := emailRow.(string)
	if !ok {
		return false
	}

	expFromReq, ok := expRow.(float64)
	if !ok {
		return false
	}

	if err != nil {
		return false
	}

	if emailFromReq == email {
		if time.Now().Unix() < int64(expFromReq) {
			return true
		}
	}
	return false
}

func AccountPage(res http.ResponseWriter, req *http.Request) {
	var jsonBody JSONBody
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(buf.Bytes(), &jsonBody)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if VerifyUser(jsonBody.AccessToken, jsonBody.Email) {
		ts, err := template.ParseFiles("./view/accountPage/page.tmpl")
		if err != nil {
			log.Println(err.Error())
			http.Error(res, "Internal Server Error", 500)
			return
		}

		err = ts.Execute(res, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(res, "Internal Server Error", 500)
			return
		}
		return
	}
	http.Error(res, "Не валидный токен", http.StatusUnauthorized)
}

func Page(res http.ResponseWriter, req *http.Request) {
	ts, err := template.ParseFiles("./view/accountPage/page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(res, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}
	return
}
