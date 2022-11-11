package indexer

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	stdtypes "github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	// ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	// ibccoretypes "github.com/cosmos/ibc-go/v3/modules/core/types"
)

type encodingConfig struct {
	InterfaceRegistry types.InterfaceRegistry
	Marshaler         codec.Codec
	TxConfig          client.TxConfig
	Amino             *codec.LegacyAmino
}

// NewEncoding registers all base protobuf types by default as well as any custom types passed in
func NewEncoding(registerInterfaces ...func(r types.InterfaceRegistry)) *encodingConfig {
	registry := types.NewInterfaceRegistry()

	// register base protobuf types
	authztypes.RegisterInterfaces(registry)
	banktypes.RegisterInterfaces(registry)
	distributiontypes.RegisterInterfaces(registry)
	// ibccoretypes.RegisterInterfaces(registry)
	// ibctransfertypes.RegisterInterfaces(registry)
	stakingtypes.RegisterInterfaces(registry)
	stdtypes.RegisterInterfaces(registry)

	// register custom protobuf types
	for _, r := range registerInterfaces {
		r(registry)
	}

	marshaler := codec.NewProtoCodec(registry)

	return &encodingConfig{
		InterfaceRegistry: registry,
		Marshaler:         marshaler,
		TxConfig:          tx.NewTxConfig(marshaler, tx.DefaultSignModes),
		Amino:             codec.NewLegacyAmino(),
	}
}

const mnemonic = "blade soda fish scale custom thumb foam garden boil enter stage cover spatial nation alert shield witness predict shaft harbor grant inmate ketchup tiger"

// configure
func configure() {
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	encConfig := NewEncoding()

	// doing this so it dumps the seeds addr and pubkey for use in arkeo txs like bond provider
	keyRing := keyring.NewInMemory(encConfig.Marshaler)
	info, err := keyRing.NewAccount("arkeo-directory", mnemonic, "", "m/44'/118'/0'/0/0", hd.Secp256k1)
	if err != nil {
		log.Errorf("error create account from mnemonic: %+v", err)
		return
	}

	pub, err := info.GetPubKey()
	if err != nil {
		log.Errorf("error getting pubkey from keyring: %+v", err)
		return
	}

	accAddr := sdk.AccAddress(pub.Address())
	addr, err := bech32.ConvertAndEncode(sdkConfig.GetBech32AccountAddrPrefix(), accAddr)
	if err != nil {
		log.Errorf("error encoding account address %+v", err)
		return
	}
	log.Infof("address: %s", addr)

	pubkey, err := bech32.ConvertAndEncode(sdkConfig.GetBech32AccountPubPrefix(), legacy.Cdc.MustMarshal(pub))
	if err != nil {
		log.Errorf("error encoding pubkey %+v", err)
		return
	}
	log.Infof("pubkey: %s", pubkey)
}
