package db

import (
	"testing"

	"github.com/google/uuid"
)

func TestInsertProvider(t *testing.T) {
	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}

	entity, err := db.InsertProvider(&ArkeoProvider{
		Pubkey: uuid.NewString(),
		Chain:  "bitcoin-mainnet",
		Bond:   "1234567890",
	})
	if err != nil {
		t.Errorf("error inserting provider: %+v", err)
		t.FailNow()
	}
	log.Infof("inserted provider %d", entity.ID)
}

func TestFindProvider(t *testing.T) {
	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}
	pubkey := "arkeopub1addwnpepqg5fsc756nx3wlrp7f4328slhgfulhu53epxnyy4q6ln3htrhxxsczgwfyf"
	chain := "btc-mainnet"
	provider, err := db.FindProvider(pubkey, chain)
	if err != nil {
		t.Errorf("error finding provider: %+v", err)
		t.FailNow()
	}
	log.Infof("found provider %d", provider.ID)

	pubkey = "nosuchthing"
	provider, err = db.FindProvider(pubkey, chain)
	if err != nil {
		t.Errorf("error finding provider: %+v", err)
		t.FailNow()
	}
	if provider != nil {
		t.Errorf("expected nil but got %v", provider)
	}
}
