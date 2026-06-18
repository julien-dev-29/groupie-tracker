package api

import (
	"slices"
	"sort"
	"strings"
)

type LocationNode struct {
	Label    string
	Value    string
	FullPath string
	Children []LocationNode
}

func ParseLocation(raw string) []string {
	return strings.Split(raw, "-")
}

func Humanize(segment string) string {
	words := strings.Split(segment, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func BuildLocationTree(allLocations []string) []LocationNode {
	countryMap := make(map[string]map[string]map[string]bool)

	for _, loc := range allLocations {
		parts := ParseLocation(loc)
		if len(parts) == 0 {
			continue
		}

		country := parts[len(parts)-1]
		if countryMap[country] == nil {
			countryMap[country] = make(map[string]map[string]bool)
		}

		if len(parts) == 1 {
			if countryMap[country][""] == nil {
				countryMap[country][""] = make(map[string]bool)
			}
			countryMap[country][""][""] = true
		} else if len(parts) == 2 {
			state := parts[0]
			if countryMap[country][state] == nil {
				countryMap[country][state] = make(map[string]bool)
			}
			countryMap[country][state][""] = true
		} else {
			state := parts[len(parts)-2]
			city := strings.Join(parts[:len(parts)-2], "-")
			if countryMap[country][state] == nil {
				countryMap[country][state] = make(map[string]bool)
			}
			countryMap[country][state][city] = true
		}
	}

	var roots []LocationNode
	countries := sortedKeys(countryMap)
	for _, country := range countries {
		countryNode := LocationNode{
			Label:    Humanize(country),
			Value:    country,
			FullPath: country,
		}
		states := countryMap[country]
		stateNames := sortedKeys(states)

		for _, state := range stateNames {
			cities := states[state]
			if state == "" {
				continue
			}

			stateFullPath := state + "-" + country
			stateNode := LocationNode{
				Label:    Humanize(state),
				Value:    state,
				FullPath: stateFullPath,
			}

			var cityNames []string
			for c := range cities {
				cityNames = append(cityNames, c)
			}
			sort.Strings(cityNames)

			for _, city := range cityNames {
				if city == "" {
					continue
				}
				cityFullPath := city + "-" + stateFullPath
				stateNode.Children = append(stateNode.Children, LocationNode{
					Label:    Humanize(city),
					Value:    city,
					FullPath: cityFullPath,
				})
			}

			countryNode.Children = append(countryNode.Children, stateNode)
		}

		roots = append(roots, countryNode)
	}
	return roots
}

func LocationMatches(artistLocations []string, selected string) bool {
	for _, loc := range artistLocations {
		if strings.HasSuffix(loc, selected) {
			if len(loc) == len(selected) || loc[len(loc)-len(selected)-1] == '-' {
				return true
			}
		}
	}
	return false
}

func CollectAllUniqueLocations(artists []EnrichedArtist) []string {
	seen := make(map[string]bool)
	var result []string
	for _, a := range artists {
		for _, loc := range a.Locations {
			if !seen[loc] {
				seen[loc] = true
				result = append(result, loc)
			}
		}
	}
	sort.Strings(result)
	return result
}

func CollectMemberCounts(artists []EnrichedArtist) []int {
	seen := make(map[int]bool)
	for _, a := range artists {
		seen[a.MemberCount] = true
	}
	var counts []int
	for c := range seen {
		counts = append(counts, c)
	}
	sort.Ints(counts)
	return counts
}

func sortedKeys[K ~string, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
