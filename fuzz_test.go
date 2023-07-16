//go:build go1.18

package base45

import (
	"bytes"
	"testing"
)

func FuzzEncode(f *testing.F) {
	f.Add([]byte("AB"))
	f.Add([]byte("Hello!!"))
	f.Add([]byte("base-45"))
	f.Add([]byte("ietf!"))
	f.Add([]byte("some data with \x00 and \ufeff"))
	f.Add([]byte("Hello, world!"))

	f.Fuzz(func(t *testing.T, data []byte) {
		encoded := make([]byte, EncodedLen(len(data)))
		Encode(encoded, data)

		decoded := make([]byte, DecodedLen(len(encoded)))
		n, err := Decode(decoded, encoded)
		if err != nil {
			t.Error(err)
		}
		decoded = decoded[:n]
		if !bytes.Equal(decoded, data) {
			t.Error("decoded result mismatch")
		}
	})
}

func FuzzDecode(f *testing.F) {
	f.Add([]byte("BB8"))
	f.Add([]byte("%69 VD92EX0"))
	f.Add([]byte("UJCLQE7W581"))
	f.Add([]byte("QED8WEX0"))
	f.Add([]byte("VQEF$DC44IECOCCE4FAWE2249440/DG743XN"))
	f.Add([]byte("%69 VDK2EV4404ESVDX0"))

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
