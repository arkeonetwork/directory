
#!/bin/bash

alicekey=arkeopub1addwnpepqtmsfstktwjyrthala365l9racv799rulta8376ytf4pls3043wvyvlgnn5

bobkey=arkeopub1addwnpepqfsnue7rdp5qppxg6kfj8kn96w6apsff70adeqc9hyw57lk693mvukvnawe

sleep 3

# bond provider
arkeod tx arkeo bond-provider $alicekey btc-mainnet 200000000 --from alice -y

sleep 5

# modify provider
arkeod tx arkeo mod-provider $alicekey btc-mainnet "https://dev.arkeo.network/provider1/metadata.json" 36 1 30 365 10 20 --from alice -y

sleep 5

# open paygo contract
#arkeod tx arkeo open-contract $alicekey btc-mainnet $bobkey 1 60 30 20 --from bob -y

sleep 5

# open subscription contract
arkeod tx arkeo open-contract $alicekey btc-mainnet $bobkey 0 300 30 10 --from bob -y

echo "done"
