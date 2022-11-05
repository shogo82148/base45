package base45

var encodeTable = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ $%*+-./:"

func Encode(dst, src []byte) {
	var pos int
	for i := 0; i+1 < len(src); i += 2 {
		v := int(src[i])*256 + int(src[i+1])
		dst[pos] = encodeTable[v%45]
		v /= 45
		dst[pos+1] = encodeTable[v%45]
		v /= 45
		dst[pos+2] = encodeTable[v]
		pos += 3
	}
	if len(src)%2 != 0 {
		v := src[len(src)-1]
		dst[pos] = encodeTable[v%45]
		dst[pos+1] = encodeTable[v/45]
	}
}

func EncodedLen(n int) int {
	return n/2*3 + n%2*2
}
