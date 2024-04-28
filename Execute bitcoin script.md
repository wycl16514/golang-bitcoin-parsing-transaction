Bitcoin script is a kind of programming langauge just like assembly language. It is used to complete some kind of task, for example some script are used to verify you are the legitimate of receipient
of given sum of bitcoin. Therefore we need to know the details about bitcoin script instructions. As we metioned before, there are two kinds of instrution, one is used to moving data around stack, the
other is used to take data from stack and do some computation.

The first kind of instruction like OP_0, OP_1, OP_2, OP_1NEGATE are used to push value 0, 1, -1 on to the stack, when push value on to stack, we need to convert the value into byte array in the order of 
little endian, and we need to pay attention is that, when we push negative value , we need to convert them into byte array with two's complement, this encoding require to set the hightest bit to 1 for 
negative number, if the aboslute value of the negative number has its most significant bit setted, we need to insert a byte with value 0x80(the most significant bit setted) at the head of the byte array.

For example for -1, its positive counterpart is 1 with binary foramt 0000 0001, the most significant one is 0, then we set 1 to 1 its the binary format of -1, that is 1000 0001, the hex value is 0x81.
for value -1233, its positive counterpart is 1234, its binary format is 0000 0100 1101 0010, then for -1234, we need to set the significant bit to 1 that is 1000 0100 1101 0010 and its hex byte array is 
0x84 0xD2, for another example -32896, its positive counter part is 32896 with hex byte array (0x80, 0x80) and its binary format is 1000 0000 1000 0000, we can see its most significant bit is already 1,
this time we need to insert 0x80 at the head which is 1000 0000 1000 0000 1000 0000, in hex byte array is 0x80, 0x80, 0x80. 

Let's see how to use code to implement this scheme, create a new file called op.go and add code like following:
```g
package transaction

type BitcoinOpCode struct {
}

func NewBitcoinOpCode() *BitcoinOpCode {
	return &BitcoinOpCode{}
}

func (b *BitcoinOpCode) EncodeNum(num int64) []byte {
	/*
		convert value to its hex byte array in little endian order,
		for negative number use two's complement
	*/
	if num == 0 {
		/*
			for 0 we don't push [0x00], but push a
			empty byte string
		*/
		return []byte("")
	}

	result := []byte{}
	absNum := num
	negative := false
	if num < 0 {
		absNum = -num
		negative = true
	}

	for absNum > 0 {
		/*
			append the last byte of asbNum into result,
			notices result will be little endian byte array of
			absNum, and the most significant byte is at the end
		*/
		result = append(result, byte(absNum&0xff))
		absNum >>= 8
	}

	/*
		check the most significant bit, notice the most
		significant byte is at the end of result
	*/
	if (result[len(result)-1] & 0x80) != 0 {
		if negative {
			/*
				for negative value with most significant bit set to 1
				we need to append 0x80 at the head for two's complement,
				because result already in littel endian, then we append
				the 0x80 byte at the end
			*/
			result = append(result, 0x80)
		} else {
			result = append(result, 0x00)
		}
	} else if negative {
		/*
			it is negative value but the most significant bit is not 1,
			then we just set the most significant bit to 1
		*/
		result[len(result)-1] |= 0x80
	}

	return result
}

func (b *BitcoinOpCode) DecodeNum(element []byte) int64 {
        //check empty byte string
	if len(element) == 0 {
		return 0
	}
	//reverse of EncodeNum
	//reverse the byte array to big endian
	bigEndian := reverseByteSlice(element)
	negative := false
	result := int64(0)
	//if the most significant bit set to 1, it is negative value
	if (bigEndian[0] & 0x80) != 0 {
		negative = true
		//reset the most significant bit to 0
		//0x7f is 0111 1111
		result = int64(bigEndian[0] & 0x7f)
	} else {
		negative = false
		result = int64(bigEndian[0])
	}

	for i := 1; i < len(bigEndian); i++ {
		result <<= 8
		result += int64(bigEndian[i])
	}

	if negative {
		return -result
	}

	return result
}

```

Then we can write code to test the above code in main.go like following:
```g
func main() {
	opCode := tx.NewBitcoinOpCode()
	encodeVal := opCode.EncodeNum(-1)
	fmt.Printf("encode -1: %x\n", encodeVal)
	fmt.Printf("decode 1: %d\n", opCode.DecodeNum(encodeVal))

	encodeVal = opCode.EncodeNum(-1234)
	fmt.Printf("encode -1: %x\n", encodeVal)
	fmt.Printf("decode 1: %d\n", opCode.DecodeNum(encodeVal))

	encodeVal = opCode.EncodeNum(-32896)
	fmt.Printf("encode -1: %x\n", encodeVal)
	fmt.Printf("decode 1: %d\n", opCode.DecodeNum(encodeVal))
}

```
Running the above code can get the following result:
```g
encode -1: 81
decode 1: -1
encode -1: d284
decode 1: -1234
encode -1: 808080
decode 1: -32896
```
The result meet our expectation and prove the logic of our code should be correct.

