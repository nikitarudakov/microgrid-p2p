package bidding

import (
	"github.com/google/uuid"
	"github.com/nikitarudakov/microenergy/internal/services/competition"
	"slices"
)

func isBidEligible(bid *Bid, usedAssets map[uuid.UUID]bool) bool {
	for _, asset := range bid.Assets {
		if usedAssets[asset] {
			return false
		}
	}
	return true
}

func filterBidsForWindow(bids []*Bid, window *competition.Window) []*Bid {
	return slices.DeleteFunc(bids, func(b *Bid) bool {
		return b.ServiceWindowID != window.ID // assuming a way to match
	})
}

func totalBidCost(b *Bid) float64 {
	return b.Capacity*b.AvailabilityPrice + b.Capacity*b.UtilizationPrice + b.ServiceFee
}

func sumWinnerPrices(winners []*Winner) float64 {
	var sum float64
	for _, w := range winners {
		sum += w.TotalPrice
	}
	return sum
}

func markAssetsUsed(bid *Bid, used map[uuid.UUID]bool) {
	for _, assetID := range bid.Assets {
		used[assetID] = true
	}
}
