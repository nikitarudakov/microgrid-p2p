package inventory

// Provider represents a flexibility provider identity
type Provider struct {
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

// Asset represents a distributed energy resource
// and describes its technical aspects and services
// it can be used for flexibility provision
type Asset struct {
	// Identifiers
	Ref  string `json:"ref"`
	Name string `json:"name"`

	// Technical details
	Service      string  `json:"service"`
	VoltageLevel float32 `json:"voltage_level"`

	// Location
	Latitude  int16 `json:"latitude"`
	Longitude int16 `json:"Longitude"`

	// Associations
	MeterID      string    `json:"meter_id"`
	ProviderName string    `json:"provider_id"`
	Provider     *Provider `json:"provider"`
}