There are details for all the bitcoin script operation code in here: https://opcodeexplained.com/, we will go to implement some of them, when you finish this section, you will have the capability to
implement those op code that are not implemented here. First let's init those op codes that will implemented by our project, in op.go add the following code : 
```g
const (
	OP_0 = 0
)
const (
	OP_1NEGATE = iota + 79
)
const (
	OP_1 = iota + 81
	OP_2
	OP_3
	OP_4
	OP_5
	OP_6
	OP_7
	OP_8
	OP_9
	OP_10
	OP_11
	OP_12
	OP_13
	OP_14
	OP_15
	OP_16
	OP_NOP
)

const (
	OP_IF = iota + 99
	OP_NOTIf
)

const (
	OP_VERIFY = iota + 105
	OP_RETURN
	OP_TOTALSTACK
	OP_FROMALTSTACK
	OP_2DROP
	OP_2DUP
	OP_3DUP
	OP_2OVER
	OP_2ROT
	OP_2SWAP
	OP_IFDUP
	OP_DEPTH
	OP_DROP
	OP_DUP
	OP_NIP
	OP_OVER
	OP_PICK
	OP_ROLL
	OP_ROT
	OP_SWAP
	OP_TUCK
)

const (
	OP_SIZE = iota + 130
)

const (
	OP_EQUAL = iota + 135
	OP_EQUALVERIFY
)

const (
	OP_1ADD = iota + 139
	OP_1SUB
)

const (
	OP_NEGATE = iota + 143
	OP_ABS
	OP_NOT
	OP_0NOTEQUAL
	OP_ADD
	OP_SUB
	OP_MUL
)

const (
	OP_BOOLAND = iota + 154
	OP_BOOLOR
	OP_NUMEQUAL
	OP_NUMEQUALVERIFY
	OP_NUMNOTEQUAL
	OP_LESSTHAN
	OP_GREATERTHAN
	OP_LESSTHANOREQUAL
	OP_GREATERTHANOREQUAL
	OP_MIN
	OP_MAX
	OP_WITHIN
	OP_RIPEMD160
	OP_SHA1
	OP_SHA256
	OP_HASH160
	OP_HASH256
)

const (
	OP_CHECKSIG = iota + 172
	OP_HECKSIGVERIFY
	OP_CHECKMULTISIG
	OP_CHECKMULTISIGVERIFY
	OP_NOP1
	OP_CHECKLOGTIMEVERIFY
	OP_CHECKSEQUENCEVERIFY
	OP_NOP4
	OP_NOP5
	OP_NOP6
	OP_NOP7
	OP_NOP8
	OP_NOP9
	OP_NOP10
)

type BitcoinOpCode struct {
	opCodeNames map[int]string
}

func NewBicoinOpCode() *BitcoinOpCode {
	opCodeNames := map[int]string{
		0:   "OP_0",
		76:  "OP_PUSHDATA1",
		77:  "OP_PUSHDATA2",
		78:  "OP_PUSHDATA4",
		79:  "OP_1NEGATE",
		81:  "OP_1",
		82:  "OP_2",
		83:  "OP_3",
		84:  "OP_4",
		85:  "OP_5",
		86:  "OP_6",
		87:  "OP_7",
		88:  "OP_8",
		89:  "OP_9",
		90:  "OP_10",
		91:  "OP_11",
		92:  "OP_12",
		93:  "OP_13",
		94:  "OP_14",
		95:  "OP_15",
		96:  "OP_16",
		97:  "OP_NOP",
		99:  "OP_IF",
		100: "OP_NOTIF",
		103: "OP_ELSE",
		104: "OP_ENDIF",
		105: "OP_VERIFY",
		106: "OP_RETURN",
		107: "OP_TOALTSTACK",
		108: "OP_FROMALTSTACK",
		109: "OP_2DROP",
		110: "OP_2DUP",
		111: "OP_3DUP",
		112: "OP_2OVER",
		113: "OP_2ROT",
		114: "OP_2SWAP",
		115: "OP_IFDUP",
		116: "OP_DEPTH",
		117: "OP_DROP",
		118: "OP_DUP",
		119: "OP_NIP",
		120: "OP_OVER",
		121: "OP_PICK",
		122: "OP_ROLL",
		123: "OP_ROT",
		124: "OP_SWAP",
		125: "OP_TUCK",
		130: "OP_SIZE",
		135: "OP_EQUAL",
		136: "OP_EQUALVERIFY",
		139: "OP_1ADD",
		140: "OP_1SUB",
		143: "OP_NEGATE",
		144: "OP_ABS",
		145: "OP_NOT",
		146: "OP_0NOTEQUAL",
		147: "OP_ADD",
		148: "OP_SUB",
		149: "OP_MUL",
		154: "OP_BOOLAND",
		155: "OP_BOOLOR",
		156: "OP_NUMEQUAL",
		157: "OP_NUMEQUALVERIFY",
		158: "OP_NUMNOTEQUAL",
		159: "OP_LESSTHAN",
		160: "OP_GREATERTHAN",
		161: "OP_LESSTHANOREQUAL",
		162: "OP_GREATERTHANOREQUAL",
		163: "OP_MIN",
		164: "OP_MAX",
		165: "OP_WITHIN",
		166: "OP_RIPEMD160",
		167: "OP_SHA1",
		168: "OP_SHA256",
		169: "OP_HASH160",
		170: "OP_HASH256",
		171: "OP_CODESEPARATOR",
		172: "OP_CHECKSIG",
		173: "OP_CHECKSIGVERIFY",
		174: "OP_CHECKMULTISIG",
		175: "OP_CHECKMULTISIGVERIFY",
		176: "OP_NOP1",
		177: "OP_CHECKLOCKTIMEVERIFY",
		178: "OP_CHECKSEQUENCEVERIFY",
		179: "OP_NOP4",
		180: "OP_NOP5",
		181: "OP_NOP6",
		182: "OP_NOP7",
		183: "OP_NOP8",
		184: "OP_NOP9",
		185: "OP_NOP10",
	}
	return &BitcoinOpCode{
		opCodeNames: opCodeNames,
	}
}
```
We define the value with its op code, you can check the link above to verify whether the value is right for the given op code.
