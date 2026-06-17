package handlers

import (
	"fmt"
	"html/template"
	"main/api"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

var (
	artistTmplOnce sync.Once
	artistTmpl     *template.Template
)

func artistTemplate() *template.Template {
	artistTmplOnce.Do(func() {
		artistTmpl = template.Must(template.New("").Funcs(template.FuncMap{
			"until": func(n int) []int {
				var arr []int
				for i := range n {
					arr = append(arr, i)
				}
				return arr
			},
			"add": func(a, b int) int {
				return a + b
			},
			"replaceReservationChars": func(s string) string {
				return strings.ReplaceAll(s, " ", "")
			},
		}).ParseFiles(
			"views/base.html",
			"views/header.html",
			"views/footer.html",
			"views/artist.html",
		))
	})
	return artistTmpl
}

func HandleArtistPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	artistId, err := strconv.Atoi(id)
	if err != nil || artistId < 1 {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	artist, dates, err := api.GetArtistById(artistId)
	if err != nil {
		fmt.Println("Error fetching artist:", err)
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	data := api.ArtistData{
		Title:  artist.Name,
		Artist: *artist,
		Dates:  *dates,
	}
	if err := artistTemplate().ExecuteTemplate(w, "base.html", data); err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
