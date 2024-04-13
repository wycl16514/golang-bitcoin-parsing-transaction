For bitcoin, the most important component is transaction. If you need to make deal with others by bitcoin, you may need to send some coins to others just like you need to pay money for 
buying goods or services from others. There are four fields that are key to the component of transaction, they are version, inputs, outputs and lock time. The binary content of 
transaction may differ for different version, we need to parse the number of version to decide how to decode the binary data from transaction. Inputs are records about bitcoins that are
sent to you from other people, and outputs are info about the money you spend and who will receive you morney. And lock time is the frozen time for the transaction, after the passing of
lock time, the transaction will accept by the bitcoin networks.

Let's see an bianry example for transaction:

0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600

Let's seperate it into the four parts above:

version: 01000000 (the first 4 bytes, it is in little endian, we need to convert it to big endian)

inputs: 
01813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff

outputs:
02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac

locktime:
19430600

The difficult part for parsing transaction is the parsing input and output, the length of input and output can be changed for different transaction. In code implementation for bitcoin
transaction, we need to pay special efforts to the parsing of inputs and outputs, first we need to build up the framework, create a new folder called transaction and add three files:
transaction.go, input.go, output.go, and we add struct define for input and output first, in input.go, we have following code:
```g
package transaction

import (
	"bufio"
)

type TransactinInput struct {
	reader *bufio.Reader
}

func NewTransactinInput(reader *bufio.Reader) *TransactinInput {

	return &TransactinInput{
		reader: reader,
	}
}


```
The following is code for output:
```g
package transaction

import (
	"bufio"
)

type TransactionOutput struct {
	reader *bufio.Reader
}

func NewTransactionOutput(reader *bufio.Reader) *TransactionOutput {
	return &TransactionOutput{
		reader: reader,
	}
}

```

we need to move the little endian conversion from elliptical curve to transaction package, add u util.go file and move the code for endian conversion to here:
```g
package transaction

import (
	"math/big"

	"github.com/tsuna/endian"
)

type LITTLE_ENDIAN_LENGTH int

const (
	LITTLE_ENDIAN_2_BYTES LITTLE_ENDIAN_LENGTH = iota
	LITTLE_ENDIAN_4_BYTES
	LITTLE_ENDIAN_8_BYTES
)

func BigIntToLittleEndian(v *big.Int, length LITTLE_ENDIAN_LENGTH) []byte {
	switch length {
	case LITTLE_ENDIAN_2_BYTES:
		val := v.Int64()
		littleEdianVal := endian.HostToNetUint16(uint16(val))
		p := big.NewInt(int64(littleEdianVal))
		return p.Bytes()
	case LITTLE_ENDIAN_4_BYTES:
		val := v.Int64()
		littleEdianVal := endian.HostToNetUint32(uint32(val))
		p := big.NewInt(int64(littleEdianVal))
		return p.Bytes()
	case LITTLE_ENDIAN_8_BYTES:
		val := v.Int64()
		littleEdianVal := endian.HostToNetUint64(uint64(val))
		p := big.NewInt(int64(littleEdianVal))
		return p.Bytes()
	}

	return nil
}

func LittleEndianToBigInt(bytes []byte, length LITTLE_ENDIAN_LENGTH) *big.Int {
	switch length {
	case LITTLE_ENDIAN_2_BYTES:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint16(uint16(p.Uint64()))
		return big.NewInt(int64(val))
	case LITTLE_ENDIAN_4_BYTES:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint32(uint32(p.Uint64()))
		return big.NewInt(int64(val))
	case LITTLE_ENDIAN_8_BYTES:
		p := new(big.Int)
		p.SetBytes(bytes)
		val := endian.NetToHostUint64(uint64(p.Uint64()))
		return big.NewInt(int64(val))
	}

	return nil
}

```
Let's create a file named transaction.go and add the skeleton for it:
```g
package transaction

import (
	"bytes"
	"fmt"
  "bufio"
	"math/big"
)

type Transaction struct {
	version   *big.Int
	txInputs  []*TransactinInput
	txOutputs []*TransactionOutput
	lockTime  *big.Int
	testnet   bool
}

func ParseTransaction(binary []byte) *Transaction {
	reader := bytes.NewReader(binary)
	bufReader := bufio.NewReader(reader)

	verBuf := make([]byte, 4)
	bufReader.Read(verBuf)

	version := LittleEndianToBigInt(verBuf, LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("transaction version: %x\n", version)

	return nil
}

```
Let's test the reading of version first:
```g
package main

import (
	"encoding/hex"
	tx "transaction"
)

func main() {
	binaryStr := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
	binary, err := hex.DecodeString(binaryStr)
	if err != nil {
		panic(err)
	}

	tx.ParseTransaction(binary)
}

```
Running the above code will output following:
```g
transaction version: 1
```
Let's see how to parse the input data, a input for the current transaction is the output for previous transaction. For example you use 10 dollar to buy a cup of coffee with price 3, then you spend 3
dollar and receive 7 dollar at return, when you use thhis 7 dolloars to buy other things, then it will be the input for new transaction. Therefore a input will do two things:

