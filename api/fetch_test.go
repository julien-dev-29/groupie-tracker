package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetArtists_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	originalURL := artistsURL
	artistsURL = ts.URL
	defer func() { artistsURL = originalURL }()

	_, err := GetArtists()
	if err == nil {
		t.Error("expected error for server error, got nil")
	}
}

func TestGetArtists_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer ts.Close()

	originalURL := artistsURL
	artistsURL = ts.URL
	defer func() { artistsURL = originalURL }()

	_, err := GetArtists()
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestGetArtistById_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	originalArtistURL := artistURL
	originalRelationURL := relationURL
	artistURL = ts.URL + "/artist"
	relationURL = ts.URL + "/relation"
	defer func() {
		artistURL = originalArtistURL
		relationURL = originalRelationURL
	}()

	_, _, err := GetArtistById(999)
	if err == nil {
		t.Error("expected error for not found, got nil")
	}
}

func TestArtistJSON(t *testing.T) {
	jsonStr := `{
		"id": 1,
		"image": "https://example.com/image.jpg",
		"name": "Test Band",
		"members": ["Alice", "Bob"],
		"creationDate": 1990,
		"firstAlbum": "01-01-1990"
	}`
	var artist Artist
	if err := json.Unmarshal([]byte(jsonStr), &artist); err != nil {
		t.Fatal("failed to unmarshal artist:", err)
	}

	if artist.Id != 1 {
		t.Errorf("expected Id 1, got %d", artist.Id)
	}
	if artist.Name != "Test Band" {
		t.Errorf("expected Name 'Test Band', got '%s'", artist.Name)
	}
	if len(artist.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(artist.Members))
	}
	if artist.CreationDate != 1990 {
		t.Errorf("expected CreationDate 1990, got %d", artist.CreationDate)
	}
}

func TestArtistDatesJSON(t *testing.T) {
	jsonStr := `{
		"id": 1,
		"datesLocations": {
			"location1": ["2024-01-01", "2024-02-01"],
			"location2": ["2024-03-01"]
		}
	}`
	var dates ArtistDates
	if err := json.Unmarshal([]byte(jsonStr), &dates); err != nil {
		t.Fatal("failed to unmarshal artist dates:", err)
	}

	if dates.Id != 1 {
		t.Errorf("expected Id 1, got %d", dates.Id)
	}
	if len(dates.DatesLocations) != 2 {
		t.Errorf("expected 2 locations, got %d", len(dates.DatesLocations))
	}
	if len(dates.DatesLocations["location1"]) != 2 {
		t.Errorf("expected 2 dates for location1, got %d", len(dates.DatesLocations["location1"]))
	}
}

func TestExtractAlbumYear_Valid(t *testing.T) {
	year := extractAlbumYear("15-09-2015")
	if year != 2015 {
		t.Errorf("expected 2015, got %d", year)
	}
}

func TestExtractAlbumYear_Invalid(t *testing.T) {
	tests := []string{"not-a-date", "15-09", "", "abc-def-ghij"}
	for _, d := range tests {
		if y := extractAlbumYear(d); y != 0 {
			t.Errorf("extractAlbumYear(%q) expected 0, got %d", d, y)
		}
	}
}

func TestGetAllRelations_MultipleArtists(t *testing.T) {
	requestCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		resp := ArtistDates{
			Id: requestCount,
			DatesLocations: map[string][]string{
				"location-a": {"2024-01-01"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	originalURL := relationURL
	relationURL = ts.URL + "/relation"
	defer func() { relationURL = originalURL }()

	results, err := GetAllRelations(3)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(results) != 3 {
		t.Errorf("expected 3 results, got %d", len(results))
	}
	for id := 1; id <= 3; id++ {
		locs, ok := results[id]
		if !ok {
			t.Errorf("missing results for id %d", id)
			continue
		}
		if len(locs) != 1 || locs[0] != "location-a" {
			t.Errorf("id %d: expected [location-a], got %v", id, locs)
		}
	}
}

func TestGetEnrichedArtists_AlbumYearParsing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/artists" {
			artists := []Artist{
				{Id: 1, Name: "Band A", FirstAlbum: "15-09-2015", Members: []string{"a", "b"}},
				{Id: 2, Name: "Band B", FirstAlbum: "01-01-1990", Members: []string{"x"}},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(artists)
		} else {
			resp := ArtistDates{
				DatesLocations: map[string][]string{"loc": {"date"}},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
		}
	}))
	defer ts.Close()

	originalArtistsURL := artistsURL
	originalRelationURL := relationURL
	artistsURL = ts.URL + "/artists"
	relationURL = ts.URL + "/relation"
	defer func() {
		artistsURL = originalArtistsURL
		relationURL = originalRelationURL
	}()

	enriched, err := GetEnrichedArtists()
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if len(enriched) != 2 {
		t.Fatalf("expected 2 artists, got %d", len(enriched))
	}
	if enriched[0].AlbumYear != 2015 {
		t.Errorf("Band A: expected AlbumYear 2015, got %d", enriched[0].AlbumYear)
	}
	if enriched[1].AlbumYear != 1990 {
		t.Errorf("Band B: expected AlbumYear 1990, got %d", enriched[1].AlbumYear)
	}
	if enriched[0].MemberCount != 2 {
		t.Errorf("Band A: expected MemberCount 2, got %d", enriched[0].MemberCount)
	}
	if enriched[1].MemberCount != 1 {
		t.Errorf("Band B: expected MemberCount 1, got %d", enriched[1].MemberCount)
	}
	if len(enriched[0].Locations) != 1 || enriched[0].Locations[0] != "loc" {
		t.Errorf("Band A: expected locations [loc], got %v", enriched[0].Locations)
	}
}
