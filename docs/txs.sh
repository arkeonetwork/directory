
#!/bin/bash

alicekey=arkeopub1addwnpepqwdy80ejutgdh8g6qwpmsl27p4lzphnga9ndm937k7axggsfeydjyhs6z65

bobkey=arkeopub1addwnpepqd53sqsa2w3d6mkgw4ymjn0xee9fpc8tucejhj0dhsmtfhf0949duvqklz7

sleep 3

# bond provider
arkeod tx arkeo bond-provider $alicekey btc-mainnet 200000000 --from alice -y

sleep 5

# modify provider
arkeod tx arkeo mod-provider $alicekey btc-mainnet "https://dev.arkeo.network/provider1/metadata.json" 36 1 10 365 10 20 --from alice -y

# sleep 5

# open paygo contract
# arkeod tx arkeo open-contract $alicekey btc-mainnet $bobkey 1 60 10 20 --from bob -y

# sleep 5

# open subscription contract
# arkeod tx arkeo open-contract $alicekey btc-mainnet $bobkey 0 300 30 10 --from bob -y

# close contract
# arkeod tx arkeo close-contract $alicekey btc-mainnet $bobkey --from alice -y

echo "done"
