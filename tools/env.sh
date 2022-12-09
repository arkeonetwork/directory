#!/bin/bash

KEY_RAW=$(arkeod keys show alice -p | jq -r .key)
PUBKEY=$(arkeod debug pubkey-raw "$KEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

export alicekey=$PUBKEY

KEY_RAW=$(arkeod keys show bob -p | jq -r .key)
PUBKEY=$(arkeod debug pubkey-raw "$KEY_RAW" | grep "Bech32 Acc" | awk '{ print $NF }')

export bobkey=$PUBKEY
