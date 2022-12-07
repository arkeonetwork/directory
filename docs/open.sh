#!/bin/bash

BASEDIR=$(dirname "$0")
source $BASEDIR/env.sh

# increment each invocation
CHAIN=btc-mainnet-fullnode

USER=bob
PROVIDER_PUBKEY=$alicekey
CLIENT_PUBKEY=$bobkey

# Error: accepts 7 arg(s), received 0
# Usage:
#   arkeod tx arkeo open-contract [provider_pubkey] [chain] [client_pubkey] [c-type] [deposit] [duration] [rate] [delegation-optional] [flags]
arkeod tx arkeo open-contract --from $USER $PROVIDER_PUBKEY eth-mainnet-fullnode $CLIENT_PUBKEY 1 100 100 20  -y
echo "done"
