package types

import "math/big"

type BondProviderEvent struct {
	Pubkey       string
	Chain        string
	TxID         string
	BondRelative *big.Int
	BondAbsolute *big.Int
}

/*
sdk.NewAttribute("pubkey", contract.ProviderPubKey.String()),

	sdk.NewAttribute("chain", contract.Chain.String()),
	sdk.NewAttribute("client", contract.Client.String()),
	sdk.NewAttribute("delegate", contract.Delegate.String()),
	sdk.NewAttribute("type", contract.Type.String()),
	sdk.NewAttribute("height", strconv.FormatInt(contract.Height, 10)),
	sdk.NewAttribute("duration", strconv.FormatInt(contract.Duration, 10)),
	sdk.NewAttribute("rate", strconv.FormatInt(contract.Rate, 10)),
	sdk.NewAttribute("open_cost", strconv.FormatInt(openCost, 10)),
*/

type ContractType string

var (
	ContractTypePayAsYouGo   ContractType = "PayAsYouGo"
	ContractTypeSubscription ContractType = "Subscription"
)

// type Contract struct {
// 	ProviderPubkey string
// 	Chain          string
// 	ClientPubkey   string
// 	DelegatePubkey string
// 	TxID           string
// 	ContractType   ContractType
// 	Height         int64
// 	Duration       int64
// 	Rate           int64
// 	OpenCost       int64
// }

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
	ProviderSortKeyAge           ProviderSortKey = "age"
	ProviderSortKeyContractCount ProviderSortKey = "contract_count"
	ProviderSortKeyAmountPaid    ProviderSortKey = "anount_paid"
)

type ProviderSearchParams struct {
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
