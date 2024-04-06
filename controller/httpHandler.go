package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

func HandleMain(res http.ResponseWriter, req *http.Request) {
	s := fmt.Sprintf("Method: %s\nHost: %s\nPath: %s",
		req.Method, req.Host, req.URL.Path)
	_, err := res.Write([]byte(s))
	if err != nil {
		fmt.Printf("Ошибка записи ответа: %s\n", err.Error())
		return
	}
}

func HandleLogin(res http.ResponseWriter, req *http.Request) {
	// Используем функцию template.ParseFiles() для чтения файла шаблона
	ts, err := template.ParseFiles("./view/login/page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "Internal Server Error", 500)
		return
	}

	// Затем мы используем метод Execute() для записи содержимого
	// шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// возможность отправки динамических данных в шаблон.
	err = ts.Execute(res, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(res, "Internal Server Error", 500)
	}
}

type msg struct {
	Message string
}

func HandlePOST(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		fmt.Printf("Email: %s\nPassword: %s\n", req.PostFormValue("email"), req.PostFormValue("password"))
		resp, err := json.Marshal(msg{
			Message: "ok",
		})
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(resp)
		return
	}
	io.WriteString(res, "Отправьте POST запрос с параметрами email и name")
}
