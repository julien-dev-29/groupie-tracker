package handlers

import (
	"fmt"
	"html/template"
	"main/api"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var artistTmpl = template.Must(template.New("").Funcs(template.FuncMap{
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

func HandleArtistPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	artistId, _ := strconv.Atoi(id)
	artist, dates, err := api.GetArtistById(artistId)
	if err != nil {
		fmt.Println("Error: ", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}
	data := api.ArtistData{
		Title:  artist.Name,
		Artist: *artist,
		Dates:  *dates,
	}
	artistTmpl.ExecuteTemplate(w, "base.html", data)
}
