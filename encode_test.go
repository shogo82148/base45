package base45

import "testing"

func TestEncode(t *testing.T) {
	tests := []struct {
		src, dst string
	}{
		{"AB", "BB8"},
		{"Hello!!", "%69 VD92EX0"},
		{"base-45", "UJCLQE7W581"},
		{"ietf!", "QED8WEX0"},
	}

	for i, tc := range tests {
		dst := make([]byte, EncodedLen(len(tc.src)))
		Encode(dst, []byte(tc.src))
		if string(dst) != tc.dst {
			t.Errorf("%d: want %q, got %q", i, tc.dst, string(dst))
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	src := make([]byte, 2953)
	dst := make([]byte, EncodedLen(len(src)))
	for i := 0; i < b.N; i++ {
		Encode(dst, src)
	}
}
