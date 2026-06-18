package viewmodels

import (
	"main/api"
	"strings"
	"testing"
)

func TestFilterQuery_Empty(t *testing.T) {
	p := PageData{DisplayMode: "card"}
	fq := p.FilterQuery()
	if fq != "" {
		t.Errorf("expected empty FilterQuery for no active filters, got %s", fq)
	}
}

func TestFilterQuery_WithQuery(t *testing.T) {
	p := PageData{DisplayMode: "card", Query: "queen"}
	fq := p.FilterQuery()
	if !strings.Contains(fq, "q=queen") {
		t.Errorf("expected FilterQuery to contain q=queen, got %s", fq)
	}
}

func TestFilterQuery_WithMember(t *testing.T) {
	p := PageData{DisplayMode: "card", Member: "freddie"}
	fq := p.FilterQuery()
	if !strings.Contains(fq, "member=freddie") {
		t.Errorf("expected FilterQuery to contain member=freddie, got %s", fq)
	}
}

func TestPageURL_IncludesFilters(t *testing.T) {
	p := PageData{DisplayMode: "card", Query: "queen", Member: "freddie", CurrentPage: 1}
	url := p.PageURL(2)
	if !strings.Contains(url, "page=2") {
		t.Errorf("expected PageURL to contain page=2, got %s", url)
	}
	if !strings.Contains(url, "q=queen") {
		t.Errorf("expected PageURL to contain q=queen, got %s", url)
	}
	if !strings.Contains(url, "member=freddie") {
		t.Errorf("expected PageURL to contain member=freddie, got %s", url)
	}
}

func TestPageURL_FirstPage(t *testing.T) {
	p := PageData{DisplayMode: "card"}
	url := p.PageURL(1)
	if !strings.Contains(url, "page=1") {
		t.Errorf("expected PageURL to contain page=1, got %s", url)
	}
}

func TestHasActiveFilter_False(t *testing.T) {
	p := PageData{}
	if p.HasActiveFilter() {
		t.Error("expected HasActiveFilter to be false for empty PageData")
	}
}

func TestHasActiveFilter_True(t *testing.T) {
	p := PageData{Query: "test"}
	if !p.HasActiveFilter() {
		t.Error("expected HasActiveFilter to be true when Query is set")
	}
}

func TestHasActiveFilter_Member(t *testing.T) {
	p := PageData{Member: "test"}
	if !p.HasActiveFilter() {
		t.Error("expected HasActiveFilter to be true when Member is set")
	}
}

func TestHasActiveFilter_MinYear(t *testing.T) {
	p := PageData{MinYear: 2000}
	if !p.HasActiveFilter() {
		t.Error("expected HasActiveFilter to be true when MinYear is set")
	}
}

func TestHasActiveFilter_MinAlbumYear(t *testing.T) {
	p := PageData{MinAlbumYear: 2000}
	if !p.HasActiveFilter() {
		t.Error("expected HasActiveFilter to be true when MinAlbumYear is set")
	}
}

func TestHasActiveFilter_SelectedMemberCounts(t *testing.T) {
	p := PageData{SelectedMemberCounts: []int{3}}
	if !p.HasActiveFilter() {
		t.Error("expected HasActiveFilter to be true with SelectedMemberCounts")
	}
}

func TestHasActiveFilter_SelectedLocations(t *testing.T) {
	p := PageData{SelectedLocations: []string{"usa"}}
	if !p.HasActiveFilter() {
		t.Error("expected HasActiveFilter to be true with SelectedLocations")
	}
}

func TestFilterQuery_WithSelectedMemberCounts(t *testing.T) {
	p := PageData{DisplayMode: "card", SelectedMemberCounts: []int{1, 3, 5}}
	fq := p.FilterQuery()
	if !strings.Contains(fq, "members=1") {
		t.Errorf("expected FilterQuery to contain members=1, got %s", fq)
	}
	if !strings.Contains(fq, "members=3") {
		t.Errorf("expected FilterQuery to contain members=3, got %s", fq)
	}
	if !strings.Contains(fq, "members=5") {
		t.Errorf("expected FilterQuery to contain members=5, got %s", fq)
	}
}

func TestFilterQuery_WithSelectedLocations(t *testing.T) {
	p := PageData{DisplayMode: "card", SelectedLocations: []string{"usa", "london-uk"}}
	fq := p.FilterQuery()
	if !strings.Contains(fq, "locations=usa") {
		t.Errorf("expected FilterQuery to contain locations=usa, got %s", fq)
	}
	if !strings.Contains(fq, "locations=london-uk") {
		t.Errorf("expected FilterQuery to contain locations=london-uk, got %s", fq)
	}
}

func TestPageURL_PreservesAlbumFilter(t *testing.T) {
	p := PageData{DisplayMode: "card", MinAlbumYear: 2000, MaxAlbumYear: 2010}
	url := p.PageURL(1)
	if !strings.Contains(url, "minAlbum=2000") {
		t.Errorf("expected PageURL to contain minAlbum=2000, got %s", url)
	}
	if !strings.Contains(url, "maxAlbum=2010") {
		t.Errorf("expected PageURL to contain maxAlbum=2010, got %s", url)
	}
}

func TestPageURL_PreservesMemberCountFilter(t *testing.T) {
	p := PageData{DisplayMode: "card", SelectedMemberCounts: []int{1, 4}}
	url := p.PageURL(2)
	if !strings.Contains(url, "members=1") {
		t.Errorf("expected PageURL to contain members=1, got %s", url)
	}
	if !strings.Contains(url, "members=4") {
		t.Errorf("expected PageURL to contain members=4, got %s", url)
	}
}

func TestPageURL_PreservesLocationFilter(t *testing.T) {
	p := PageData{DisplayMode: "card", SelectedLocations: []string{"usa"}}
	url := p.PageURL(1)
	if !strings.Contains(url, "locations=usa") {
		t.Errorf("expected PageURL to contain locations=usa, got %s", url)
	}
}

func TestResultSummary_NoFilter(t *testing.T) {
	p := PageData{}
	if s := p.ResultSummary(); s != "" {
		t.Errorf("expected empty summary, got '%s'", s)
	}
}

func TestResultSummary_WithFilter(t *testing.T) {
	p := PageData{
		Query: "queen",
		Artists: []api.EnrichedArtist{
			{Artist: api.Artist{Id: 1}},
			{Artist: api.Artist{Id: 2}},
		},
	}
	s := p.ResultSummary()
	if !strings.Contains(s, "2") {
		t.Errorf("expected summary to contain '2', got '%s'", s)
	}
}
