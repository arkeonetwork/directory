package types

import "math/big"

type BondProviderEvent struct {
	Pubkey       string
	Chain        string
	Height       int64
	TxID         string
	BondRelative *big.Int
	BondAbsolute *big.Int
}

type ContractType string

var (
	ContractTypePayAsYouGo   ContractType = "PayAsYouGo"
	ContractTypeSubscription ContractType = "Subscription"
)

type OpenContractEvent struct {
	ProviderPubkey string
	Chain          string
	ClientPubkey   string
	DelegatePubkey string
	TxID           string
	ContractType   ContractType
	Height         int64
	Duration       int64
	Rate           int64
	OpenCost       int64
}

type ProviderStatus string

var (
	ProviderStatusOnline  ProviderStatus = "Online"
	ProviderStatusOffline ProviderStatus = "Offline"
)

type ModProviderEvent struct {
	Pubkey              string
	Chain               string
	Height              int64
	TxID                string
	MetadataURI         string
	MetadataNonce       uint64
	Status              ProviderStatus
	MinContractDuration int64
	MaxContractDuration int64
	SubscriptionRate    int64
	PayAsYouGoRate      int64
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type ProviderSortKey string

var (
	ProviderSortKeyNone          ProviderSortKey = ""
	ProviderSortKeyAge           ProviderSortKey = "age"
	ProviderSortKeyContractCount ProviderSortKey = "contract_count"
	ProviderSortKeyAmountPaid    ProviderSortKey = "anount_paid"
)

type ProviderSearchParams struct {
	Pubkey                    string
	Chain                     string
	SortKey                   ProviderSortKey
	MaxDistance               int64
	IsMaxDistanceSet          bool
	Coordinates               Coordinates
	MinValidatorPayments      int64
	IsMinValidatorPaymentsSet bool
	MinProviderAge            int64
	IsMinProviderAgeSet       bool
	MinRateLimit              int64
	IsMinRateLimitSet         bool
	MinOpenContracts          int64
	IsMinOpenContractsSet     bool
}
