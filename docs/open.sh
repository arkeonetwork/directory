#!/bin/bash

source ./env.sh

# increment each invocation
NONCE=10
CHAIN=eth-mainnet-fullnode

USER=bob
PROVIDER_PUBKEY=$alicekey
CLIENT_PUBKEY=$bobkey

arkeod tx arkeo open-contract --from $USER $PROVIDER_PUBKEY eth-mainnet-fullnode $CLIENT_PUBKEY 1 60 100 20  -y
echo "done"
