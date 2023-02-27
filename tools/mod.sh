#!/bin/bash

BASEDIR=$(dirname "$0")
source $BASEDIR/env.sh

# CHAIN=eth-mainnet-fullnode
CHAIN=gaia-mainnet-rpc-archive

USER=alice
PROVIDER_PUBKEY=$alicekey
METAURL="http://testnet-sentinel.arkeo.shapeshift.com:3636/metadata.json"
METADATA_NONCE=0
# modify provider
# Error: accepts 9 arg(s), received 0
# Usage:
#   arkeod tx arkeo mod-provider [pubkey] [chain] [metatadata-uri] [metadata-nonce] [status] [min-contract-duration] [max-contract-duration] [subscription-rate] [pay-as-you-go-rate] [flags]
arkeod --node $ARKEOD_HOST tx arkeo mod-provider $PROVIDER_PUBKEY $CHAIN $METAURL $METADATA_NONCE 1 10 5256000 10 20 --from alice -y

echo "done"
