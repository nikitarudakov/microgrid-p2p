package matching

import (
	"github.com/nikitarudakov/microenergy/internal/services/competition"
	"github.com/nikitarudakov/microenergy/internal/services/inventory"
	"math"
	"slices"
)

const EarthRadiusMeters = 6371000 // Earth's radius in meters

func Filter(comp *competition.Competition, assets []*inventory.Asset) []*inventory.Asset {
	eligibleByVoltage := filter(comp, assets, func(asset *inventory.Asset, comp *competition.Competition) bool {
		return isWithinVoltageLimit(asset.VoltageLevel, comp.MinVoltage, comp.MaxVoltage)
	})

	eligibleByDistance := filter(comp, eligibleByVoltage, func(asset *inventory.Asset, comp *competition.Competition) bool {
		return isWithinRadius(asset.Latitude, asset.Longitude, comp.Latitude, comp.Longitude, comp.Radius)
	})

	eligibleByAnyWindow := filter(comp, eligibleByDistance, func(asset *inventory.Asset, comp *competition.Competition) bool {
		return isEligibleForAnyWindow(asset, comp.Windows)
	})

	return eligibleByAnyWindow
}

func isEligibleForAnyWindow(asset *inventory.Asset, windows []*competition.Window) bool {
	index := slices.IndexFunc(windows, func(window *competition.Window) bool {
		return isEligibleServiceWindow(asset, window)
	})

	return !(index == -1)
}

func isEligibleServiceWindow(asset *inventory.Asset, window *competition.Window) bool {
	if asset.MaxCapacity < window.Capacity {
		return false
	}
	if asset.MaxRuntimeMinutes < window.MinRuntimeMinutes {
		return false
	}
	return true
}

func isWithinVoltageLimit(aVolt, compMin, compMax float32) bool {
	return compMin <= aVolt && aVolt <= compMax
}

func isWithinRadius(aLat, aLon, compLat, compLon float64, radius float64) bool {
	dist := haversineDistance(aLat, aLon, compLat, compLon)
	return dist <= radius
}

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
