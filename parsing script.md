Bitcoin has smart contract which is the composition of script in the transaction we methioned in last sections. The script is a kind of executable instructions like assembly language.
But is has its own specilities. It is stack-based, and it can't have loop, in each instruction of bitcoin script, it can only has one operand, and all data that in the running of the
script should be put in the stack.

Bitcoin script is used to run some process to unlock the fund in transaction or verify you are the right recipient of the fund. Just like when you deposite your money in to your bank
account, There are two kinds of instruction in bitcoin script, one is called elements and the other is operations. Elements are used to push data onto the stack, and operations are used
to pop data from stack, do some computation and push the result on the stack. For example instrution OP_DUP is used to copy the data on the top of stack and push the same result on the 
stack. If the there is a value X on the top of stack, and after the execution of OP_DUP then there are two elements with value X on the top of stack.

Beside instruct OP_DUP, the instruct OP_HASH160 take the element on the top of stack, then do a sha256 computation followed by a ripemd160 and push the result on the top of stack. The other
important operation is OP_CHECKSIG, it takes two element from the top of stack, the first as public key, and the second as signature and check whether the public key can be used to verify the
signatrue. If it can then it pushes a value of 1 on the stack, otherwise it pushes 0 on to the stack. By the way each instruction conrresponding to a number, for example OP_DUP with value 118.

Let's see how we can use code to parse the script field, the parsing process is the same for ScriptSig and ScriptPubKey. First we check the first byte in the beginning of the field, if its value 
between 0x01 and 0x4b(75), by which we use n to indicate the value of this byte, then we need to read the following n bytes after the beginning byte and push them on to the stack, 
if the the value of n is not in the range of [0x01, 0x4b], which means this byte represents an operation, then we need to check which operation it takes, following are values for some operations:

0x00: OP_0
0x51: OP_1
0x60: OP_16
0x76: OP_DUP
0x93: OP_ADD
0xa9: OP_HASH160
0xac: OP_CHECKSIG

you can look for more details about instruction here: https://en.bitcoin.it/wiki/Script

as we mentionded above, if the byte with value in [0x01, 0x4b] means we need to push the following bytes with length given by the value on to the stack, which means we can only push at most 0x4b(75) bytes
on the stack, what if we want to push more bytes on to stack? for such situation we need three special operation, if the value with given byte is 0x4c(76) which means it is an instruction OP_PUSHDATA1,
then we need to read the value of following byte to how how many bytes we need to push on to the stack, the value of the following byte is in the range of [0x4b, 0xff], if the byte takes value 77, which means
it is an instruction OP_PUSHDATA2, then we need to read the following two bytes in little endian as the length of the bytes we need to push on the stack, 
this instruction eanble us to push at most 520 bytes on to the stack, any data with length more than 520 is not allowed to transfer by the network. 

You may confused by the process above, let's give out the code that may clear the cloud around you head. In script_sig.go we add following code:
```g
package transaction

import (
	"bufio"
)

type ScriptSig struct {
	cmds []byte
}

const (
	SCRIPT_DATA_LENGTH_BEGIN = 1
	SCRIPT_DATA_LENGTH_END   = 75
	OP_PUSHDATA1             = 76
	OP_PUSHDATA2             = 77
)

func NewScriptSig(reader *bufio.Reader) *ScriptSig {
	cmds := []byte{}
	/*
		At the beginning is the total length for script field, and
		we need to get the length just like how we get the transaction
		input count
	*/
	scriptLen := ReadVarint(reader).Int64()
	count := int64(0)
	current := make([]byte, 1)
	var current_byte byte
	for count < scriptLen {
		reader.Read(current)
		count += 1
		current_byte = current[0]
		if current_byte >= SCRIPT_DATA_LENGTH_BEGIN &&
			current_byte <= SCRIPT_DATA_LENGTH_END {
			//push the following bytes as data on the stack
			data := make([]byte, current_byte)
			reader.Read(data)
			cmds = append(cmds, data...)
			count += int64(current_byte)
		} else if current_byte == OP_PUSHDATA1 {
			/*
				it is instruction OP_PUSHDATA1, we need to read the following
				bytes as the length of data, then push following bytes with
				given length on to stack
			*/
			length := make([]byte, 1)
			reader.Read(length)
			data := make([]byte, length[0])
			reader.Read(data)
			cmds = append(cmds, data...)
			count += int64(length[0] + 1)
		} else if current_byte == OP_PUSHDATA2 {
			/*
				we need to read the following 2 bytes in little endian
				as the length of data byte we need to read at following
			*/
			lenBuf := make([]byte, 2)
			reader.Read(lenBuf)
			length := LittleEndianToBigInt(lenBuf, LITTLE_ENDIAN_2_BYTES)
			data := make([]byte, length.Int64())
			cmds = append(cmds, data...)
			count += length.Int64() + 2
		} else {
			//current byte is an instruction
			cmds = append(cmds, current_byte)

		}
	}

	if count != scriptLen {
		panic("parsing script field failed")
	}

	return &ScriptSig{
		cmds: cmds,
	}
}

```








