package base45

import "testing"

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

func BenchmarkDecode(b *testing.B) {
	src := make([]byte, EncodedLen(2953))
	dst := make([]byte, 2953)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		Decode(dst, src)
	}
}
