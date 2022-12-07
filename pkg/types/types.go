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

type OpenContractEvent struct {
	ProviderPubkey string       `mapstructure:"pubkey"`
	Chain          string       `mapstructure:"chain"`
	ClientPubkey   string       `mapstructure:"client"`
	DelegatePubkey string       `mapstructure:"client"`
	TxID           string       `mapstructure:"txID"`
	ContractType   ContractType `mapstructure:"type"`
	Height         int64        `mapstructure:"height"`
	Duration       int64        `mapstructure:"duration"`
	Rate           int64        `mapstructure:"rate"`
	OpenCost       int64        `mapstructure:"open_cost"`
}

type ContractSettlementEvent struct {
	ProviderPubkey string `mapstructure:"pubkey"`
	Chain          string `mapstructure:"chain"`
	ClientPubkey   string `mapstructure:"client"`
	Paid           string `mapstructure:"paid"`
	Height         string `mapstructure:"height"`
	Nonce          string `mapstructure:"nonce"`
	Reserve        string `mapstructure:"reserve"`
}

type ValidatorPayoutEvent struct {
	Validator string `mapstructure:"validator"`
	Paid      string `mapstructure:"paid"`
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
