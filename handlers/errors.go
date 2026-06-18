package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

var (
	errorTmplOnce sync.Once
	errorTmpl     *template.Template
)

func errorTemplate() *template.Template {
	errorTmplOnce.Do(func() {
		errorTmpl = template.Must(template.ParseFiles(
			"views/base.html",
			"views/header.html",
			"views/footer.html",
			"views/error.html",
		))
	})
	return errorTmpl
}

type ErrorData struct {
	Title   string
	Message string
	Code    int
}

func renderError(w http.ResponseWriter, status int, title, message string) {
	w.WriteHeader(status)
	data := ErrorData{Title: title, Message: message, Code: status}
	if err := errorTemplate().ExecuteTemplate(w, "base.html", data); err != nil {
		fmt.Println("Error executing error template:", err)
		http.Error(w, message, status)
	}
}

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	renderError(w, http.StatusNotFound, "404 - Page Not Found", "The page you are looking for does not exist.")
}
