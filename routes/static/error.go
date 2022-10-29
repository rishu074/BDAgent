package static

import (
	"html/template"
	"net/http"
)

type ErrorPageData struct {
	Code    int
	Message string
}

func ErrorRouteHandler(w http.ResponseWriter, r *http.Request, msg string, code int) {
	pageTemplate := template.Must(template.ParseFiles("./public/error.html"))

	w.WriteHeader(http.StatusNotFound)
	pageTemplate.Execute(w, ErrorPageData{
		Code:    code,
		Message: msg,
	})
}
