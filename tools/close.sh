#!/bin/bash

BASEDIR=$(dirname "$0")
source $BASEDIR/env.sh

CHAIN=gaia-mainnet-rpc-archive

USER=bob
PROVIDER_PUBKEY=$alicekey
CLIENT_PUBKEY=$bobkey

# Error: accepts 3 arg(s), received 0
# Usage:
#   arkeod tx arkeo close-contract [pubkey] [chain] [client] [delegate-optional] [flags]
arkeod tx arkeo close-contract --from $USER -y $PROVIDER_PUBKEY $CHAIN $CLIENT_PUBKEY

# arkeod tx arkeo open-contract --from $USER $PROVIDER_PUBKEY eth-mainnet-fullnode $CLIENT_PUBKEY 1 100 100 20  -y
echo "done"
