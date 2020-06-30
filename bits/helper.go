package bits

//GetNybbles Returns the Nybbles for a slice of bytes
func GetNybbles(in []byte) []byte {
	output := []byte{}
	for _, v := range in {
		output = append(output, in[v]>>4)
		output = append(output, in[v]&0x0f)
	}
	return output
}

//GetNybble Returns the Nybbles for a byte
func GetNybble(in byte) []byte {
	output := []byte{}
	output = append(output, in>>4)
	output = append(output, in&0x0f)
	return output
}

func BytesToUInt(b []byte)uint{
	out uint := 
	for i, v := range b {
		j := (i % 2) 

	} 
}
