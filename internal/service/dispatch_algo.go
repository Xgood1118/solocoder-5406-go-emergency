package service

import (
	"math"
	"sort"

	"emergency-dispatch/internal/models"
	"emergency-dispatch/internal/store"
)

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371.0
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
		math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func GetRequiredVehicleType(illnessType models.IllnessType) models.VehicleType {
	switch illnessType {
	case models.IllnessStroke, models.IllnessHeart, models.IllnessSevere:
		return models.VehicleTypeSevere
	case models.IllnessNeonatal:
		return models.VehicleTypeNeonatal
	default:
		return models.VehicleTypeNormal
	}
}

type DispatchCandidate struct {
	Vehicle  *models.Vehicle
	Distance float64
	Score    float64
}

func FindNearestVehicles(lat, lon float64, illnessType models.IllnessType, limit int) []DispatchCandidate {
	requiredType := GetRequiredVehicleType(illnessType)
	vehicles := store.GlobalStore.ListVehiclesByTypeAndStatus(requiredType, models.VehicleStatusStandby)

	var candidates []DispatchCandidate
	for _, v := range vehicles {
		dist := haversine(lat, lon, v.Latitude, v.Longitude)
		score := 1.0 / (dist + 0.1)
		candidates = append(candidates, DispatchCandidate{
			Vehicle:  v,
			Distance: dist,
			Score:    score,
		})
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Distance < candidates[j].Distance
	})

	if len(candidates) > limit {
		candidates = candidates[:limit]
	}

	return candidates
}

func FindNearestVehiclesAllTypes(lat, lon float64, limit int) []DispatchCandidate {
	var allTypes := []models.VehicleType{
		models.VehicleTypeSevere,
		models.VehicleTypeNormal,
		models.VehicleTypeNeonatal,
	}

	var allCandidates := make(map[models.VehicleType][]DispatchCandidate)
	for _, vType := range allTypes {
		vehicles := store.GlobalStore.ListVehiclesByTypeAndStatus(vType, models.VehicleStatusStandby)
		var candidates []DispatchCandidate
		for _, v := range vehicles {
			dist := haversine(lat, lon, v.Latitude, v.Longitude)
			score := 1.0 / (dist + 0.1)
			candidates = append(candidates, DispatchCandidate{
				Vehicle:  v,
				Distance: dist,
				Score:    score,
			})
		}
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Distance < candidates[j].Distance
		})
		if len(candidates) > limit {
			candidates = candidates[:limit]
		}
		allCandidates[vType] = candidates
	}

	var result []DispatchCandidate
	for _, vType := range allTypes {
		result = append(result, allCandidates[vType]...)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Distance < result[j].Distance
	})

	if len(result) > limit*2 {
		result = result[:limit*2]
	}

	return result
}

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	return haversine(lat1, lon1, lat2, lon2)
}
