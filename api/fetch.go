package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var artistsURL = "https://groupietrackers.herokuapp.com/api/artists"
var artistURL = "https://groupietrackers.herokuapp.com/api/artists/"
var relationURL = "https://groupietrackers.herokuapp.com/api/relation/"

func GetArtists() ([]Artist, error) {
	resp, err := http.Get(artistsURL)
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
	artistResp, err := http.Get(artistURL + strconv.Itoa(id))
	if err != nil {
		return nil, nil, err
	}
	defer artistResp.Body.Close()

	if artistResp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("artist not found: %d", id)
	}

	var artist Artist
	if err := json.NewDecoder(artistResp.Body).Decode(&artist); err != nil {
		return nil, nil, err
	}

	relationResp, err := http.Get(relationURL + strconv.Itoa(id))
	if err != nil {
		return nil, nil, err
	}
	defer relationResp.Body.Close()

	var dates ArtistDates
	if err := json.NewDecoder(relationResp.Body).Decode(&dates); err != nil {
		return nil, nil, err
	}

	return &artist, &dates, nil
}

func GetEnrichedArtists() ([]EnrichedArtist, error) {
	artists, err := GetArtists()
	if err != nil {
		return nil, err
	}

	relations, err := GetAllRelations(len(artists))
	if err != nil {
		return nil, err
	}

	enriched := make([]EnrichedArtist, len(artists))
	for i, a := range artists {
		loc := relations[a.Id]
		if loc == nil {
			loc = []string{}
		}
		enriched[i] = EnrichedArtist{
			Artist:      a,
			Locations:   loc,
			MemberCount: len(a.Members),
			AlbumYear:   extractAlbumYear(a.FirstAlbum),
		}
	}
	return enriched, nil
}

func GetAllRelations(expectedCount int) (map[int][]string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	results := make(map[int][]string)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errCh := make(chan error, expectedCount)

	for i := 1; i <= expectedCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			locations, err := fetchRelationLocations(client, id)
			if err != nil {
				errCh <- err
				return
			}
			mu.Lock()
			results[id] = locations
			mu.Unlock()
		}(i)
	}

	wg.Wait()
	close(errCh)

	if err, ok := <-errCh; ok {
		return nil, err
	}

	return results, nil
}

func fetchRelationLocations(client *http.Client, id int) ([]string, error) {
	resp, err := client.Get(relationURL + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dates ArtistDates
	if err := json.NewDecoder(resp.Body).Decode(&dates); err != nil {
		return nil, err
	}

	locations := make([]string, 0, len(dates.DatesLocations))
	for loc := range dates.DatesLocations {
		locations = append(locations, loc)
	}
	return locations, nil
}

func extractAlbumYear(dateStr string) int {
	parts := strings.Split(dateStr, "-")
	if len(parts) != 3 {
		return 0
	}
	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0
	}
	return year
}
