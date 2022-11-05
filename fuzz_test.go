package base45

import (
	"bytes"
	"testing"
)

func FuzzDecode(f *testing.F) {
	f.Add([]byte("BB8"))
	f.Add([]byte("%69 VD92EX0"))
	f.Add([]byte("UJCLQE7W581"))
	f.Add([]byte("QED8WEX0"))

	f.Fuzz(func(t *testing.T, data []byte) {
		decoded := make([]byte, DecodedLen(len(data)))
		n, err := Decode(decoded, data)
		if err != nil {
			return
		}
		decoded = decoded[:n]

		encoded := make([]byte, EncodedLen(len(decoded)))
		Encode(encoded, decoded)
		if !bytes.Equal(encoded, data) {
			t.Error("encoded result mismatch")
		}
	})
}
