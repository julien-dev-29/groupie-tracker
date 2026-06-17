package viewmodels

import (
	"fmt"
	"main/api"
	"net/url"
	"strconv"
)

type PageData struct {
	Artists               []api.EnrichedArtist
	Title                 string
	DisplayMode           string
	CurrentPage           int
	TotalPages            int
	HasPrevious           bool
	HasNext               bool
	PreviousPage          int
	NextPage              int
	Query                 string
	Member                string
	MinYear               int
	MaxYear               int
	MinAlbumYear          int
	MaxAlbumYear          int
	SelectedMemberCounts  []int
	SelectedLocations     []string
	AvailableMemberCounts []int
	AvailableLocations    []api.LocationNode
}

func (p PageData) FilterQuery() string {
	vals := url.Values{}
	if p.Query != "" {
		vals.Set("q", p.Query)
	}
	if p.Member != "" {
		vals.Set("member", p.Member)
	}
	if p.MinYear > 0 {
		vals.Set("minYear", strconv.Itoa(p.MinYear))
	}
	if p.MaxYear > 0 {
		vals.Set("maxYear", strconv.Itoa(p.MaxYear))
	}
	if p.MinAlbumYear > 0 {
		vals.Set("minAlbum", strconv.Itoa(p.MinAlbumYear))
	}
	if p.MaxAlbumYear > 0 {
		vals.Set("maxAlbum", strconv.Itoa(p.MaxAlbumYear))
	}
	if len(p.SelectedMemberCounts) > 0 {
		for _, c := range p.SelectedMemberCounts {
			vals.Add("members", strconv.Itoa(c))
		}
	}
	if len(p.SelectedLocations) > 0 {
		vals["locations"] = p.SelectedLocations
	}
	vals.Set("display", p.DisplayMode)
	if len(vals) > 0 {
		return "&" + vals.Encode()
	}
	return ""
}

func (p PageData) PageURL(pageNum int) string {
	vals := url.Values{}
	vals.Set("page", strconv.Itoa(pageNum))
	vals.Set("display", p.DisplayMode)
	if p.Query != "" {
		vals.Set("q", p.Query)
	}
	if p.Member != "" {
		vals.Set("member", p.Member)
	}
	if p.MinYear > 0 {
		vals.Set("minYear", strconv.Itoa(p.MinYear))
	}
	if p.MaxYear > 0 {
		vals.Set("maxYear", strconv.Itoa(p.MaxYear))
	}
	if p.MinAlbumYear > 0 {
		vals.Set("minAlbum", strconv.Itoa(p.MinAlbumYear))
	}
	if p.MaxAlbumYear > 0 {
		vals.Set("maxAlbum", strconv.Itoa(p.MaxAlbumYear))
	}
	if len(p.SelectedMemberCounts) > 0 {
		for _, c := range p.SelectedMemberCounts {
			vals.Add("members", strconv.Itoa(c))
		}
	}
	if len(p.SelectedLocations) > 0 {
		vals["locations"] = p.SelectedLocations
	}
	return "/?" + vals.Encode()
}

func (p PageData) HasActiveFilter() bool {
	return p.Query != "" ||
		p.Member != "" ||
		p.MinYear > 0 ||
		p.MaxYear > 0 ||
		p.MinAlbumYear > 0 ||
		p.MaxAlbumYear > 0 ||
		len(p.SelectedMemberCounts) > 0 ||
		len(p.SelectedLocations) > 0
}

func (p PageData) ResultSummary() string {
	if !p.HasActiveFilter() {
		return ""
	}
	return fmt.Sprintf("Found %d result(s)", len(p.Artists))
}
