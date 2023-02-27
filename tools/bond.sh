#!/bin/bash

BASEDIR=$(dirname "$0")
source $BASEDIR/env.sh

# CHAIN=eth-mainnet-fullnode
CHAIN=gaia-mainnet-rpc-archive

USER=alice
PROVIDER_PUBKEY=$alicekey
AMT=100000000000

arkeod --node $ARKEOD_HOST tx arkeo bond-provider --from $USER -y $PROVIDER_PUBKEY $CHAIN $AMT
