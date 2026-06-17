package main

import (
	"fmt"
	"log"
	"main/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()

	mux.HandleFunc("/", handlers.HandleHome)
	mux.HandleFunc("/artist/{id}", handlers.HandleArtistPage)

	mux.NotFoundHandler = http.HandlerFunc(handlers.HandleNotFound)

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Println("http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
