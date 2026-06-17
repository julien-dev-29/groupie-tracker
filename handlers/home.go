package handlers

import (
	"fmt"
	"html/template"
	"main/api"
	"main/viewmodels"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	homeTmplOnce sync.Once
	homeTmpl     *template.Template
)

func homeTemplate() *template.Template {
	homeTmplOnce.Do(func() {
		homeTmpl = template.Must(template.New("").Funcs(template.FuncMap{
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
			"stringInSlice": func(s string, slice []string) bool {
				for _, v := range slice {
					if v == s {
						return true
					}
				}
				return false
			},
			"intInSlice": func(n int, slice []int) bool {
				for _, v := range slice {
					if v == n {
						return true
					}
				}
				return false
			},
			"title": func(s string) string {
				if len(s) == 0 {
					return s
				}
				return strings.ToUpper(s[:1]) + s[1:]
			},
		}).ParseFiles(
			"views/base.html",
			"views/header.html",
			"views/footer.html",
			"views/index.html",
		))
	})
	return homeTmpl
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	artists, err := api.GetEnrichedArtists()
	if err != nil {
		fmt.Println("Error fetching artists:", err)
		http.Error(w, "Failed to load artists", http.StatusInternalServerError)
		return
	}

	displayMode := r.URL.Query().Get("display")
	if displayMode == "" {
		displayMode = "card"
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))
	memberFilter := strings.TrimSpace(r.URL.Query().Get("member"))
	minYear, _ := strconv.Atoi(r.URL.Query().Get("minYear"))
	maxYear, _ := strconv.Atoi(r.URL.Query().Get("maxYear"))
	minAlbum, _ := strconv.Atoi(r.URL.Query().Get("minAlbum"))
	maxAlbum, _ := strconv.Atoi(r.URL.Query().Get("maxAlbum"))

	var selectedMemberCounts []int
	for _, s := range r.URL.Query()["members"] {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			selectedMemberCounts = append(selectedMemberCounts, n)
		}
	}

	selectedLocations := r.URL.Query()["locations"]

	availableMemberCounts := api.CollectMemberCounts(artists)
	allLocations := api.CollectAllUniqueLocations(artists)
	availableLocations := api.BuildLocationTree(allLocations)

	hasFilters := query != "" ||
		memberFilter != "" ||
		minYear > 0 ||
		maxYear > 0 ||
		minAlbum > 0 ||
		maxAlbum > 0 ||
		len(selectedMemberCounts) > 0 ||
		len(selectedLocations) > 0

	if hasFilters {
		var filtered []api.EnrichedArtist
		queryLower := strings.ToLower(query)
		memberLower := strings.ToLower(memberFilter)

		for _, a := range artists {
			if query != "" && !strings.Contains(strings.ToLower(a.Name), queryLower) {
				continue
			}
			if memberFilter != "" {
				memberMatch := false
				for _, m := range a.Members {
					if strings.Contains(strings.ToLower(m), memberLower) {
						memberMatch = true
						break
					}
				}
				if !memberMatch {
					continue
				}
			}
			if minYear > 0 && a.CreationDate < minYear {
				continue
			}
			if maxYear > 0 && a.CreationDate > maxYear {
				continue
			}
			if minAlbum > 0 && a.AlbumYear < minAlbum {
				continue
			}
			if maxAlbum > 0 && a.AlbumYear > maxAlbum {
				continue
			}
			if len(selectedMemberCounts) > 0 {
				countMatch := false
				for _, c := range selectedMemberCounts {
					if a.MemberCount == c {
						countMatch = true
						break
					}
				}
				if !countMatch {
					continue
				}
			}
			if len(selectedLocations) > 0 {
				locMatch := false
				for _, sel := range selectedLocations {
					if api.LocationMatches(a.Locations, sel) {
						locMatch = true
						break
					}
				}
				if !locMatch {
					continue
				}
			}
			filtered = append(filtered, a)
		}
		artists = filtered
	}

	page := 1
	if p := r.URL.Query().Get("page"); p != "" {
		if parsedPage, err := strconv.Atoi(p); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	itemsPerPage := 12
	totalItems := len(artists)
	totalPages := (totalItems + itemsPerPage - 1) / itemsPerPage
	if totalPages == 0 {
		totalPages = 1
	}

	if page > totalPages {
		page = totalPages
	}

	startIdx := (page - 1) * itemsPerPage
	endIdx := min(startIdx+itemsPerPage, totalItems)

	var displayedArtists []api.EnrichedArtist
	if totalItems > 0 {
		displayedArtists = artists[startIdx:endIdx]
	}

	data := viewmodels.PageData{
		Artists:               displayedArtists,
		Title:                 "Groupie Tracker",
		DisplayMode:           displayMode,
		CurrentPage:           page,
		TotalPages:            totalPages,
		HasPrevious:           page > 1,
		HasNext:               page < totalPages,
		PreviousPage:          page - 1,
		NextPage:              page + 1,
		Query:                 query,
		Member:                memberFilter,
		MinYear:               minYear,
		MaxYear:               maxYear,
		MinAlbumYear:          minAlbum,
		MaxAlbumYear:          maxAlbum,
		SelectedMemberCounts:  selectedMemberCounts,
		SelectedLocations:     selectedLocations,
		AvailableMemberCounts: availableMemberCounts,
		AvailableLocations:    availableLocations,
	}

	if err := homeTemplate().ExecuteTemplate(w, "base.html", data); err != nil {
		fmt.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
