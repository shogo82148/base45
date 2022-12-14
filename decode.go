package base45

import (
	"errors"
	"io"
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
		if a < 0 || b < 0 {
			return 0, errors.New("base45: invalid character")
		}
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

// DecodedLen returns the maximum length in bytes of the decoded data
// corresponding to n bytes of base45-encoded data.
func DecodedLen(n int) int {
	return n/3*2 + // triple chars are decoded to a byte pair
		n%3/2 // pair chars are decode a byte
}

type decoder struct {
	r io.Reader

	err     error
	readErr error      // error from r.Read
	buf     [1024]byte // leftover input
	nbuf    int
	out     []byte // leftover decoded output
	outbuf  [1024 / 3 * 2]byte
}

func NewDecoder(r io.Reader) io.Reader {
	return &decoder{r: r}
}

func (dec *decoder) Read(buf []byte) (n int, err error) {
	// Use leftover decoded output from last read.
	if len(dec.out) > 0 {
		n = copy(buf, dec.out)
		dec.out = dec.out[n:]
		return n, nil
	}
	if dec.err != nil {
		return 0, dec.err
	}

	// Refill buffer.
	for dec.nbuf < 3 && dec.readErr == nil {
		nn := len(buf) / 2 * 3
		if nn < 3 {
			nn = 3
		}
		if nn > len(dec.buf) {
			nn = len(dec.buf)
		}
		nn, dec.readErr = dec.r.Read(dec.buf[dec.nbuf:nn])
		dec.nbuf += nn
	}

	if dec.nbuf < 3 {
		// decode final fragment
		var nw int
		nw, dec.err = Decode(dec.outbuf[:], dec.buf[:dec.nbuf])
		dec.nbuf = 0
		dec.out = dec.outbuf[:nw]
		n = copy(buf, dec.out)
		dec.out = dec.out[n:]
		if n > 0 || len(buf) == 0 && len(dec.out) > 0 {
			return n, nil
		}
		if dec.err != nil {
			return 0, dec.err
		}
		dec.err = dec.readErr
		if errors.Is(dec.err, io.EOF) && dec.nbuf > 0 {
			dec.err = io.ErrUnexpectedEOF
		}
		return 0, dec.err
	}

	// Decode chunk into p, or d.out and then p if p is too small.
	nr := dec.nbuf / 3 * 3
	nw := dec.nbuf / 3 * 2
	if nw > len(buf) {
		nw, dec.err = Decode(dec.outbuf[:], dec.buf[:nr])
		dec.out = dec.outbuf[:nw]
		n = copy(buf, dec.out)
		dec.out = dec.out[n:]
	} else {
		n, dec.err = Decode(buf, dec.buf[:nr])
	}
	dec.nbuf -= nr
	copy(dec.buf[:dec.nbuf], dec.buf[nr:])
	return n, dec.err
}
