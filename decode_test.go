package base45

import (
	"io"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		src, dst string
	}{
		{"BB8", "AB"},
		{"%69 VD92EX0", "Hello!!"},
		{"UJCLQE7W581", "base-45"},
		{"QED8WEX0", "ietf!"},
	}

	for i, tc := range tests {
		dst := make([]byte, DecodedLen(len(tc.src)))
		if _, err := Decode(dst, []byte(tc.src)); err != nil {
			t.Error(err)
			continue
		}
		if string(dst) != tc.dst {
			t.Errorf("%d: want %q, got %q", i, tc.dst, string(dst))
		}
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		src, dst string
	}{
		{"BB8", "AB"},
		{"%69 VD92EX0", "Hello!!"},
		{"UJCLQE7W581", "base-45"},
		{"QED8WEX0", "ietf!"},
	}

	for i, tc := range tests {
		dst, err := DecodeString(tc.src)
		if err != nil {
			t.Error(err)
			continue
		}
		if string(dst) != tc.dst {
			t.Errorf("%d: want %q, got %q", i, tc.dst, string(dst))
		}
	}
}

func BenchmarkDecode(b *testing.B) {
	src := make([]byte, EncodedLen(8192))
	dst := make([]byte, 8192)
	Encode(src, dst)

	b.SetBytes(int64(len(src)))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := Decode(dst, src); err != nil {
			b.Error(err)
		}
	}
}

func TestDecoder(t *testing.T) {
	tests := []struct {
		src, dst string
	}{
		{"BB8", "AB"},
		{"%69 VD92EX0", "Hello!!"},
		{"UJCLQE7W581", "base-45"},
		{"QED8WEX0", "ietf!"},
	}

	for i, tc := range tests {
		r := strings.NewReader(tc.src)
		dec := NewDecoder(r)
		dst, err := io.ReadAll(dec)
		if err != nil {
			t.Error(err)
			continue
		}
		if string(dst) != tc.dst {
			t.Errorf("%d: want %q, got %q", i, tc.dst, string(dst))
		}
	}
}