1, reference to previous transaction

2, proof the the money for the input is indeed belongs to you

The second point will need the help of Signing from elliptical curve we developed. As you can see the txInputs filed in Transaction, its an slice that means there may more than one instance of input
in a transaction. For the example above, if there has bill note with value of 7, then it will create on input, but we don't have bill note with face value of 7, in order to pay you 7 dollars, the 
shop may give you a note with 5 dolloar and one with 2 dollar, then you will have two inputs for next transaction, if the shop give you one 5 dollor bill and two one dollar bill, then you will have 
three inputs for next transaction, if the shop repay you 70o pennies, then you have 700 inputs for the next transaction.

Let's see how to decode the binary data of input:
inputs: 
01813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff

The first job we need to do is get the number of inputs, we need some effort to parse the number for inputs, the bytes used for encoding the number of input can be vary, we will use the following 
steps for parsing the number of inputs:
1, if the value of first byte for input is smaller than 253(0xfd), then the value of this byte is the number for input

2, if the value of first byte for input is 0xfd, then we need to read the following 2 bytes for the number of input, these two bytes are encoded as little endian.

3, if the value of first byte for input is 0xfe, then we need to read the following 4 bytes for the number of input, these four bytes are encode as little endian.

4, if the value of first byte for input is 0xff, then we need to read the following 8 bytes for the number of input

Let's put this scheme into code in util.go:
```g

func ReadVarint(reader *bytes.Reader) *big.Int {
	//if the value of first byte is < 0xfd, then the value is one byte
	i := make([]byte, 1)
	reader.Read(i)
	v := new(big.Int)
	v.SetBytes(i)
	if v.Cmp(big.NewInt(int64(0xfd))) < 0 {
		return v
	}

	if v.Cmp(big.NewInt(int64(0xfd))) == 0 {
		//first byte is 0xfd, read the following 2 bytes for number of input
		i1 := make([]byte, 2)
		reader.Read(i1)
		return LittleEndianToBigInt(i1, LITTLE_ENDIAN_2_BYTES)
	}

	if v.Cmp(big.NewInt(int64(0xfe))) == 0 {
		//first byte 0xfe, read the following 4 bytes to get the number of input
		i1 := make([]byte, 4)
		return LittleEndianToBigInt(i1, LITTLE_ENDIAN_4_BYTES)
	}

	//first byte is 0xff, read following 8 bytes for the number of input
	i1 := make([]byte, 8)
	reader.Read(i1)
	return LittleEndianToBigInt(i, LITTLE_ENDIAN_8_BYTES)
}

func EncodeVarint(v *big.Int) []byte {
	//if the value < 0xfd, one byte is enough
	if v.Cmp(big.NewInt(int64(0xfd))) < 0 {
		vBytes := v.Bytes()
		return []byte{vBytes[0]}
	} else if v.Cmp(big.NewInt(int64(0x10000))) < 0 {
		//if value >= 0xfd and < 0x10000, then need 2 bytes
		buf := []byte{0xfd}
		vBuf := BigIntToLittleEndian(v, LITTLE_ENDIAN_2_BYTES)
		buf = append(buf, vBuf...)
		return buf
	} else if v.Cmp(big.NewInt(int64(0x100000000))) < 0 {
		//value >= 0xFFFF and <= 0xFFFFFFFF, then need 4 bytes
		buf := []byte{0xfe}
		vBuf := BigIntToLittleEndian(v, LITTLE_ENDIAN_4_BYTES)
		buf = append(buf, vBuf...)
		return buf
	}

	p := new(big.Int)
	p.SetString("10000000000000000", 16)
	if v.Cmp(p) < 0 {
		//need 8 bytes
		buf := []byte{0xff}
		vBuf := BigIntToLittleEndian(v, LITTLE_ENDIAN_8_BYTES)
		buf = append(buf, vBuf...)
		return buf
	}

	panic(fmt.Sprintf("integer too large: %x\n", v))
}

```

