The p2pk (pay to pubkey) script we executed in last section has some problems. First the uncompressed public is too long, it is 65 bytes in binary data, when encode it into string,the length can be doubled, even the compressed SEC format is still too long and difficult to be handled by wallet. 

Second it is no secure enough when machine exachange info usinguncompressed SEC format. And Bitcoin nodes need to indexed all outputs for certain purpose, 
the uncompressed SEC format will requires more disk resources on the part of nodes.

In order to handle problems in p2pk script, there is a updated script format called p2pkh, in this script ,it will use sha256 and ripemd160 hashing to make the 
address much shorter compared with the uncompressed public key. The structure of p2pkh script is like following:


![bitcoin_script](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/e330c0ee-c047-4fa7-b736-10c1e71611a8)

the first two element of the script commands are data elements therefore they are push to stack directyly:

![bitcoin_script](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/0a69fe30-96be-4be6-b8b3-df1b59e78b8d)

Notice we push the signature first and then the pubkey, that's why the pubkey is at the top of stack. Then we execute the OP_DUP, this command duplicate the top 
elment on the stack and push it to the stack that cause the stack has two pubkey element at the top:

![bitcoin_script](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/73a16f6c-0eb7-4c70-8076-00813fee939e)


Now take the OP_HASH160 command from the script and take the top element from the stack, do hash160 operation on the elment and push the result to the top of the 
stack:

![bitcoin_script (1)](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/1c926c3c-6ccc-44de-a634-7482627a5247)

Now there is a hash element on the script, then we push it directly to the stack, this results in two hash elements on the stack:

![bitcoin_script (2)](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/d4dee602-911a-4927-8a46-2112ab5651a5)

Now there is OP_EQUALVERIFY command on the top of script, it will take the top two elements from the stack and check they are equal or not. If they are equal, 
the script continue is execution, otherwise the script stops and fail out. we assume the top two elements are equal then they are taken from the top of the stack:

![bitcoin_script (3)](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/dc531c92-5946-4337-a248-2d3a27dde75f)

Finnaly there is a OP_CHECKSIG on the script and pubkey and signature (in DER binary format) on the stack, this is what we have done in prevoius section, execute the
command, if the signature can be verified by the pubkey, there will be only one element with value 1 on the stack:

![bitcoin_script (4)](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/46226c0b-6ff2-4e90-8c21-ed91f990e346)

There are two ways the script can fail, the first is when we executing the OP_EQUALVERIFY command, when the top two elements on the stack are not equal. The second is
when we are executing the OP_CHECKSIG command, if the pubkey can't verify the signature. It seems it is much troublesome than p2pk script right? We will see why it is
better than the p2pk script in later section.

Let's create a script that can test the process above, in main.go, we have the following code:
```
package main

import (
	"bufio"
	"bytes"
	ecc "elliptic_curve"
	"fmt"
	"math/big"
	tx "transaction"
)

func main() {
	e := new(big.Int)
	e.SetBytes([]byte(string("my secrect")))
	z := new(big.Int)
	z.SetBytes([]byte(string("my message")))
	priveKey := ecc.NewPrivateKey(e)
	signature := priveKey.Sign(z)
	sigDER := signature.Der()
	//append byte as hash type at the end
	sigDER = append(sigDER, 0x01)
	fmt.Printf("len of sigDER is :%d\n", len(sigDER))
	fmt.Printf("content of sigDER is: %x\n", sigDER)

	pubKey := priveKey.GetPublicKey()
	_, pubKeySec := pubKey.Sec(true)
	fmt.Printf("len of pub key sec: %d\n", len(pubKeySec))
	fmt.Printf("pub key sec compressed: %x\n", pubKeySec)
	pubKeySecHash160 := ecc.Hash160(pubKeySec)
	fmt.Printf("pub key sec compressed with hash160:%x\n", pubKeySecHash160)

	script := make([]byte, 0)
	script = append(script, byte(len(sigDER)))
	script = append(script, sigDER...)
	script = append(script, byte(len(pubKeySec)))
	script = append(script, pubKeySec...)
	script = append(script, tx.OP_DUP)
	script = append(script, tx.OP_HASH160)
	script = append(script, byte(len(pubKeySecHash160)))
	script = append(script, pubKeySecHash160...)
	script = append(script, tx.OP_EQUALVERIFY)
	script = append(script, tx.OP_CHECKSIG)
	scriptLen := len(script)
	totalLen := tx.EncodeVarint(big.NewInt(int64(scriptLen)))
	script = append(totalLen, script...)
	reader := bytes.NewReader(script)
	bufReader := bufio.NewReader(reader)
	scriptSig := tx.NewScriptSig(bufReader)
	fmt.Printf("serialize of the script object: %x\n", scriptSig.Serialize())
}

```
Run the above code and check the serialization result is the same as the binary content of script byte slice, if so which means we construct the ScriptSig object
successfully.
