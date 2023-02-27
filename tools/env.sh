#!/bin/bash
export ARKEOD_HOST=http://testnet-seed.arkeo.shapeshift.com:26657
export ARKEOD_HOST_LCD=http://testnet-seed.arkeo.shapeshift.com:1317

# export ARKEOD_HOST=http://localhost:26657
# export ARKEOD_HOST_LCD=http://localhost:1317


export ARKEOD_KEYRING_BACKEND=test

KEY_RAW=$(arkeod keys show alice -p | jq -r .key)
PUBKEY=$(arkeod debug pubkey-raw "$KEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

export alicekey=$PUBKEY

KEY_RAW=$(arkeod keys show bob -p | jq -r .key)
PUBKEY=$(arkeod debug pubkey-raw "$KEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

export bobkey=$PUBKEY
