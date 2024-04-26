Here is the implementation of the reverse of the ReadVarint function, in util.go of transaction add the following code here:
```g
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
