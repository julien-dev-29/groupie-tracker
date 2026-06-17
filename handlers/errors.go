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

func HandleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	data := struct {
		Title   string
		Message string
		Code    int
	}{
		Title:   "404 - Page Not Found",
		Message: "The page you are looking for does not exist.",
		Code:    404,
	}
	if err := errorTemplate().ExecuteTemplate(w, "base.html", data); err != nil {
		fmt.Println("Error executing 404 template:", err)
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
