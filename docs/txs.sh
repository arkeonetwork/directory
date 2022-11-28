
#!/bin/bash

alicekey=arkeopub1addwnpepqg0t7ze7s6fahfewl5df9fs50cud977qpr9agdfyv80y3mkqmntcksmpep6
bobkey=arkeopub1addwnpepq020773rsennrvmk2w9sdyxlcaylttun8x6jdmd5l3xtve7tun8zv4hhwwe

sleep 3

# bond provider
arkeod tx arkeo bond-provider $alicekey btc-mainnet-fullnode 200000000 --from alice -y

sleep 5

# modify provider
arkeod tx arkeo mod-provider $alicekey btc-mainnet-fullnode "https://dev.arkeo.network/provider1/metadata.json" 36 1 10 365 10 20 --from alice -y

# sleep 5

# open paygo contract
# arkeod tx arkeo open-contract $alicekey btc-mainnet-fullnode $bobkey 1 60 10 20 --from bob -y

# sleep 5

# open subscription contract
# arkeod tx arkeo open-contract $alicekey btc-mainnet-fullnode $bobkey 0 300 30 10 --from bob -y

# close contract
# arkeod tx arkeo close-contract $alicekey btc-mainnet-fullnode $bobkey --from alice -y

echo "done"
