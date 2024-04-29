
The following code is used to parse raw binary data into a signature object, this will be used in design the bitcoin script virtual machine, in util.go of 
package elliptic-curve, add the following code:
```g
func ParseSigBin(sigBin []byte) *Signature {
	reader := bytes.NewReader(sigBin)
	bufReader := bufio.NewReader(reader)
	//first byte should be 0x30
	firstByte := make([]byte, 1)
	bufReader.Read(firstByte)
	if firstByte[0] != 0x30 {
		panic("Bad Signature, fisrt byte is not 0x30")
	}
	//second byte is the length of r and s
	lenBuf := make([]byte, 1)
	bufReader.Read(lenBuf)
	//first two byte with the length of r and s should be the total length of sigBin
	if lenBuf[0]+2 != byte(len(sigBin)) {
		panic("Bad Signature length")
	}

	//marker 0x02 as the beginning of r
	marker := make([]byte, 1)
	bufReader.Read(marker)
	if marker[0] != 0x02 {
		panic("signature marker for r is not 0x02")
	}
	//following is the length of r bin
	lenBuf = make([]byte, 1)
	bufReader.Read(lenBuf)
	rLength := lenBuf[0]
	rBin := make([]byte, rLength)
	//it may have 0x00 append at the head but it dose not affect the value of r
	bufReader.Read(rBin)
	r := new(big.Int)
	r.SetBytes(rBin)
	//markder 0x02 for the beginning of s
	marker = make([]byte, 1)
	bufReader.Read(marker)
	if marker[0] != 0x02 {
		panic("signature marker for s is not 0x02")
	}
	//following is length of s bin
	lenBuf = make([]byte, 1)
	bufReader.Read(lenBuf)
	sLength := lenBuf[0]
	sBin := make([]byte, sLength)
	bufReader.Read(sBin)
	s := new(big.Int)
	s.SetBytes(sBin)
	if len(sigBin) != int(6+rLength+sLength) {
		panic("signature wrong length")
	}

	n := GetBitcoinValueN()
	return NewSignature(NewFieldElement(n, r), NewFieldElement(n, s))
}

```
Then in main.go we used following code to test the code above:
```g
package main

import (
	ecc "elliptic_curve"
	"fmt"
	"math/big"
)

func main() {
	n := ecc.GetBitcoinValueN()
	rVal := new(big.Int)
	rVal.SetString("37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6", 16)
	rField := ecc.NewFieldElement(n, rVal)

	sVal := new(big.Int)
	sVal.SetString("8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec", 16)
	sField := ecc.NewFieldElement(n, sVal)

	sig := ecc.NewSignature(rField, sField)
	sigDER := sig.Der()
	sig2 := ecc.ParseSigBin(sigDER)
	fmt.Printf("signature parsed from raw binary data is: %s\n", sig2)
}

```

Running the above code we get the following result:
```g
signature parsed from raw binary data is: Signature(r: {FieldElement{order: fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f, num: 37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6}}, s:{FieldElement{order: fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f, num: 8ca63759c1157ebeaec0d03cecca119fc9a75bf8e6d0fa65c841c8e2738cdaec}})
```
check with the rVal and sVal, we can sure the parsing is correct.