Let's try to parse the version of input , first we need to know there are two kinds of transaction, one is called legacy transaction , the binary we have seen above is for this catagory, the second
kind is called segwit transaction, the difference of the two is that, if the byte follow by the first 8 bytes which is the version is 0, then it is the second kind transaction, we need to skip the 
following 2 bytes(the values of these two bytes are 0x00, 0x01) in order to get the input count, for example you can use the following command to get raw data for the segwit transaction:
```g
 curl https://docs-demo.btc.quiknode.pro/ \
    -X POST \
    -H "Content-Type: application/json" \
    --data '{"method": "getrawtransaction", "params": ["464dd72e17069e6b55487ae75971df57280091988fbae9ca8fc1abedafddfcbc", 0]}'

```
The result for the command is :
```g
{"result":"01000000000102197393122da5beff963907ff11e4041af10780c868188aad754cc73e3cc35cd9010000001716001462c61a14835b032d5acbe190291d80d0cc5ca28e00000000feae2204104ffe542f30a20012a5b8e2b54a6f61f592520b511801b2237b5ed80100000017160014b30be91e50402cda780c56a3e1c350b1086c80af000000000200a3e111000000001976a914e60c9ac5f72d1d620287a0fc35656bceae5e2ab988ac525d35130000000017a9144795995aff558cc538669ebfecffbe5c9837d5ca870247304402207dd1e7c6c596041276b5285dd3747f586ad819a24acdf0ad60b1faa82af00d3b022046a22dd57df4b72ac165e05b4a6cf8dbecfcfad8f16ae7353df56638ebbf5d1f012103a1a226c5047672af98b2e673751dc69f0140b957753d9c1a789c243100292c6f024730440220670625143c3dfc7a862659a79cbf4ad0f84ff1509bd052cfbfbcdba7adf501f9022015f14a6ee1ae7a8f9fec1070d8a97195422b76a317286c816392cb150d7eb76d012102c910a40bf5726168acc5a8318b0505375e877d4d74448f32ef48156794e657f900000000","error":null,"id":null}
```
you can see after the first 8 bytes, it followed by byte 0x00, which means it is a segwit transaction, skip the following two bytes (00,01), then we get 02, this is the number for the input of this
transaction, we can goto blocakchain.com and search the transaction using its hash id:464dd72e17069e6b55487ae75971df57280091988fbae9ca8fc1abedafddfcbc, and you get following:

![截屏2024-04-13 14 15 16](https://github.com/wycl16514/golang-bitcoin-parsing-transaction/assets/7506958/fdd6e748-505d-4d4d-aabe-0f539ac42d2d)

Let's write the code to parse the count of input, in transaction.go, add the following code:
```g
func getInputCount(bufReader *bufio.Reader) *big.Int {
	/*
		if the first byte is 0, then its segwit transaction, we need to skip the
		first two byes
	*/
	firstByte, err := bufReader.Peek(1)
	if err != nil {
		panic(err)
	}

	if firstByte[0] == 0x00 {
		//skip the first two bytes
		skipBuf := make([]byte, 2)
		_, err = bufReader.Read(skipBuf)
		if err != nil {
			panic(err)
		}
	}

	count := ReadVarint(bufReader)
	fmt.Printf("input count is : %x\n", count)
	return count
}

func ParseTransaction(binary []byte) *Transaction {
  ...
  getInputCount(bufReader)

	return nil
}
```
Running the above code we get:
transaction version: 1
input count is : 1

If we change the raw transaction data to the segwit transaction we got above:
```
func main() {
	//legacy transaction
	//binaryStr := "0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600"
	//segwit transaction
	binaryStr := "01000000000102197393122da5beff963907ff11e4041af10780c868188aad754cc73e3cc35cd9010000001716001462c61a14835b032d5acbe190291d80d0cc5ca28e00000000feae2204104ffe542f30a20012a5b8e2b54a6f61f592520b511801b2237b5ed80100000017160014b30be91e50402cda780c56a3e1c350b1086c80af000000000200a3e111000000001976a914e60c9ac5f72d1d620287a0fc35656bceae5e2ab988ac525d35130000000017a9144795995aff558cc538669ebfecffbe5c9837d5ca870247304402207dd1e7c6c596041276b5285dd3747f586ad819a24acdf0ad60b1faa82af00d3b022046a22dd57df4b72ac165e05b4a6cf8dbecfcfad8f16ae7353df56638ebbf5d1f012103a1a226c5047672af98b2e673751dc69f0140b957753d9c1a789c243100292c6f024730440220670625143c3dfc7a862659a79cbf4ad0f84ff1509bd052cfbfbcdba7adf501f9022015f14a6ee1ae7a8f9fec1070d8a97195422b76a317286c816392cb150d7eb76d012102c910a40bf5726168acc5a8318b0505375e877d4d74448f32ef48156794e657f900000000"
	binary, err := hex.DecodeString(binaryStr)
	if err != nil {
		panic(err)
	}

	tx.ParseTransaction(binary)
}
```
Then we will get the following result:
```g
transaction version: 1
input count is : 2
```

