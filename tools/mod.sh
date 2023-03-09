#!/bin/bash

BASEDIR=$(dirname "$0")
source $BASEDIR/env.sh

# CHAIN=eth-mainnet-fullnode
CHAIN=gaia-mainnet-rpc-archive

USER=alice
PROVIDER_PUBKEY=$alicekey
METADATA_NONCE=0
arkeod tx arkeo mod-provider $PROVIDER_PUBKEY $CHAIN $METAURL $METADATA_NONCE 1 10 5256000 10 20 10 --from alice -y

echo "done"
