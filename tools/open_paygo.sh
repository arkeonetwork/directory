#!/bin/bash

export ARKEOD_KEYRING_BACKEND=test
BASEDIR=$(dirname "$0")
source $BASEDIR/env.sh

CTYPE=1 # contract type, 0 is subscription, 1 is pay-as-you-go
CHAIN=gaia-mainnet-rpc-archive
DEPOSIT=400 # amount of tokens you want to deposit. Subscriptions should make sense in that duration and rate equal deposit
DURATION=20 # number of blocks to make a subscription. There are lower and higher limits to this number
SETTLEMENT_DURATION=10
RATE=20 # should equal the porvider's rate which you can lookup at (`curl http://seed.arkeo.network:3636/metadata.json | jq .`)
FROMPUBKEY=$bobkey
USER=bob
arkeod tx arkeo open-contract -y --from $USER -- $alicekey $CHAIN "$FROMPUBKEY" "$CTYPE" "$DEPOSIT" "$DURATION" $RATE "$SETTLEMENT_DURATION"
