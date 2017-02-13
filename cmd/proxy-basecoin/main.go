/*
proxy-basecoin is an example script that sets up a proxy
that knows about the basic basecoin types (sendtx and accounts).

It can be extended to add support for basecoin plugins,
or as a source of inspiration to configure a proxy for your
own app.

If you run the basecoin demo app locally (from the data dir),
try the following:

proxy-basecoin -chain=test_chain_id -rpc=localhost:46657

curl http://localhost:8108/keys/
curl -XPOST http://localhost:8108/keys/ -d '{"name": "alice", "passphrase": "1234567890"}'
curl -XPOST http://localhost:8108/keys/ -d '{"name": "bobby", "passphrase": "1234567890"}'
curl http://localhost:8108/keys/

# -> at this point, grab that address, but it in the genesis for
# basecoin, so you are rich, and restart the basecoin server ;)

## TODO: working examples here

# query no data
curl http://localhost:8108/query/store/01234567

# get an account (by path, or knowing the special prefix)
curl http://localhost:8108/query/account/1B1BE55F969F54064628A63B9559E7C21C925165
curl http://localhost:8108/query/store/626173652f612f1B1BE55F969F54064628A63B9559E7C21C925165

# 626173652f612f <- this is the magic base/a/ prefix for accounts in hex
# 1B1BE55F969F54064628A63B9559E7C21C925165 <- address with coins

# get proof by complete key
# TODO: currently fails, complaining about validator sigs
curl http://localhost:8108/proof/626173652f612f1B1BE55F969F54064628A63B9559E7C21C925165

# post a tx (not yet implemented)
# use addressed returned from your keys call above
# input is alice, output is bob
curl -XPOST http://localhost:8108/txs/ -d '{
  "name": "alice",
  "passphrase": "1234567890",
  "data": {
    "gas": 22,
    "fee": {"denom": "ETH", "amount": 1},
    "inputs": [{
      "address": "4d8908785ec867139ca02e71a717c01fa506b96a",
      "coins": [{"denom": "ETH", "amount": 21}],
      "sequence": 1,
      "pub_key": [1, "d7fb176319af0c126c4c4c7851cf7c66340e7df8763f0aa9700ebae32a955401"]
    }],
    "outputs": [{
      "address": "9f31a3ac6b1468402aac5593ae9e82a041457f5f",
      "coins": [{"denom": "ETH", "amount": 20}]
    }]
  }
}'

*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tendermint/light-client/cryptostore"
	"github.com/tendermint/light-client/extensions/basecoin"
	"github.com/tendermint/light-client/proxy"
	"github.com/tendermint/light-client/rpc"
	"github.com/tendermint/light-client/storage/filestorage"
)

var (
	port    = flag.Int("port", 8108, "port for proxy server")
	rpcAddr = flag.String("rpc", "", "url for tendermint rpc server")
	chainID = flag.String("chain", "", "id of the blockchain")
	keydir  = flag.String("keydir", ".keys", "directory to store private keys")
)

// TODO: add cors and unix-socket support like over in verifier
func main() {
	flag.Parse()

	// TODO: make these actually do something
	vr := basecoin.BasecoinValues{}
	sr := basecoin.BasecoinTx{ChainID: *chainID}

	if *chainID == "" {
		fmt.Println("You must specify -chain with the chain_id")
		return
	}
	if *rpcAddr == "" {
		fmt.Println("You must specify -rpc with the location of a tendermint node")
		return
	}

	// set up all the pieces based on config
	r := mux.NewRouter()
	cstore := cryptostore.New(
		cryptostore.GenEd25519,
		cryptostore.SecretBox,
		filestorage.New(*keydir),
	)
	node := rpc.NewNode(*rpcAddr, *chainID, vr)
	proxy.RegisterDefault(r, cstore, node, sr, vr)

	// TODO: add some awesome middlewares...

	// now, just run the server and bind to localhost for security
	bind := fmt.Sprintf("127.0.0.1:%d", *port)
	log.Fatal(http.ListenAndServe(bind, r))

}
