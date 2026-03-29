package api

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func GetArtists() ([]Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return []Artist{}, err
	}
	defer resp.Body.Close()
	var artists []Artist

	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return []Artist{}, err
	}
	return artists, nil
}

func GetArtistById(id int) (*Artist, *ArtistDates, error) {
	// Get artist data
	artistResp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists/" + strconv.Itoa(id))
	if err != nil {
		return &Artist{}, &ArtistDates{}, err
	}
	defer artistResp.Body.Close()
	var artist Artist
	if err := json.NewDecoder(artistResp.Body).Decode(&artist); err != nil {
		return &Artist{}, &ArtistDates{}, err
	}

	// Get relation data (dates and locations)
	relationResp, err := http.Get("https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(id))
	if err != nil {
		return &Artist{}, &ArtistDates{}, err
	}
	defer relationResp.Body.Close()
	var dates ArtistDates
	if err := json.NewDecoder(relationResp.Body).Decode(&dates); err != nil {
		return &Artist{}, &ArtistDates{}, err
	}

	return &artist, &dates, nil
}
