package transaction

import (
	"bufio"
	"bytes"
	"fmt"
	"math/big"
)

type Transaction struct {
	version   *big.Int
	txIputs   []*TransactionInput
	txOutputs []*TransactionOutput
	lockTime  *big.Int
	testnet   bool
}

func getInputCount(bufReader *bufio.Reader) *big.Int {
	/*
		if the first byte of input is 0, then witness transaction,
		we need to skip the first two bytes(0x00, 0x01)
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
	fmt.Printf("input count is: %x\n", count)
	return count
}

func ParseTransaction(binary []byte) *Transaction {
	transaction := &Transaction{}
	reader := bytes.NewReader(binary)
	bufReader := bufio.NewReader(reader)

	verBuf := make([]byte, 4)
	bufReader.Read(verBuf)

	version := LittleEndianToBigInt(verBuf, LITTLE_ENDIAN_4_BYTES)
	fmt.Printf("transaction version:%x\n", version)
	transaction.version = version

	inputs := getInputCount(bufReader)
	transactionInputs := []*TransactionInput{}
	for i := 0; i < int(inputs.Int64()); i++ {
		input := NewTractionInput(bufReader)
		transactionInputs = append(transactionInputs, input)
	}
	transaction.txIputs = transactionInputs

	//read output counts
	outputs := ReadVarint(bufReader)
	transactionOutputs := []*TransactionOutput{}
	for i := 0; i < int(outputs.Int64()); i++ {
		output := NewTractionOutput(bufReader)
		transactionOutputs = append(transactionOutputs, output)
	}
	transaction.txOutputs = transactionOutputs

	//get last four bytes for lock time
	lockTimeBytes := make([]byte, 4)
	bufReader.Read(lockTimeBytes)
	transaction.lockTime = LittleEndianToBigInt(lockTimeBytes, LITTLE_ENDIAN_4_BYTES)

	return transaction

}

func (t *Transaction) GetScript(idx int, testnet bool) *ScriptSig {
	if idx < 0 || idx >= len(t.txIputs) {
		panic("invalid idx for transaction input")
	}

	txInput := t.txIputs[idx]
	return txInput.Script(testnet)
}
