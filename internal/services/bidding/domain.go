package bidding

// Voltage represent voltage requirements for a competition
type Voltage struct {
	Min float32 `json:"min"`
	Max float32 `json:"max"`
}

// ServiceWindow represents a flexibility time window
type ServiceWindow struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
}

// Competition represents published contests by consumers (SOs)
// All of them have bids that were put by FSPs
type Competition struct {
	ConsumerName   string          `json:"consumer_name"`
	Capacity       float32         `json:"capacity"`
	Lifespan       int32           `json:"lifespan"`
	Voltage        Voltage         `json:"voltage"`
	ServiceWindows []ServiceWindow `json:"service_windows"`
	Bids           []Bid           `json:"bids"`
}

// PricingType is used to describe a pricing used in the FSP's bid
type PricingType struct {
	Type  string  `json:"type"`
	Value float32 `json:"value"`
}

// Bid represents a FSPs bids for a competition.
// It includes all fields that are needed for participation in the bidding for the competition.
type Bid struct {
	// Status
	Status bool `json:"status"` // rejected or accepted

	// Details
	Capacity      float32       `json:"capacity"`
	MaxRuntime    int64         `json:"max_runtime"`
	Pricing       []PricingType `json:"pricing"`
	ServiceWindow ServiceWindow `json:"service_window"`

	// Associations
	AssetRef     string       `json:"asset_ref"`
	ProviderName string       `json:"provider_name"`
	Competition  *Competition `json:"competition"`
}
