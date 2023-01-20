#!/bin/bash

BASEDIR=$(dirname "$0")
source $BASEDIR/env.sh

# increment each invocation
NONCE=$1
if [ -z $NONCE ]
then
  NONCE=1
fi
echo "using nonce $NONCE"
CHAIN=btc-mainnet-fullnode
CTYPE=1 # contract type, 0 is subscription, 1 is pay-as-you-go
DEPOSIT=1000 # amount of tokens you want to deposit. Subscriptions should make sense in that duration and rate equal deposit
DURATION=200 # number of blocks to make a subscription. There are lower and higher limits to this number
RATE=100 # should equal the porvider's rate which you can lookup at (`curl http://seed.arkeo.network:3636/metadata.json | jq .`)
FROMPUBKEY=arkeopub1addwnpepq2tjwwcpqwswatymx7uw3q75sqhljmp0qw3pz2fvnavkv7k2jvknq9k9lr0
USER=adam

# USER=bob
PROVIDER_PUBKEY=arkeopub1addwnpepqtrc0rrpkwn2esula68zl3dvqqfxfjhr5dyfxy3uq97dssntrq8twhy9nvu
CLIENT_PUBKEY=arkeopub1addwnpepq2tjwwcpqwswatymx7uw3q75sqhljmp0qw3pz2fvnavkv7k2jvknq9k9lr0

HEIGHT=$(curl -s $ARKEOD_HOST:1317/arkeo/contract/"$PROVIDER_PUBKEY"/"$CHAIN"/"$CLIENT_PUBKEY" | jq -r .contract.height)

echo "Using height $HEIGHT"
SIGNATURE=$(signhere -u "$USER" -m "$PROVIDER_PUBKEY:$CHAIN:$CLIENT_PUBKEY:$HEIGHT:$NONCE")
echo "SIGNATURE: $SIGNATURE"
echo "executing arkeod tx arkeo claim-contract-income --from $USER $PROVIDER_PUBKEY $CHAIN $CLIENT_PUBKEY $NONCE $HEIGHT $SIGNATURE"

arkeod tx arkeo claim-contract-income -y --from $USER -- "$PROVIDER_PUBKEY" "$CHAIN" "$CLIENT_PUBKEY" "$NONCE" "$HEIGHT" "$SIGNATURE"

echo "done"
