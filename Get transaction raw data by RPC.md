We put much attetion on parsing the binary data of bitcoin which makes our work seems working along. Bitcoin is a distributed network system which means it need nodes around the world to work together,
and even a bitcoin core library needs to connect to other nodes for help, just like quering the transaction fee.

There is no free lanuch at the world. When you using bitcoin system for transaction you need to pay some fee to nodes for helping you process the transaction data. Just like banks need to pay its employee
for services. To get the right amount for transaction fee is not easy and there are some providers named RPC nodes that are providing several services for quering such info.We can use http request to query
the needed info from RPC nodes, following is link for a RPC provider for quering info for testnet:

https://blockstream.info/testnet/api/

and following is the URL for quering info for mainnet:

https://blockstream.info/api

you can get details about the api for the RPC provider by the link:
https://github.com/Blockstream/esplora/blob/master/API.md

Let's see an example, given a hash id we get from raw transaction data in previous section:
d95cc33c3ec74c75ad8a1868c88007f11a04e411ff073996ffbea52d12937319
we can query its info by using the following link:

https://blockstream.info/api/tx/d95cc33c3ec74c75ad8a1868c88007f11a04e411ff073996ffbea52d12937319/hex

put the above link in browser we can following:
![截屏2024-04-17 18 15 50](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/20e84bc0-33e7-4d2c-bc7a-892bc8e350f6)


Let's see how we can use code to do that, add a file named tx_fetcher.go and add following code:

```g
package transaction

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

type TransactionFetcher struct{}

func NewTransactionFetcher() *TransactionFetcher {
	return &TransactionFetcher{}
}

func (t *TransactionFetcher) getURL(testnet bool) string {
	if testnet {
		return "https://blockstream.info/testnet/api/tx"
	}

	return "https://blockstream.info/api/tx"
}

func (t *TransactionFetcher) Fetch(txID string, testnet bool) []byte {
	url := fmt.Sprintf("%s/%s/hex", t.getURL(testnet), txID)
	fmt.Printf("fetching url: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		panic(fmt.Sprintf("fetch transaction err: %v\n", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Sprintf("read body err:%v\n", err))
	}

	fmt.Printf("response: %s\n", string(body))

	buf, err := hex.DecodeString(string(body))
	if err != nil {
		panic(err)
	}

	return buf
}

```
TransactionFetcher will used to get the raw data for given transaction with given hash id. It receives the hash id of given transaction and construct a request,
if ask transaction from testnet, the url would like following:

https://blockstream.info/testnet/api/tx/id/hex

for mainnet, the request url like following:
https://blockstream.info/api/tx/id/hex

let's have a try for it, remember we parsed the legacy transaction and get the prevous transaction hash id is: 
d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81

then the code using TransactonFetch to get the raw data of the transaction is :
```g
func main() {
	fetcher := tx.NewTransactionFetcher()
	fetcher.Fetch("d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81", false)
}
```
The above code will give the following result:
```g
response: 0100000002137c53f0fb48f83666fcfd2fe9f12d13e94ee109c5aeabbfa32bb9e02538f4cb000000006a47304402207e6009ad86367fc4b166bc80bf10cf1e78832a01e9bb491c6d126ee8aa436cb502200e29e6dd7708ed419cd5ba798981c960f0cc811b24e894bff072fea8074a7c4c012103bc9e7397f739c70f424aa7dcce9d2e521eb228b0ccba619cd6a0b9691da796a1ffffffff517472e77bc29ae59a914f55211f05024556812a2dd7d8df293265acd8330159010000006b483045022100f4bfdb0b3185c778cf28acbaf115376352f091ad9e27225e6f3f350b847579c702200d69177773cd2bb993a816a5ae08e77a6270cf46b33f8f79d45b0cd1244d9c4c0121031c0b0b95b522805ea9d0225b1946ecaeb1727c0b36c7e34165769fd8ed860bf5ffffffff027a958802000000001976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88aca5515e00000000001976a914e82bd75c9c662c3f5700b33fec8a676b6e9391d588ac00000000
```
Then we can parse the given transaction raw data, get its amount of  transaction output by using the previousTransactionIndex of current transaction, let's add 
a new function for input:
```g
func (t *TransactinInput) Value(testnet bool) *big.Int {
	previousTxID := fmt.Sprintf("%x", t.previousTransaction)
	previousTX := t.fetcher.Fetch(previous, testnet)
	tx := ParseTransaction(previousTX)
	//get the amount of previous transaction output
	return tx.txOutputs[t.previousTransactionIdex.Int64()].amount
}
```
We can't run above code yet because we don't know how to parse the ScriptSig object, we will dive deep into it in next several sections
