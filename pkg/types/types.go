package types

import "math/big"

type ProviderBondEvent struct {
	Pubkey       string
	Chain        string
	BondRelative *big.Int
	BondAbsolute *big.Int
}
