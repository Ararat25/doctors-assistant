package login

import (
	"html/template"
	"log"
	"net/http"
)

func Page(res http.ResponseWriter, req *http.Request) {
	// Используем функцию template.ParseFiles() для чтения файла шаблона
	ts, err := template.ParseFiles("./view/loginPage/page.tmpl")
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
