package base45_test

import (
	"fmt"
	"os"

	"github.com/shogo82148/base45"
)

func Example() {
	msg := "Hello, 世界"
	encoded := base45.EncodeToString([]byte(msg))
	fmt.Println(encoded)
	decoded, err := base45.DecodeString(encoded)
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	fmt.Println(string(decoded))
	// Output:
	// %69 VDK2E5744FNKCT53
	// Hello, 世界
}

func ExampleEncodeToString() {
	data := []byte("any + old & data")
	str := base45.EncodeToString(data)
	fmt.Println(str)
	// Output:
	// CEC3EFFK5*3ERTC+ 42VC3WE
}

func ExampleEncode() {
	data := []byte("Hello, world!")
	dst := make([]byte, base45.EncodedLen(len(data)))
	base45.Encode(dst, data)
	fmt.Println(string(dst))
	// Output:
	// %69 VDK2EV4404ESVDX0
}

func ExampleDecodeString() {
	str := "VQEF$DC44IECOCCE4FAWE2249440/DG743XN"
	data, err := base45.DecodeString(str)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("%q\n", data)
	// Output:
	// "some data with \x00 and \ufeff"
}

func ExampleDecode() {
	str := "%69 VDK2EV4404ESVDX0"
	dst := make([]byte, base45.DecodedLen(len(str)))
	n, err := base45.Decode(dst, []byte(str))
	if err != nil {
		fmt.Println("decode error:", err)
		return
	}
	dst = dst[:n]
	fmt.Printf("%q\n", dst)
	// Output:
	// "Hello, world!"
}

func ExampleNewEncoder() {
	input := []byte("foo\x00bar")
	encoder := base45.NewEncoder(os.Stdout)
	encoder.Write(input)
	// Must close the encoder when finished to flush any partial blocks.
	// If you comment out the following line, the last partial block "r"
	// won't be encoded.
	encoder.Close()

	// Output:
	// X.CL1EUJCO2
}
