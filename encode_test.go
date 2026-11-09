package base45

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		src, dst string
	}{
		// examples from RFC 9285
		{"AB", "BB8"},
		{"Hello!!", "%69 VD92EX0"},
		{"base-45", "UJCLQE7W581"},
		{"ietf!", "QED8WEX0"},

		// examples
		{"some data with \x00 and \ufeff", "VQEF$DC44IECOCCE4FAWE2249440/DG743XN"},
		{"Hello, world!", "%69 VDK2EV4404ESVDX0"},
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
	src := make([]byte, 8192)
	dst := make([]byte, EncodedLen(len(src)))
	b.SetBytes(int64(len(src)))
	for i := 0; i < b.N; i++ {
		Encode(dst, src)
	}
}

func TestEncoder(t *testing.T) {
	tests := []struct {
		src, dst string
	}{
		// examples from RFC 9285
		{"AB", "BB8"},
		{"Hello!!", "%69 VD92EX0"},
		{"base-45", "UJCLQE7W581"},
		{"ietf!", "QED8WEX0"},

		// long input
		{strings.Repeat("AB", 1024), strings.Repeat("BB8", 1024)},
	}

	for i, tc := range tests {
		var dst bytes.Buffer
		enc := NewEncoder(&dst)

		n, err := enc.Write([]byte(tc.src))
		if err != nil {
			t.Error(err)
			continue
		}
		if n != len(tc.src) {
			t.Errorf("%d: unexpected wrote bytes: %d", i, n)
		}

		if err := enc.Close(); err != nil {
			t.Error(err)
			continue
		}

		if dst.String() != tc.dst {
			t.Errorf("unexpected result: want %q, got %q", tc.dst, dst.String())
		}
	}
}

func TestEncoder_EmptyWriteAfterOddWrite(t *testing.T) {
	var dst bytes.Buffer
	enc := NewEncoder(&dst)

	if n, err := enc.Write([]byte("A")); err != nil || n != 1 {
		t.Fatalf("Write(%q) = %d, %v; want 1, nil", "A", n, err)
	}
	if n, err := enc.Write(nil); err != nil || n != 0 {
		t.Fatalf("Write(nil) = %d, %v; want 0, nil", n, err)
	}
	if n, err := enc.Write([]byte{}); err != nil || n != 0 {
		t.Fatalf("Write([]byte{}) = %d, %v; want 0, nil", n, err)
	}
	if n, err := enc.Write([]byte("B")); err != nil || n != 1 {
		t.Fatalf("Write(%q) = %d, %v; want 1, nil", "B", n, err)
	}
	if err := enc.Close(); err != nil {
		t.Fatal(err)
	}
	if got, want := dst.String(), "BB8"; got != want {
		t.Errorf("unexpected result: want %q, got %q", want, got)
	}
}
