package indexer

import (
	"fmt"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
)

const (
	bech32PrefixAccAddr = "rko"
	bech32PrefixAccPub  = "rkopub"
)

// list keys in $HOME/.arkeo/keyring-test with addr and pubkey
func TestAccountDetails(t *testing.T) {
	arkeoHome := fmt.Sprintf("%s/.arkeo", os.Getenv("HOME"))
	addr, pubkey, err := accountDetails("alice", arkeoHome)
	if err != nil {
		t.Errorf("error getting details for alice: %+v", err)
	}
	log.Infof("alice: addr: %s pubkey: %s", addr, pubkey)

	addr, pubkey, err = accountDetails("bob", arkeoHome)
	if err != nil {
		t.Errorf("error getting details for bob: %+v", err)
	}
	log.Infof("bob: addr: %s pubkey: %s", addr, pubkey)
}

func accountDetails(keyName, keyringPath string) (addr, pubkey string, err error) {
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(bech32PrefixAccAddr, bech32PrefixAccPub)
	encConfig := NewEncoding()

	keyRing, err := keyring.New("arkeo", "test", keyringPath, nil, encConfig.Marshaler)
	if err != nil {
		log.Errorf("error opening keyring: %+v", err)
		return
	}
	all, err := keyRing.List()
	if err != nil {
		log.Errorf("error listing keys: %+v", err)
		return
	}
	for _, v := range all {
		pub, perr := v.GetPubKey()
		if perr != nil {
			log.Errorf("error getting \"%s\" pubkey from keyring: %+v", v.Name, err)
			continue
		}
		accAddr := sdk.AccAddress(pub.Address())
		addr, err = bech32.ConvertAndEncode(sdkConfig.GetBech32AccountAddrPrefix(), accAddr)
		if err != nil {
			log.Errorf("error encoding account address %+v", err)
			continue
		}
		pubkey, err = bech32.ConvertAndEncode(sdkConfig.GetBech32AccountPubPrefix(), legacy.Cdc.MustMarshal(pub))
		if err != nil {
			log.Errorf("error encoding pubkey %+v", err)
			return
		}
	}

	info, err := keyRing.Key(keyName)
	if err != nil {
		log.Errorf("error getting alice key: %+v", err)
		return
	}

	pub, err := info.GetPubKey()
	if err != nil {
		log.Errorf("error getting pubkey from keyring: %+v", err)
		return
	}

	accAddr := sdk.AccAddress(pub.Address())
	addr, err = bech32.ConvertAndEncode(sdkConfig.GetBech32AccountAddrPrefix(), accAddr)
	if err != nil {
		log.Errorf("error encoding account address %+v", err)
		return
	}

	pubkey, err = bech32.ConvertAndEncode(sdkConfig.GetBech32AccountPubPrefix(), legacy.Cdc.MustMarshal(pub))
	if err != nil {
		log.Errorf("error encoding pubkey %+v", err)
		return
	}

	return
}
