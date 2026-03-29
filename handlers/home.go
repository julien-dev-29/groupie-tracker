package handlers

import (
	"fmt"
	"html/template"
	"main/api"
	"main/viewmodels"
	"net/http"
	"strconv"
	"strings"
)

var homeTmpl = template.Must(template.New("").Funcs(template.FuncMap{
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
	"views/index.html",
))

func HandleHome(w http.ResponseWriter, r *http.Request) {
	artists, err := api.GetArtists()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	displayMode := r.URL.Query().Get("display")
	if displayMode == "" {
		displayMode = "card"
	}

	// Pagination
	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	itemsPerPage := 12
	totalItems := len(artists)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage

	// Validate page number
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	// Calculate start and end indices
	startIdx := (page - 1) * itemsPerPage
	endIdx := min(startIdx+itemsPerPage, totalItems)

	var displayedArtists []api.Artist
	if totalItems > 0 {
		displayedArtists = artists[startIdx:endIdx]
	}

	data := viewmodels.PageData{
		Artists:      displayedArtists,
		Title:        "Groupie Tracker",
		DisplayMode:  displayMode,
		CurrentPage:  page,
		TotalPages:   totalPages,
		HasPrevious:  page > 1,
		HasNext:      page < totalPages,
		PreviousPage: page - 1,
		NextPage:     page + 1,
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	homeTmpl.ExecuteTemplate(w, "base.html", data)
}
