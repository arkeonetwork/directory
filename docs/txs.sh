
#!/bin/bash

alicekey=arkeopub1addwnpepq0dvgx6z3phvw3n2v3yfmpyqpvjux2uz4w6872ekxu3k8zg6gj8ps7qd26k
bobkey=arkeopub1addwnpepqdcpy9afcq3wmas7xd436f53yyfm7v0ggyfck550759a9378ucvg7hkkcvw

# bond provider
arkeod tx arkeo bond-provider $alicekey btc-mainnet-fullnode 100000000 --from alice -y
sleep 6


arkeod tx arkeo bond-provider $alicekey eth-mainnet-fullnode 100000000 --from alice -y
sleep 6

# modify provider
# Error: accepts 9 arg(s), received 0
# Usage:
#   arkeod tx arkeo mod-provider [pubkey] [chain] [metatadata-uri] [metadata-nonce] [status] [min-contract-duration] [max-contract-duration] [subscription-rate] [pay-as-you-go-rate] [flags]
arkeod tx arkeo mod-provider $alicekey btc-mainnet-fullnode "https://dev.arkeo.network/provider1/metadata-btc.json" 0 1 10 5256000 10 20 --from alice -y
sleep 6

arkeod tx arkeo mod-provider $alicekey eth-mainnet-fullnode "https://dev.arkeo.network/provider1/metadata-eth.json" 0 1 10 5256000 10 20 --from alice -y
sleep 6

# sleep 5

# open paygo contract
# Error: accepts 7 arg(s), received 0
# Usage:
#   arkeod tx arkeo open-contract [provider_pubkey] [chain] [client_pubkey] [c-type] [deposit] [duration] [rate] [delegation-optional] [flags]
# arkeod tx arkeo open-contract $alicekey btc-mainnet-fullnode $bobkey 1 60 100 20 --from bob -y

# sleep 5

# open subscription contract
arkeod tx arkeo open-contract $alicekey btc-mainnet-fullnode $bobkey 0 300 30 10 --from bob -y

# close contract
# arkeod tx arkeo close-contract $alicekey btc-mainnet-fullnode $bobkey --from alice -y

echo "done"
