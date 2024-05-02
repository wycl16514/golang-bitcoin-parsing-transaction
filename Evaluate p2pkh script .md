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
Run the above code and check the serialization result is the same as the binary content of script byte slice, if so which means we construct the ScriptSig object successfully.

Now we need to add code to handle command like OP_EQUALVERIFY, OP_DUP, OP_HASH160, add the following code in op.go, in the code, we do some
refactory, we move stack, altStack ad cmds from ScriptSig to BitcoinOperationCode:
```g
type BitcoinOpCode struct {
	opCodeNames map[int]string
	stack       [][]byte
	altStack    [][]byte
	cmds        [][]byte
}

func NewBitcoinOpCode() *BitcoinOpCode {
    ....
    return &BitcoinOpCode{
		opCodeNames: opCodeNames,
		stack:       make([][]byte, 0),
		altStack:    make([][]byte, 0),
		cmds:        make([][]byte, 0),
	}
}

func (b *BitcoinOpCode) opDup() bool {
	if len(b.stack) < 1 {
		return false
	}

	b.stack = append(b.stack, b.stack[len(b.stack)-1])
	return true
}

func (b *BitcoinOpCode) opHash160() bool {
	if len(b.stack) < 1 {
		return false
	}

	element := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	hash160 := ecc.Hash160(element)
	b.stack = append(b.stack, hash160)
	return true
}

func (b *BitcoinOpCode) opEqual() bool {
	if len(b.stack) < 2 {
		return false
	}

	elem1 := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	elem2 := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	if bytes.Equal(elem1, elem2) {
		b.stack = append(b.stack, b.EncodeNum(1))
	} else {
		b.stack = append(b.stack, b.EncodeNum(0))
	}
	return true
}

func (b *BitcoinOpCode) opVerify() bool {
	if len(b.stack) < 1 {
		return false
	}

	elem := b.stack[len(b.stack)-1]
	b.stack = b.stack[0 : len(b.stack)-1]
	if b.DecodeNum(elem) == 0 {
		return false
	}

	return true
}

func (b *BitcoinOpCode) opEqualVerify() bool {
	resEqual := b.opEqual()
	resVerify := b.opVerify()
	return resEqual && resVerify
}

func (b *BitcoinOpCode) RemoveCmd() []byte {
	cmd := b.cmds[0]
	b.cmds = b.cmds[1:]
	return cmd
}

func (b *BitcoinOpCode) HasCmd() bool {
	return len(b.cmds) > 0
}

func (b *BitcoinOpCode) AppendDataElement(element []byte) {
	b.stack = append(b.stack, element)
}
```

opDup is for handling of OP_DUP command, it just copy the top element on the stack and push it on to the stack, this will result in two 
elements with same value on the stack. opHash160 is for command OP_HASH160, it takes the top elemenet from stack and compute its hash160
and push the result onto the stack.

opEqual is for OP_EQUAL command, it takes the top two elements from the top of stack, compare their value, if they are the same, push value
1 on the stack, if not push value 0 on the stack, opVerify is for command OP_VERIFY, it checks the element on the top of stack, if its value
is 0, then return false, otherwise return true.

opEqualVerify is for command OP_EQUALVERIFY, it is the combination of OP_EQUAL and OP_VERIFY, we add three new helper functions to 
BitcoinOperationCode, RemoveCmd takes one command from the command array, HasCmd check whether there are commands need to be handle,
and AppendDataElement append data element to the command stack.

Then we goto script_sig.go and modify its code like following:
```g
func InitScriptSig(cmds [][]byte) *ScriptSig {
	bitcoinOpCode := NewBitcoinOpCode()
	bitcoinOpCode.cmds = cmds
	return &ScriptSig{
		bitcoinOpCode: bitcoinOpCode,
	}
}

func (s *ScriptSig) Evaluate(z []byte) bool {
	for s.bitcoinOpCode.HasCmd() {
		cmd := s.bitcoinOpCode.RemoveCmd()
		if len(cmd) == 1 {
			//this is an op code, run it
			opRes := s.bitcoinOpCode.ExecuteOperation(int(cmd[0]), z)
			if opRes != true {
				return false
			}
		} else {
			s.bitcoinOpCode.AppendDataElement(cmd)
		}
	}

	/*
		After running all operation in the script and the stack is empty,
		then the evaluation fail, otherwise we check the top element of stack
		if it is not 0 then success, otherwise fail
	*/
	if len(s.bitcoinOpCode.stack) == 0 {
		return false
	}

	if len(s.bitcoinOpCode.stack[0]) == 0 {
		return false
	}

	return true
}
```

