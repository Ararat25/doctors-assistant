package account

import (
	"html/template"
	"log"
	"net/http"
)

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
