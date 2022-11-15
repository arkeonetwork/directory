package types

import "math/big"

type BondProviderEvent struct {
	Pubkey       string
	Chain        string
	BondRelative *big.Int
	BondAbsolute *big.Int
}

type ProviderStatus string

var (
	ProviderStatusOnline  ProviderStatus = "Online"
	ProviderStatusOffline ProviderStatus = "Offline"
)

type ModProviderEvent struct {
	Pubkey              string
	Chain               string
	MetadataURI         string
	MetadataNonce       uint64
	Status              ProviderStatus
	MinContractDuration int64
	MaxContractDuration int64
	SubscriptionRate    int64
	PayAsYouGoRate      int64
}

type Coordinates struct {
	Latitude  float32
	Longitude float32
}
