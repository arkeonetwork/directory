#!/bin/bash

BASEDIR=$(dirname "$0")
source $BASEDIR/env.sh

CHAIN=eth-mainnet-fullnode

USER=alice
PROVIDER_PUBKEY=$alicekey
AMT=100000000

arkeod tx arkeo bond-provider --from $USER -y $PROVIDER_PUBKEY $CHAIN $AMT

echo "done"
