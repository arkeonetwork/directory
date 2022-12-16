package types

type BondProviderEvent struct {
	Pubkey       string `mapstructure:"pubkey"`
	Chain        string `mapstructure:"chain"`
	Height       int64  `mapstructure:"height"`
	TxID         string `mapstructure:"hash"`
	BondRelative string `mapstructure:"bond_rel"`
	BondAbsolute string `mapstructure:"bond_abs"`
}

type ContractType string

var (
	ContractTypePayAsYouGo   ContractType = "PayAsYouGo"
	ContractTypeSubscription ContractType = "Subscription"
)

type BaseContractEvent struct {
	ProviderPubkey string `mapstructure:"pubkey"`
	Chain          string `mapstructure:"chain"`
	ClientPubkey   string `mapstructure:"client"`
	DelegatePubkey string `mapstructure:"delegate"` // see BaseContractEvent.GetDelegatePubkey()
	TxID           string `mapstructure:"hash"`
	Height         int64  `mapstructure:"height"`
	EventHeight    int64  `mapstructure:"eventHeight"`
}

// get the delegate pubkey falling back to client pubkey if undefined
func (b BaseContractEvent) GetDelegatePubkey() string {
	if b.DelegatePubkey != "" {
		return b.DelegatePubkey
	}
	return b.ClientPubkey
}

type OpenContractEvent struct {
	BaseContractEvent `mapstructure:",squash"`
	Duration          int64        `mapstructure:"duration"`
	ContractType      ContractType `mapstructure:"type"`
	Rate              int64        `mapstructure:"rate"`
	OpenCost          int64        `mapstructure:"open_cost"`
}

type ContractSettlementEvent struct {
	BaseContractEvent `mapstructure:",squash"`
	Nonce             string `mapstructure:"nonce"`
	Paid              string `mapstructure:"paid"`
	Reserve           string `mapstructure:"reserve"`
}

type CloseContractEvent struct {
	ContractSettlementEvent `mapstructure:",squash"`
}

type ValidatorPayoutEvent struct {
	Validator string `mapstructure:"validator"`
	Height    int64  `mapstructure:"height"`
	TxID      string `mapstructure:"hash"`
	Paid      int64  `mapstructure:"paid"`
}

type ProviderStatus string

var (
	ProviderStatusOnline  ProviderStatus = "Online"
	ProviderStatusOffline ProviderStatus = "Offline"
)

type ModProviderEvent struct {
	Pubkey              string         `mapstructure:"pubkey"`
	Chain               string         `mapstructure:"chain"`
	Height              int64          `mapstructure:"height"`
	TxID                string         `mapstructure:"hash"`
	MetadataURI         string         `mapstructure:"metadata_uri"`
	MetadataNonce       uint64         `mapstructure:"metadata_nonce"`
	Status              ProviderStatus `mapstructure:"status"`
	MinContractDuration int64          `mapstructure:"min_contract_duration"`
	MaxContractDuration int64          `mapstructure:"max_contract_duration"`
	SubscriptionRate    int64          `mapstructure:"subscription_rate"`
	PayAsYouGoRate      int64          `mapstructure:"pay-as-you-go_rate"`
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
