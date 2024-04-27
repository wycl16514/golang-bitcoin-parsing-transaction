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
