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
# CHAIN=eth-mainnet-fullnode
CHAIN=gaia-mainnet-rpc-archive
FROMPUBKEY=$bobkey
USER=bob

PROVIDER_PUBKEY=$alicekey
CLIENT_PUBKEY=$bobkey

CONTRACT_ID=$(curl -s $ARKEOD_HOST_LCD/arkeo/active-contract/$CLIENT_PUBKEY/$PROVIDER_PUBKEY/$CHAIN | jq -r .contract.id)
SIGNATURE=$(signhere -u "$USER" -m "$CONTRACT_ID:$CLIENT_PUBKEY:$NONCE")

echo "executing arkeod tx arkeo claim-contract-income --from $USER $CHAIN $PROVIDER_PUBKEY $CONTRACT_ID  $CLIENT_PUBKEY $NONCE $SIGNATURE"
arkeod tx arkeo claim-contract-income -y --from $USER -- "$CONTRACT_ID" "$CLIENT_PUBKEY" "$NONCE" "$SIGNATURE"

echo "done"