we simplify the inputs for ExecuteOperation of BitcoinOperationCode because we move those two stacks and command array into it,
by the refactory, the code can be clearer than before. Now go to main.go to evaluate the script:
```
func main() {
    ....
    evalRes := scriptSig.Evaluate(z.Bytes())
    fmt.Printf("result of script evaluation is %v\n", evalRes)
}
```

Run the code and we can get the following result:
```g
result of script evaluation is true
```

We have handcraft the script by using our hand, actually a script is constructed by using the ScriptSig of current transaction and ScriptPubKey from the output of 
previous trnascation like following:

![bitcoin_script (5)](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/576de25a-c1f1-4189-9aa7-7184bb978be7)

Let's see how we can construct the script from current transaction and previous transaction, fist we need add code to parse the transaction output:
```g
package transaction

import (
	"bufio"
	"math/big"
)

type TransactionOutput struct {
	amount       *big.Int
	scriptPubKey *ScriptSig
}

func NewTransactionOutput(reader *bufio.Reader) *TransactionOutput {
        /*
          amount is in stashi which is 1/100,000,000 of one bitcoin
        */
	amountBuf := make([]byte, 8)
	reader.Read(amountBuf)
	amount := LittleEndianToBigInt(amountBuf, LITTLE_ENDIAN_8_BYTES)
	script := NewScriptSig(reader)
	return &TransactionOutput{
		amount:       amount,
		scriptPubKey: script,
	}
}

func (t *TransactionOutput) Serialize() []byte {
	result := make([]byte, 0)
	result = append(result,
		BigIntToLittleEndian(t.amount, LITTLE_ENDIAN_8_BYTES)...)
	result = append(result, t.scriptPubKey.Serialize()...)
	return result
}

```

and let's add the serialize method for transaction input too:
```g
func (t *TransactinInput) Serialize() []byte {
	result := make([]byte, 0)
	result = append(result, reverseByteSlice(t.previousTransaction)...)
	result = append(result,
		BigIntToLittleEndian(t.previousTransactionIdex,
			LITTLE_ENDIAN_4_BYTES)...)
	result = append(result, t.scriptSig.Serialize()...)
	result = append(result,
		BigIntToLittleEndian(t.sequence, LITTLE_ENDIAN_4_BYTES)...)
	return result
}

```

We add a method to combine the commands of two ScriptSig object:
```g
func (s *ScriptSig) Add(script *ScriptSig) *ScriptSig {
	cmds := make([][]byte, 0)
	cmds = append(cmds, s.bitcoinOpCode.cmds...)
	cmds = append(cmds, script.bitcoinOpCode.cmds...)
	return InitScriptSig(cmds)
}
```
Now we get the previous transaction, extract its transaction output script oject, using the method above to construct the script we need to evaluate:
```g
func (t *Transaction) GetScript(idx int, testnet bool) *ScriptSig {
	if idx < 0 || idx >= len(t.txInputs) {
		panic("invalid index for transaction input")
	}

	txInput := t.txInputs[idx]
	return txInput.Script(testnet)
}
```
Then we can goto the main.go and test the code above:
```g

func main() {
	//legacy transaction
	binaryStr := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"

	binary, err := hex.DecodeString(binaryStr)
	if err != nil {
		panic(err)
	}

	transaction := tx.ParseTransaction(binary)

	script := transaction.GetScript(0, false)
	//this is not our transaction and we don't have its message
	script.Evaluate([]byte{})
}
```
We don't have the message for the given sgnature because the transaction is not created by us. By debugging the above code we can find out the script for previous
transaction output has command like OP_DUP, OP_HASH160, OP_EQUALVERIFY and OP_CHECKSIG, the evaluate will fail at the step of OP_CHECKSIG because we don't have its
message, but we can make sure the combined script is correct.


