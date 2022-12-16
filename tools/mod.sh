#!/bin/bash

BASEDIR=$(dirname "$0")
source $BASEDIR/env.sh

CHAIN=eth-mainnet-fullnode

USER=alice
PROVIDER_PUBKEY=$alicekey
METAURL="https://dev.arkeo.network/provider1/metadata-btc.json"
# modify provider
# Error: accepts 9 arg(s), received 0
# Usage:
#   arkeod tx arkeo mod-provider [pubkey] [chain] [metatadata-uri] [metadata-nonce] [status] [min-contract-duration] [max-contract-duration] [subscription-rate] [pay-as-you-go-rate] [flags]
arkeod tx arkeo mod-provider $PROVIDER_PUBKEY $CHAIN $METAURL 0 1 3 5256000 10 20 --from alice -y

echo "done"
