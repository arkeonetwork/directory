#!/bin/bash

export ARKEOD_KEYRING_BACKEND=test

CTYPE=0 # contract type, 0 is subscription, 1 is pay-as-you-go
DEPOSIT=2000 # amount of tokens you want to deposit. Subscriptions should make sense in that duration and rate equal deposit
DURATION=200 # number of blocks to make a subscription. There are lower and higher limits to this number
RATE=10 # should equal the porvider's rate which you can lookup at (`curl http://seed.arkeo.network:3636/metadata.json | jq .`)
FROMPUBKEY=arkeopub1addwnpepq2tjwwcpqwswatymx7uw3q75sqhljmp0qw3pz2fvnavkv7k2jvknq9k9lr0
USER=adam
arkeod tx arkeo open-contract -y --from $USER --keyring-backend $ARKEOD_KEYRING_BACKEND --node "tcp://seed.arkeo.network:26657" -- arkeopub1addwnpepqtrc0rrpkwn2esula68zl3dvqqfxfjhr5dyfxy3uq97dssntrq8twhy9nvu btc-mainnet-fullnode "$FROMPUBKEY" "$CTYPE" "$DEPOSIT" "$DURATION" $RATE



# CTYPE=1 # contract type, 0 is subscription, 1 is pay-as-you-go
# DEPOSIT=1000 # amount of tokens you want to deposit. Subscriptions should make sense in that duration and rate equal deposit
# DURATION=200 # number of blocks to make a subscription. There are lower and higher limits to this number
# RATE=100 # should equal the porvider's rate which you can lookup at (`curl http://seed.arkeo.network:3636/metadata.json | jq .`)
# FROMPUBKEY=arkeopub1addwnpepq2tjwwcpqwswatymx7uw3q75sqhljmp0qw3pz2fvnavkv7k2jvknq9k9lr0
# USER=adam
# arkeod tx arkeo open-contract -y --from $USER --keyring-backend file --node "tcp://seed.arkeo.network:26657" -- arkeopub1addwnpepqtrc0rrpkwn2esula68zl3dvqqfxfjhr5dyfxy3uq97dssntrq8twhy9nvu btc-mainnet-fullnode "$FROMPUBKEY" "$CTYPE" "$DEPOSIT" "$DURATION" $RATE
