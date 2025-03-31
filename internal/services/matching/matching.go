package matching

import (
	"github.com/google/uuid"
	"github.com/nikitarudakov/microenergy/internal/services/competition"
	"github.com/nikitarudakov/microenergy/internal/services/inventory"
	"math"
	"math/rand"
	"slices"
)

type Match struct {
	Window *competition.Window
	Asset  *inventory.Asset
	Score  float64
}

type Service struct{}

func (s *Service) Filter(comp *competition.Competition, assets []*inventory.Asset) []*inventory.Asset {
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

func (s *Service) Match(comp *competition.Competition, assets []*inventory.Asset) []*Match {
	filteredAssets := s.Filter(comp, assets)

	usedAssets := map[uuid.UUID]bool{}
	var matches []*Match

	for _, window := range comp.Windows {
		var bestMatch *Match
		var bestScore float64

		for _, asset := range filteredAssets {
			// Skip if already used
			if usedAssets[asset.ID] {
				continue
			}
			if !isEligibleServiceWindow(asset, window) {
				continue
			}

			score := s.Score(asset, window, comp)
			if bestMatch == nil || score > bestScore {
				bestMatch = &Match{Window: window, Asset: asset, Score: score}
				bestScore = score
			}
		}

		if bestMatch != nil {
			matches = append(matches, bestMatch)
			usedAssets[bestMatch.Asset.ID] = true // Mark asset as used
		}
	}

	return matches
}

func (s *Service) Score(asset *inventory.Asset, window *competition.Window, comp *competition.Competition) float64 {
	// Distance score (closer = better)
	dist := haversineDistance(asset.Latitude, asset.Longitude, comp.Latitude, comp.Longitude)
	distScore := 1.0 - math.Min(dist/comp.Radius, 1.0) // normalized 0.0–1.0

	// Capacity fit: how close the asset is to the required capacity
	capacityRatio := float64(window.Capacity) / float64(asset.MaxCapacity)
	if capacityRatio > 1.0 {
		capacityRatio = 0.0 // can't meet requirement
	}
	capacityScore := capacityRatio // closer to 1.0 = better

	// Runtime fit (how much runtime margin we have)
	runtimeRatio := float64(asset.MaxRuntimeMinutes-window.MinRuntimeMinutes) / float64(asset.MaxRuntimeMinutes)
	runtimeScore := math.Max(0, runtimeRatio)

	// Randomness to break ties
	random := rand.Float64() * 0.05

	return distScore*0.4 + capacityScore*0.4 + runtimeScore*0.15 + random
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
