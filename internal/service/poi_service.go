package service

import (
	"sort"
	"strings"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/store"
)

type POIMatchResult struct {
	POI     *models.POI
	Score   int
}

func MatchPOIs(query string, limit int) []POIMatchResult {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}
	queryLower := strings.ToLower(query)
	queryRunes := []rune(queryLower)

	pois := store.GlobalStore.ListPOIs()
	var results []POIMatchResult

	for _, poi := range pois {
		score := calculateMatchScore(poi, queryLower, queryRunes)
		if score > 0 {
			results = append(results, POIMatchResult{
				POI:   poi,
				Score: score,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	return results
}

func calculateMatchScore(poi *models.POI, queryLower string, queryRunes []rune) int {
	score := 0
	poiNameLower := strings.ToLower(poi.Name)
	poiAddrLower := strings.ToLower(poi.Address)

	if poiNameLower == queryLower {
		score += 100
	}
	if poiAddrLower == queryLower {
		score += 90
	}

	if strings.HasPrefix(poiNameLower, queryLower) {
		score += 80
	}

	if strings.Contains(poiNameLower, queryLower) {
		score += 50
	}
	if strings.Contains(poiAddrLower, queryLower) {
		score += 30
	}

	score += substringMatchScore(poiNameLower, queryRunes) * 5
	score += substringMatchScore(poiAddrLower, queryRunes) * 3

	return score
}

func substringMatchScore(target string, queryRunes []rune) int {
	if len(queryRunes) == 0 {
		return 0
	}
	targetRunes := []rune(target)
	matches := 0
	idx := 0
	for _, qr := range queryRunes {
		for i := idx; i < len(targetRunes); i++ {
			if targetRunes[i] == qr {
				matches++
				idx = i + 1
				break
			}
		}
	}
	return matches
}
