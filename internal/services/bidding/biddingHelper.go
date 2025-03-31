package bidding

import (
	"github.com/nikitarudakov/microenergy/internal/services/competition"
	"slices"
)

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
