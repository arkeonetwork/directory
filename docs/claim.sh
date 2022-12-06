#!/bin/bash

source ./env.sh

# increment each invocation
NONCE=7
CHAIN=eth-mainnet-fullnode

USER=bob
PROVIDER_PUBKEY=$alicekey
CLIENT_PUBKEY=$bobkey



HEIGHT=$(curl -s localhost:1317/arkeo/contract/"$PROVIDER_PUBKEY"/"$CHAIN"/"$CLIENT_PUBKEY" | jq -r .contract.height)

echo "Using height $HEIGHT"
SIGNATURE=$(signhere -u "$USER" -m "$PROVIDER_PUBKEY:$CHAIN:$CLIENT_PUBKEY:$HEIGHT:$NONCE")
echo "SIGNATURE: $SIGNATURE"
echo "executing arkeod tx arkeo claim-contract-income --from $USER $PROVIDER_PUBKEY $CHAIN $CLIENT_PUBKEY $NONCE $HEIGHT $SIGNATURE"

###NEWONE
arkeod tx arkeo claim-contract-income -y --from $USER -- "$PROVIDER_PUBKEY" "$CHAIN" "$CLIENT_PUBKEY" "$NONCE" "$HEIGHT" "$SIGNATURE"

# arkeod tx arkeo claim-contract-income --from $USER $PUBKEY $CHAIN $CLIENT_PUBKEY $NONCE $HEIGHT $SIGNATURE
echo "done"