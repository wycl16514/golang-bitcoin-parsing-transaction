Inputs and Outputs are main components in transaction, and we need some effort to get info out from their raw data. 


There are four components in the raw data of input:
1, Previous trnasaction ID, this is 32 bytes data in little endian , it is hash256 hash for previous transaction, for example if your company pays you 10 dollars,
this paying is the previous transacton. you use this 10 dollars to buy a 3 dollar coffee, then the buying is current transaction, 
the previous transaction id is the hash256 result of transaction that you company paying you 10 dollars.

2, Previous transaction index.  Your company may pay you twice, one for you salary (10 dollars), one for you bounus (20 dollars), then you will have two inputs for the current transaction. The previous
ID let you choose which payment you will use to pay the current coffee.

3. ScriptSig, the use of this filed just like the signing on you check or input password when you swap your credit card.It is used to verify you are the real owner of the input money in this transaction.
   You don't want others steal you money for their own spend. This field is a chunk of data with variable length, it has to do with bitcoin's smart contract language, we will focus on this in later sections

4. Sequence, this field is depricated, we just ingore it for saving brain cells.

Let's see the raw data for input and extract those fielsd out, following is the raw data for transaction and data in { and } are for input:

0100000001{ 813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfc
f21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e
213bf016b278afeffffff }02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75
f40df79fea1288ac19430600 

The first 32 bytes from { are previous transaction id in little endian which is following:

813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1

The following 4 bytes is prevous transaction ID in little endian:
00000000

Then follows the raw data for ScriptSig, it is variable length, different transaction will have different length of ScriptSig:
6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f
67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a

Finally the last four bytes are for sequence:
feffffff

Let's see how to use code for the parsing of input, first we create a ScriptSig Object and we will pay attention to the parsing of it in later sections, create a file name script_sig.go in transaction
folder and input the following code:
```g
package transaction

import (
	"bufio"
)

type ScriptSig struct {
	reader *bufio.Reader
}

func NewScriptSig(reader *bufio.Reader) *ScriptSig {

	return &ScriptSig{
		reader: reader,
	}
}

```
Then we goto the transaction.go, create TrasactionInput objects according to its count:
```g
func ParseTransaction(binary []byte) *Transaction {
	transcion := Transaction{}
	reader := bytes.NewReader(binary)
	bufReader := bufio.NewReader(reader)

	verBuf := make([]byte, 4)
	bufReader.Read(verBuf)

	version := LittleEndianToBigInt(verBuf, LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("transaction version: %x\n", version)

	transcion.version = version

	inputs := getInputCount(bufReader)
	transactionInputs := []*TransactinInput{}
	//create transaction input object by its count
	for i := 0; i < int(inputs.Int64()); i++ {
		input := NewTransactinInput(bufReader)
		transactionInputs = append(transactionInputs, input)
	}
	transcion.txInputs = transactionInputs

	return &transcion
}
```

Go to input.go, we will extract those four filds from raw data and generate the transaction input object:
```g
package transaction

import (
	"bufio"
	"fmt"
	"math/big"
)

type TransactinInput struct {
	previousTransaction     []byte
	previousTransactionIdex *big.Int
	scriptSig               *ScriptSig
	sequence                *big.Int
}

func reverseByteSlice(bytes []byte) []byte {
	reverseBytes := []byte{}
	for i := len(bytes) - 1; i >= 0; i-- {
		reverseBytes = append(reverseBytes, bytes[i])
	}

	return reverseBytes
}

func NewTransactinInput(reader *bufio.Reader) *TransactinInput {
	transactionInput := TransactinInput{}
	//Convert the first 32 bytes from little endian to big endian by reversing it
	previousTransaction := make([]byte, 32)
	reader.Read(previousTransaction)
	transactionInput.previousTransaction = reverseByteSlice(previousTransaction)
	fmt.Printf("previous transaction id: %x\n", transactionInput.previousTransaction)

	//following four bytes are previous transaction index
	idx := make([]byte, 4)
	reader.Read(idx)
	transactionInput.previousTransactionIdex = LittleEndianToBigInt(idx, LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("previous transaction idex: %x\n", transactionInput.previousTransactionIdex)

	transactionInput.scriptSig = NewScriptSig(reader)

	seqBytes := make([]byte, 4)
	reader.Read(seqBytes)
	transactionInput.sequence = LittleEndianToBigInt(seqBytes, LITTLE_ENDIAN_4_BYTES)

	return &transactionInput
}

```
Let's try the parsing on the raw data of legacy transaction:
```g
func main() {
	//legacy transaction
	binaryStr := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
	//segwit transaction
	//binaryStr := "01000000000102197393122da5beff963907ff11e4041af10780c868188aad754cc73e3cc35cd9010000001716001462c61a14835b032d5acbe190291d80d0cc5ca28e00000000feae2204104ffe542f30a20012a5b8e2b54a6f61f592520b511801b2237b5ed80100000017160014b30be91e50402cda780c56a3e1c350b1086c80af000000000200a3e111000000001976a914e60c9ac5f72d1d620287a0fc35656bceae5e2ab988ac525d35130000000017a9144795995aff558cc538669ebfecffbe5c9837d5ca870247304402207dd1e7c6c596041276b5285dd3747f586ad819a24acdf0ad60b1faa82af00d3b022046a22dd57df4b72ac165e05b4a6cf8dbecfcfad8f16ae7353df56638ebbf5d1f012103a1a226c5047672af98b2e673751dc69f0140b957753d9c1a789c243100292c6f024730440220670625143c3dfc7a862659a79cbf4ad0f84ff1509bd052cfbfbcdba7adf501f9022015f14a6ee1ae7a8f9fec1070d8a97195422b76a317286c816392cb150d7eb76d012102c910a40bf5726168acc5a8318b0505375e877d4d74448f32ef48156794e657f900000000"
	binary, err := hex.DecodeString(binaryStr)
	if err != nil {
		panic(err)
	}

	tx.ParseTransaction(binary)
```
Then we get the following result:
```g
transaction version: 1
input count is : 1
previous transaction id: d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81
previous transaction idex: 0
```

   
