package base45

import (
	"errors"
)

var decodeTable = [256]int{
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	36, -1, -1, -1, 37, 38, -1, -1,
	-1, -1, 39, 40, -1, 41, 42, 43,
	0, 1, 2, 3, 4, 5, 6, 7, 8,
	9, 44, -1, -1, -1, -1, -1,
	-1, 10, 11, 12, 13, 14, 15, 16,
	17, 18, 19, 20, 21, 22, 23, 24,
	25, 26, 27, 28, 29, 30, 31, 32,
	33, 34, 35, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
	-1, -1, -1, -1, -1, -1, -1, -1,
}

func Decode(dst, src []byte) (n int, err error) {
	if len(src)%3 == 1 {
		return 0, errors.New("base45: invalid length")
	}

	pos := 0
	for i := 0; i+2 < len(src); i += 3 {
		a := decodeTable[src[i]]
		b := decodeTable[src[i+1]]
		c := decodeTable[src[i+2]]
		if a < 0 || b < 0 || c < 0 {
			return 0, errors.New("base45: invalid character")
		}
		v := a + b*45 + c*45*45
		if v >= 0x10000 {
			return 0, errors.New("base45: invalid value")
		}
		dst[pos+1] = byte(v % 256)
		v /= 256
		dst[pos] = byte(v)
		pos += 2
	}
	if len(src)%3 == 2 {
		a := decodeTable[src[len(src)-2]]
		b := decodeTable[src[len(src)-1]]
		v := a + b*45
		if v >= 0x100 {
			return 0, errors.New("base45: invalid value")
		}
		dst[pos] = byte(v)
		pos++
	}
	return pos, nil
}

// DecodeString returns the bytes represented by the base45 string s.
func DecodeString(s string) ([]byte, error) {
	buf := make([]byte, DecodedLen(len(s)))
	n, err := Decode(buf, []byte(s))
	return buf[:n], err
}

func DecodedLen(n int) int {
	return n/3*2 + n%3/2
}
