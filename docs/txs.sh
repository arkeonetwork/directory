
#!/bin/bash

source ./env.sh

# bond provider
arkeod tx arkeo bond-provider $alicekey btc-mainnet-fullnode 100000000 --from alice -y
sleep 6


arkeod tx arkeo bond-provider $alicekey eth-mainnet-fullnode 100000000 --from alice -y
sleep 6

# modify provider
# Error: accepts 9 arg(s), received 0
# Usage:
#   arkeod tx arkeo mod-provider [pubkey] [chain] [metatadata-uri] [metadata-nonce] [status] [min-contract-duration] [max-contract-duration] [subscription-rate] [pay-as-you-go-rate] [flags]
arkeod tx arkeo mod-provider $alicekey btc-mainnet-fullnode "https://dev.arkeo.network/provider1/metadata-btc.json" 1 1 10 5256000 10 20 --from alice -y
sleep 6

arkeod tx arkeo mod-provider $alicekey eth-mainnet-fullnode "https://dev.arkeo.network/provider1/metadata-eth.json" 1 1 10 5256000 10 20 --from alice -y
sleep 6

# open paygo contract
# Error: accepts 7 arg(s), received 0
# Usage:
#   arkeod tx arkeo open-contract [provider_pubkey] [chain] [client_pubkey] [c-type] [deposit] [duration] [rate] [delegation-optional] [flags]
# arkeod tx arkeo open-contract $alicekey btc-mainnet-fullnode $bobkey 1 60 10 20 --from bob -y
arkeod tx arkeo open-contract $alicekey eth-mainnet-fullnode $bobkey 1 60 100 20 --from bob -y

sleep 6

# open subscription contract
arkeod tx arkeo open-contract $alicekey btc-mainnet-fullnode $bobkey 0 3000 300 10 --from bob -y

# close contract
# arkeod tx arkeo close-contract $alicekey btc-mainnet-fullnode $bobkey --from alice -y

# claim contract income see claim.sh
# Error: accepts 6 arg(s), received 0
# Usage:
#   arkeod tx arkeo claim-contract-income [pubkey] [chain] [client] [nonce] [height] [signature] [flags]


echo "done"
