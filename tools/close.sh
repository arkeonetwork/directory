#!/bin/bash

BASEDIR=$(dirname "$0")
source $BASEDIR/env.sh

CHAIN=gaia-mainnet-rpc-archive

USER=bob
PROVIDER_PUBKEY=$alicekey
CLIENT_PUBKEY=$bobkey

CONTRACT_ID=$(curl -s $ARKEOD_HOST_LCD/arkeo/active-contract/$CLIENT_PUBKEY/$PROVIDER_PUBKEY/$CHAIN | jq -r .contract.id)

echo "Using contractID $CONTRACT_ID"

arkeod tx arkeo close-contract --from $USER -y $CONTRACT_ID

echo "done"
