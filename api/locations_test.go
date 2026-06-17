package api

import (
	"testing"
)

func TestParseLocation(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{"seattle-washington-usa", []string{"seattle", "washington", "usa"}},
		{"london-uk", []string{"london", "uk"}},
		{"usa", []string{"usa"}},
		{"", []string{""}},
	}
	for _, tt := range tests {
		got := ParseLocation(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("ParseLocation(%q) returned %v, want %v", tt.input, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("ParseLocation(%q) = %v, want %v", tt.input, got, tt.want)
				break
			}
		}
	}
}

func TestHumanize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"seattle", "Seattle"},
		{"los-angeles", "Los Angeles"},
		{"usa", "Usa"},
		{"north-carolina", "North Carolina"},
	}
	for _, tt := range tests {
		got := Humanize(tt.input)
		if got != tt.want {
			t.Errorf("Humanize(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBuildLocationTree(t *testing.T) {
	locations := []string{
		"seattle-washington-usa",
		"los-angeles-california-usa",
		"london-uk",
		"glasgow-uk",
	}
	tree := BuildLocationTree(locations)

	if len(tree) != 2 {
		t.Fatalf("expected 2 countries, got %d", len(tree))
	}

	var usaNode, ukNode *LocationNode
	for i := range tree {
		switch tree[i].Value {
		case "usa":
			usaNode = &tree[i]
		case "uk":
			ukNode = &tree[i]
		}
	}

	if usaNode == nil {
		t.Fatal("expected usa node in tree")
	}
	if usaNode.Label != "Usa" {
		t.Errorf("expected usa label 'Usa', got %q", usaNode.Label)
	}
	if len(usaNode.Children) != 2 {
		t.Fatalf("expected 2 states under usa, got %d", len(usaNode.Children))
	}

	if ukNode == nil {
		t.Fatal("expected uk node in tree")
	}
	if len(ukNode.Children) != 2 {
		t.Fatalf("expected 2 cities under uk, got %d", len(ukNode.Children))
	}
}

func TestLocationMatches(t *testing.T) {
	artistLocs := []string{
		"seattle-washington-usa",
		"london-uk",
	}

	tests := []struct {
		selected string
		want     bool
	}{
		{"usa", true},
		{"washington-usa", true},
		{"seattle-washington-usa", true},
		{"uk", true},
		{"london-uk", true},
		{"canada", false},
		{"washington", false},
	}
	for _, tt := range tests {
		got := LocationMatches(artistLocs, tt.selected)
		if got != tt.want {
			t.Errorf("LocationMatches(%v, %q) = %v, want %v", artistLocs, tt.selected, got, tt.want)
		}
	}
}

func TestCollectAllUniqueLocations(t *testing.T) {
	artists := []EnrichedArtist{
		{Locations: []string{"seattle-washington-usa", "london-uk"}},
		{Locations: []string{"los-angeles-california-usa", "london-uk"}},
		{Locations: []string{"tokyo-japan"}},
	}
	locs := CollectAllUniqueLocations(artists)

	expected := []string{
		"london-uk",
		"los-angeles-california-usa",
		"seattle-washington-usa",
		"tokyo-japan",
	}
	if len(locs) != len(expected) {
		t.Fatalf("expected %d unique locations, got %d: %v", len(expected), len(locs), locs)
	}
	for i := range expected {
		if locs[i] != expected[i] {
			t.Errorf("expected locs[%d] = %q, got %q", i, expected[i], locs[i])
		}
	}
}

func TestCollectMemberCounts(t *testing.T) {
	artists := []EnrichedArtist{
		{MemberCount: 1},
		{MemberCount: 4},
		{MemberCount: 5},
		{MemberCount: 1},
	}
	counts := CollectMemberCounts(artists)

	expected := []int{1, 4, 5}
	if len(counts) != len(expected) {
		t.Fatalf("expected %d unique counts, got %d: %v", len(expected), len(counts), counts)
	}
	for i := range expected {
		if counts[i] != expected[i] {
			t.Errorf("expected counts[%d] = %d, got %d", i, expected[i], counts[i])
		}
	}
}
