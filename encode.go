package base45

import "io"

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

// EncodeToString returns the base45 encoding of src.
func EncodeToString(src []byte) string {
	dst := make([]byte, EncodedLen(len(src)))
	Encode(dst, src)
	return string(dst)
}

// EncodedLen returns the length in bytes of the base45 encoding
// of an input buffer of length n.
func EncodedLen(n int) int {
	return n/2*3 + // 3 chars for each byte pair
		n%2*2 // 2 chars for the remain.
}

type encoder struct {
	err  error
	w    io.Writer
	buf  [2]byte    // buffered data waiting to be encoded
	nbuf int        // number of bytes in buf
	out  [1024]byte // output buffer
}

func NewEncoder(w io.Writer) io.WriteCloser {
	return &encoder{w: w}
}

func (enc *encoder) Write(data []byte) (n int, err error) {
	if enc.err != nil {
		return 0, enc.err
	}

	// Leading fringe.
	if enc.nbuf > 0 {
		enc.buf[enc.nbuf] = data[0]
		enc.nbuf++
		n++
		data = data[1:]
		Encode(enc.out[:], enc.buf[:])
		if _, enc.err = enc.w.Write(enc.out[:3]); enc.err != nil {
			return n, enc.err
		}
		enc.nbuf = 0
	}

	// Large interior chunks.
	for len(data) >= 2 {
		nn := len(enc.out) / 3 * 2
		if nn > len(data) {
			nn = len(data)
			nn -= nn % 2
		}
		Encode(enc.out[:], data[:nn])
		if _, enc.err = enc.w.Write(enc.out[0 : nn/2*3]); enc.err != nil {
			return n, enc.err
		}
		n += nn
		data = data[nn:]
	}

	// Trailing fringe.
	if len(data) > 0 {
		enc.buf[0] = data[0]
		enc.nbuf = 1
		n++
	}
	return
}

// Close flushes any pending output from the encoder.
// It is an error to call Write after calling Close.
func (enc *encoder) Close() error {
	if enc.err != nil {
		return enc.err
	}

	// If there's anything left in the buffer, flush it out
	if enc.nbuf > 0 {
		Encode(enc.out[:], enc.buf[:1])
		_, enc.err = enc.w.Write(enc.out[:2])
	}
	return nil
}
