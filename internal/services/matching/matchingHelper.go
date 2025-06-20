package matching

import (
	"github.com/nikitarudakov/microenergy/internal/services/competition"
	"github.com/nikitarudakov/microenergy/internal/services/inventory"
	"math"
)

const EarthRadiusMeters = 6371000 // Earth's radius in meters

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	toRadians := func(deg float64) float64 {
		return deg * math.Pi / 180
	}

	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)

	lat1 = toRadians(lat1)
	lat2 = toRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusMeters * c
}

func filter(
	comp *competition.Competition,
	assets []*inventory.Asset,
	f func(asset *inventory.Asset, comp *competition.Competition) bool,
) (ret []*inventory.Asset) {
	for _, asset := range assets {
		if f(asset, comp) {
			ret = append(ret, asset)
		}
	}
	return
}
